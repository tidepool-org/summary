package summarizer

import (
	"log"
	"time"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
)

//GlucoseSummarizer creates summaries of upload activity
type GlucoseSummarizer struct {
	Histograms []*Histogramer // histograms for each time range
	Periods    []api.SummaryPeriod
	Normalizer GlucoseNormalizer
	Request    api.SummaryRequest
}

// NewGlucoseSummarizer creates a Summarizer for the given request
func NewGlucoseSummarizer(request api.SummaryRequest) *GlucoseSummarizer {
	histograms := make([]*Histogramer, request.Period.NumPeriods)
	quantiles := make([]QuantileInfo, len(request.Quantiles))

	for i := range quantiles {
		quantiles[i].Name = request.Quantiles[i].Name
		quantiles[i].Threshold = float64(request.Quantiles[i].Threshold)
	}

	periods := make([]api.SummaryPeriod, request.Period.NumPeriods)

	now := time.Now()
	ending := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
	var nDays int
	switch request.Period.Length {
	case "day":
		nDays = 1
	case "week":
		nDays = 7
	}

	for i := range histograms {
		histograms[i] = NewHistogramer(quantiles)
		periods[i].End = ending
		ending = ending.AddDate(0, 0, -nDays)
		periods[i].Start = ending
		periods[i].Length = request.Period.Length
	}

	return &GlucoseSummarizer{
		Histograms: histograms,
		Request:    request,
		Normalizer: GlucoseNormalizer{},
		Periods:    periods,
	}
}

//Process a glucose sample
func (s *GlucoseSummarizer) Process(v *data.Blood) {
	now := time.Now()
	if v.Value == nil || v.Units == nil {
		log.Printf("skipping entry with missing value or units %v\n", v)
		return
	}
	t, err := time.Parse(Layout, *v.Time)

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
}

//Summary return summary report
func (s *GlucoseSummarizer) Summary() []api.GlucoseSummary {
	reports := make([]api.GlucoseSummary, len(s.Periods))

	for i, period := range s.Periods {
		histogram := s.Histograms[i]
		quantiles := make([]api.Quantile, len(histogram.Info))
		for j, info := range histogram.Info {
			quantiles[j].Count = new(int)
			*quantiles[j].Count = info.Count
			quantiles[j].Percentage = float32(info.Count) / float32(histogram.Count)
			quantiles[j].Threshold = s.Request.Quantiles[j].Threshold
			quantiles[j].Name = s.Request.Quantiles[j].Name
		}
		reports[i] = api.GlucoseSummary{
			Period: period,
			Stats: api.SummaryStatistics{
				Count:     histogram.Count,
				Mean:      float32(histogram.Mean()),
				Units:     s.Request.Units,
				Quantiles: quantiles,
			},
		}
	}
	return reports
}
