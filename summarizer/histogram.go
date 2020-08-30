package summarizer

import (
	"log"
	"math"

	"github.com/tidepool-org/tdigest"
)

// QuantileInfo aggregates partial info leading to a quantile summary
type QuantileInfo struct {
	Count     float64
	Mean      float64
	Threshold float64
	Name      string
}

// Histogramer create a summary of blood glucose data
type Histogramer struct {
	Info   []QuantileInfo
	Digest *tdigest.TDigest
}

//Bins return the number of bins
func (s *Histogramer) Bins() int {
	return len(s.Info) - 1
}

//OverallCount total number of samples received
func (s *Histogramer) OverallCount() float64 {
	return s.Count(s.Bins())
}

//OverallMean is the overall mean
func (s *Histogramer) OverallMean() float64 {
	return s.Mean(s.Bins())
}

//Percentage return the percentage of items in bin
func (s *Histogramer) Percentage(index int) float64 {
	total := s.OverallCount()
	if total == 0.0 {
		return 0.0
	}
	result := s.Info[index].Count / total
	digestResult := s.Digest.CDF(s.Info[index].Threshold)
	absErr := result - digestResult
	errPts := int(100 * math.Abs(absErr))
	if errPts > 0 {
		log.Printf("histogram percentage %0.3f err %0.4f count %v totalCount %v", result, absErr, s.Info[index].Count, total)
	}
	return result
}

//Count return count for given bin
func (s *Histogramer) Count(index int) float64 {
	return s.Info[index].Count
}

//Mean return mean for given bin
func (s *Histogramer) Mean(index int) float64 {
	return s.Info[index].Mean
}

//Threshold return threshold for the bin
func (s *Histogramer) Threshold(index int) float64 {
	return s.Info[index].Threshold
}

//Name return name for the bin
func (s *Histogramer) Name(index int) string {
	return s.Info[index].Name
}

// NewHistogramer create a Histogramer
func NewHistogramer(src []QuantileInfo) *Histogramer {
	dst := make([]QuantileInfo, len(src)+1)
	copy(dst, src)
	dst[len(src)].Threshold = math.MaxFloat64

	return &Histogramer{
		Info:   dst,
		Digest: tdigest.NewWithCompression(100.0),
	}
}

// Add adds an EGV to the summary
func (s *Histogramer) Add(value float64) {
	s.Digest.Add(value, 1.0)
	for i, quantile := range s.Info {
		if value < quantile.Threshold {
			s.Info[i].Count++
			s.Info[i].Mean = (s.Info[i].Mean)*((s.Info[i].Count-1.0)/s.Info[i].Count) + (value / s.Info[i].Count)
		}
	}
}
