package api

// QuantileInfo aggregates partial info leading to a quantile summary
type QuantileInfo struct {
	Count     int
	Sum       float64
	Threshold float64
	Name      string
}

// Histogramer create a summary of blood glucose data
type Histogramer struct {
	Quantiles []QuantileInfo
	Count     int
	Sum       float64
}

// NewHistogramer create a Histogramer
func NewHistogramer(quantiles []QuantileInfo) *Histogramer {
	return &Histogramer{
		Quantiles: quantiles,
	}
}

// Add adds a blood sample to the summary
func (s *Histogramer) Add(value float64) {
	for i, quantile := range s.Quantiles {
		if value < quantile.Threshold {
			s.Quantiles[i].Count++
			s.Quantiles[i].Sum += float64(value)
		}
	}
	s.Count++
	s.Sum += float64(value)
}

// Remove adds a blood sample to the summary
func (s *Histogramer) Remove(value float64) {
	for i, quantile := range s.Quantiles {
		if value < quantile.Threshold {
			s.Quantiles[i].Count--
			s.Quantiles[i].Sum -= float64(value)
		}
	}
	s.Count--
	s.Sum -= float64(value)
}

// QuantileReport returns the summary report
func (s *Histogramer) QuantileReport() []QuantileInfo {
	return s.Quantiles
}
