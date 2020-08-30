package summarizer

import (
	"fmt"
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
	BGStart    time.Time
	BGEnd      time.Time
	Type       *string
}

// NewGlucoseSummarizer creates a Summarizer for the given request
func NewGlucoseSummarizer(request api.SummaryRequest, periods []api.SummaryPeriod) *GlucoseSummarizer {
	histograms := make([]*Histogramer, len(periods))
	quantiles := quantilesFromRequest(request)
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
func (s *GlucoseSummarizer) Process(v *data.Blood) error {
	if v.Value == nil || v.Units == nil {
		return fmt.Errorf("missing value or units from egv value: %v", v)
	}

	t, err := time.Parse(Layout, *v.Time)
	if err != nil {
		return fmt.Errorf("skipping entry with bad date %v", v)
	}

	if s.BGStart.IsZero() || t.Before(s.BGStart) {
		s.BGStart = t
	}
	if s.BGEnd.IsZero() || t.After(s.BGEnd) {
		s.BGEnd = t
	}

	s.Type = &v.Base.Type

	standardized := s.Normalizer.ToStandard(float32(*v.Value), *v.Units)

	for i, p := range s.Periods {
		if (!t.Before(p.Start)) && p.End.After(t) {
			s.Histograms[i].Add(float64(standardized))
		}
	}
	return nil
}

//Summary return summary report
func (s *GlucoseSummarizer) Summary() []api.GlucoseSummary {
	reports := make([]api.GlucoseSummary, len(s.Periods))
	now := time.Now()
	for i, period := range s.Periods {
		histogram := s.Histograms[i]
		period.Updated = now
		reports[i] = api.GlucoseSummary{
			Period: period,
			Stats: api.SummaryStatistics{
				Count:     int(histogram.OverallCount()),
				Mean:      float32(histogram.OverallMean()),
				Units:     s.Request.Units,
				Quantiles: quantilesFromHistogram(histogram),
			},
		}
	}
	return reports
}

//quantilesFromRequest makes the quantiles for a histogrammer, including a quantile to capture all the data
func quantilesFromRequest(request api.SummaryRequest) []QuantileInfo {
	quantiles := make([]QuantileInfo, len(request.Quantiles))
	for i, requested := range request.Quantiles {
		quantiles[i] = QuantileInfo{Name: requested.Name, Threshold: float64(requested.Threshold)}
	}
	return quantiles
}

func quantilesFromHistogram(histogram *Histogramer) []api.Quantile {
	nBins := histogram.Bins()
	quantiles := make([]api.Quantile, nBins)
	for j := 0; j != nBins; j++ {
		quantiles[j].Count = new(int)
		*quantiles[j].Count = int(histogram.Count(j))
		quantiles[j].Percentage = float32(histogram.Percentage(j))
		quantiles[j].Threshold = float32(histogram.Threshold(j))
		quantiles[j].Name = histogram.Name(j)
	}
	return quantiles
}
