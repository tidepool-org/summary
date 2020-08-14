package summarizer

import (
	"log"
	"time"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
	"github.com/tidepool-org/summary/histogram"
	"github.com/tidepool-org/summary/normalize"
)

//Summarizer creates summaries
type Summarizer struct {
	Histograms []*histogram.Histogramer // histograms for each time range
	Periods    []api.SummaryPeriod
	Normalizer normalize.UnitNormalizer
	Request    api.SummaryRequest
}

// NewSummarizer creates a Summarizer for the given request
func NewSummarizer(request api.SummaryRequest) *Summarizer {
	histograms := make([]*histogram.Histogramer, request.Period.NumPeriods)
	quantiles := make([]histogram.QuantileInfo, len(request.Quantiles))

	for i := range quantiles {
		quantiles[i].Name = request.Quantiles[i].Name
		quantiles[i].Threshold = float64(request.Quantiles[i].Threshold)
	}

	periods := make([]api.SummaryPeriod, request.Period.NumPeriods)

	now := time.Now()
	ending := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
	var duration time.Duration
	switch request.Period.Length {
	case "day":
		duration = 24 * time.Hour
	case "week":
		duration = 7 * 24 * time.Hour
	}

	for i := range histograms {
		histograms[i] = histogram.NewHistogramer(quantiles)
		periods[i].End = ending
		ending = ending.Add(-1 * duration)
		periods[i].Start = ending
		periods[i].Length = request.Period.Length
	}

	return &Summarizer{
		Histograms: histograms,
		Request:    request,
		Normalizer: &normalize.BloodGlucoseNormalizer{},
		Periods:    periods,
	}
}

//Process an event
func (s *Summarizer) Process(rec interface{}) {

	now := time.Now()
	switch v := rec.(type) {
	case data.Upload:
		uploadId := v.UploadID
	case data.Blood:
		if v.Value == nil || v.Units == nil {
			log.Printf("skipping entry with missing value or units %v\n", v)
			return
		}
		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, *v.Time)

		if err != nil {
			log.Printf("skipping entry with bad date %v\n", v)
			return
		}

		standardized := s.Normalizer.ToStandard(float32(*v.Value), *v.Units)

		if v.Active {
			for i, p := range s.Periods {
				if (!t.Before(p.Start)) && p.End.After(t) {
					s.Histograms[i].Add(float64(standardized))
					p.Updated = now
				}
			}
		}
	default:
		log.Printf("skipping  %v\n", v)
	}
}

//Summary return summary report
func (s *Summarizer) Summary() api.SummaryResponse {

	if s.Count == 0 {
		return &api.SummaryResponse{}
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
