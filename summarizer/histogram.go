package summarizer

// QuantileInfo aggregates partial info leading to a quantile summary
type QuantileInfo struct {
	Count     int
	Sum       float64
	Threshold float64
	Name      string
}

// Histogramer create a summary of blood glucose data
type Histogramer struct {
	Info  []QuantileInfo
	Count int
	Sum   float64
}

// NewHistogramer create a Histogramer
func NewHistogramer(quantiles []QuantileInfo) *Histogramer {
	return &Histogramer{
		Info: quantiles,
	}
}

// Add adds a blood sample to the summary
func (s *Histogramer) Add(value float64) {
	for i, quantile := range s.Info {
		if value < quantile.Threshold {
			s.Info[i].Count++
			s.Info[i].Sum += float64(value)
		}
	}
	s.Count++
	s.Sum += float64(value)
}

// Remove adds a blood sample to the summary
func (s *Histogramer) Remove(value float64) {
	for i, quantile := range s.Info {
		if value < quantile.Threshold {
			s.Info[i].Count--
			s.Info[i].Sum -= float64(value)
		}
	}
	s.Count--
	s.Sum -= float64(value)
}

//Mean value
func (s *Histogramer) Mean() float64 {
	if s.Count > 0 {
		return s.Sum / float64(s.Count)
	}
	return 0.0
}
