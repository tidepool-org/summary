package api

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/tidepool-org/summary/data"
)

//Summarizer creates summaries
type Summarizer struct {
	Histograms []*Histogramer // histograms for each time range
	Periods    []SummaryPeriod
	Normalizer UnitNormalizer
	Request    SummaryRequest
}

// NewSummarizer creates a Summarizer for the given request
func NewSummarizer(request SummaryRequest) *Summarizer {
	histograms := make([]*Histogramer, request.Period.NumPeriods)
	quantiles := make([]QuantileInfo, len(request.Quantiles))

	for i := range quantiles {
		quantiles[i].Name = request.Quantiles[i].Name
		quantiles[i].Threshold = float64(request.Quantiles[i].Threshold)
	}

	for i := range histograms {
		histograms[i] = NewHistogramer(quantiles)
	}

	return &Summarizer{
		Histograms: histograms,
		Request:    request,
		Normalizer: &BloodGlucoseNormalizer{},
	}
}

// ChangeNotificationEvent is a notification of a change given event time
type ChangeNotificationEvent struct {
	Date   time.Time `json:"date"`   // event time of original record
	UserID string    `json:"userid"` // userid
	Kind   string    `kind:"kind"`   // enum: cbg, smbg, profile
}

//Run runs the summarizer until the context is closed or the input channel is closed
func (s *Summarizer) Run(ctx context.Context, in <-chan *ChangeNotificationEvent, out chan<- SummaryResponse) {
	defer close(out)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-in:
			s.Process(msg)
		}
	}
}

type BloodGlucose struct {
	Date   time.Time
	UserID string
	Value  float32
	Units  string
}

type BGProvider interface {
	Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BloodGlucose)
}

type MongoBGProvider struct {
}

var _ BGProvider = &MongoBGProvider{}

//Get provide BG values on a channel, close channel when no more values
func (b *MongoBGProvider) Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BloodGlucose) {
}

//Process an event
// TODO handle delete and update (usually active to inactive or vice-versa)
func (s *Summarizer) Process(rec *ChangeNotificationEvent) {
	var after data.Blood

	if err := json.Unmarshal([]byte(rec.Payload.After), &after); err != nil {
		log.Println("Error Unmarshalling after field", err)
		return
	}

	if after.Type == "cbg" || after.Type == "smbg" {
		if after.Value == nil || after.Units == nil {
			log.Printf("skipping entry with missing value or units %v\n", after)
			return
		}
		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, *after.Time)

		if err != nil {
			log.Printf("skipping entry with bad date %v\n", after)
			return
		}

		var standardizedAfter float32
		op := rec.Payload.Op
		if op == "r" || op == "c" {
			standardizedAfter = s.Normalizer.ToStandard(float32(*after.Value), *after.Units)
		}

		if after.Active {
			for i, p := range s.Periods {
				if (!t.Before(p.Start)) && p.End.After(t) {
					if op == "r" || op == "c" {
						s.Histograms[i].Add(float64(standardizedAfter))
					}
				}
			}
		}
	} else {
		log.Printf("skipping type %v\n", after.Type)
	}
}

/*
	if s.Count == 0 {
		return &SummaryStatistics{
			Units:     s.Request.Units,
			Quantiles: make([]Quantile, 0),
		}
	}
	quantiles := make([]Quantile, len(s.Request.Quantiles))
	for i, quantile := range s.Quantiles {
		quantiles[i].Count = new(int)
		*quantiles[i].Count = quantile.Count
		quantiles[i].Percentage = float32(quantile.Count) / float32(s.Count)
		quantiles[i].Threshold = s.Request.Quantiles[i].Threshold
		quantiles[i].Name = s.Request.Quantiles[i].Name
	}
	return &SummaryStatistics{
		Count:     s.Count,
		Mean:      float32(s.Sum / float64(s.Count)),
		Units:     s.Request.Units,
		Quantiles: quantiles,
	}
}
*/
