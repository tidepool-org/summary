package api

import "github.com/tidepool-org/summary/data"

// ConverstionRatio is the mmol to mg/dL ratio
const ConverstionRatio = 18.0182

// QuantileInfo aggregates partial info leading to a quantile summary
type QuantileInfo struct {
	Count int
	Sum   float64
}

// Summarizer create a summary of blood glucose data
type Summarizer struct {
	Request   *SummaryRequest
	Quantiles []QuantileInfo
	Count     int
	Sum       float64
}

// NewSummarizer create a summarizer
func NewSummarizer(request *SummaryRequest) *Summarizer {
	return &Summarizer{
		Request:   request,
		Quantiles: make([]QuantileInfo, len(request.Quantiles)),
	}
}

// Convert value in units of summarizer
func (s *Summarizer) Convert(val float64, units SummaryRequestUnits) float64 {
	if units == s.Request.Units {
		return val
	}
	if units == SummaryRequestUnits_mmol_L || units == SummaryRequestUnits_mmol_l {
		return val * ConverstionRatio
	}
	return val / ConverstionRatio
}

// Add adds a blood sample to the summary
func (s *Summarizer) Add(d data.Blood) {
	if d.Units == nil || d.Value == nil {
		return
	}
	value := s.Convert(*d.Value, SummaryRequestUnits(*d.Units))
	for i, quantile := range s.Request.Quantiles {
		if value < float64(quantile.Threshold) {
			s.Quantiles[i].Count++
			s.Quantiles[i].Sum += value
		}
	}
	s.Count++
	s.Sum += value
}

// Remove adds a blood sample to the summary
func (s *Summarizer) Remove(d data.Blood) {
	if d.Units == nil || d.Value == nil {
		return
	}
	value := s.Convert(*d.Value, SummaryRequestUnits(*d.Units))
	for i, quantile := range s.Request.Quantiles {
		if value < float64(quantile.Threshold) {
			s.Quantiles[i].Count--
			s.Quantiles[i].Sum -= value
		}
	}
	s.Count--
	s.Sum -= value
}

// Summary returns the summary report
func (s *Summarizer) Summary() *SummaryStatistics {
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
