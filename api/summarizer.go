package api

import (
	"context"
	"encoding/json"
	"log"

	"github.com/tidepool-org/summary/data"
	"github.com/tidepool-org/summary/debezium"
)

//Summarizer creates summaries
type Summarizer struct {
	Histograms []*Histogramer // histograms for each time range
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

//Run runs the summarizer until the context is closed or the input channel is closed
func (s *Summarizer) Run(ctx context.Context, in <-chan *debezium.MongoDBEvent, out chan<- SummaryResponse) {
	defer close(out)

	for {
		select {
		case <-ctx.Done():
			return
		case s := <-in:
			switch s.Payload.Op {
			case "c":

			case "r":
			case "u":
			case "d":
			}
		}
	}
}

//Add add an event
func (s *Summarizer) Add(rec *debezium.MongoDBEvent) {
	var d data.Blood
	if err := json.Unmarshal([]byte(rec.Payload.After), &d); err != nil {
		log.Println("Error Unmarshalling after field", err)
	} else {
		if d.Type == "cbg" || d.Type == "smbg" {
			log.Printf("%v\n", d)
			if d.Value == nil || d.Units == nil {
				log.Printf("skipping entry with missing value or units %v\n", d)
				return
			}
			//standardized := s.Normalizer.ToStandard(float32(*d.Value), *d.Units) XXX
			// place into all matching bins
		} else {
			log.Printf("skipping type %v\n", d.Type)
		}
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
