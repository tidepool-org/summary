package summarizer

// QuantileInfo aggregates partial info leading to a quantile summary
type QuantileInfo struct {
	Count     float64
	Mean      float64
	Threshold float64
	Name      string
}

// Histogramer create a summary of blood glucose data
type Histogramer struct {
	Info []QuantileInfo
}

// NewHistogramer create a Histogramer
func NewHistogramer(src []QuantileInfo) *Histogramer {
	dst := make([]QuantileInfo, len(src))
	copy(dst, src)
	return &Histogramer{Info: dst}
}

// Add adds a blood sample to the summary
func (s *Histogramer) Add(value float64) {
	for i, quantile := range s.Info {
		if value < quantile.Threshold {
			s.Info[i].Count++
			s.Info[i].Mean = (s.Info[i].Mean)*((s.Info[i].Count-1.0)/s.Info[i].Count) + (value / s.Info[i].Count)
		}
	}
}

// Remove adds a blood sample to the summary
func (s *Histogramer) Remove(value float64) {
	for i, quantile := range s.Info {
		if value < quantile.Threshold {
			s.Info[i].Count--
			s.Info[i].Mean = (s.Info[i].Mean)*((s.Info[i].Count)/s.Info[i].Count+1) - (value / (s.Info[i].Count + 1))
		}
	}
}
