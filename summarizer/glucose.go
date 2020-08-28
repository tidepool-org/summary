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
	Start      time.Time
	End        time.Time
}

// NewGlucoseSummarizer creates a Summarizer for the given request
func NewGlucoseSummarizer(request api.SummaryRequest, periods []api.SummaryPeriod) *GlucoseSummarizer {
	histograms := make([]*Histogramer, len(periods))
	quantiles := make([]QuantileInfo, len(request.Quantiles))

	for i := range quantiles {
		quantiles[i].Name = request.Quantiles[i].Name
		quantiles[i].Threshold = float64(request.Quantiles[i].Threshold)
	}

	for i := range periods {
		histograms[i] = NewHistogramer(quantiles)
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

	for i, p := range s.Periods {
		if (!t.Before(p.Start)) && p.End.After(t) {
			s.Histograms[i].Add(float64(standardized))
		}
	}
}

//Summary return summary report
func (s *GlucoseSummarizer) Summary() []api.GlucoseSummary {
	reports := make([]api.GlucoseSummary, len(s.Periods))
	now := time.Now()
	for i, period := range s.Periods {
		histogram := s.Histograms[i]
		period.Updated = now
		quantiles := make([]api.Quantile, len(histogram.Info))
		for j, info := range histogram.Info {
			quantiles[j].Count = new(int)
			*quantiles[j].Count = info.Count
			if histogram.Count > 0 {
				quantiles[j].Percentage = (100.0 * float32(info.Count)) / float32(histogram.Count)
			} else {
				quantiles[j].Percentage = 0.0
			}
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
	log.Printf("reports %v", reports)
	return reports
}
