package summarizer

import "github.com/tidepool-org/summary/api"

// ConverstionRatio is the mmol to mg/dL ratio
const ConverstionRatio = 18.0182

// UnitNormalizer converts values to/from standard units
type UnitNormalizer interface {
	ToStandard(value float32, unit string) float32
	FromStandard(value float32, unit string) float32
}

// GlucoseNormalizer converts to/from the internal representatation
type GlucoseNormalizer struct {
}

var _ UnitNormalizer = &GlucoseNormalizer{}

// ToStandard Converts the given value to the standard units
func (*GlucoseNormalizer) ToStandard(value float32, unit string) float32 {
	switch api.Units(unit) {
	case api.Units_mmol_l:
		return value * ConverstionRatio
	case api.Units_mmol_L:
		return value * ConverstionRatio
	case api.Units_mg_dL:
		return value
	case api.Units_mg_dl:
		return value
	default:
		return value
	}
}

// FromStandard Converts the given value from the standard units into the units given
func (*GlucoseNormalizer) FromStandard(value float32, unit string) float32 {
	switch api.Units(unit) {
	case api.Units_mmol_l:
		return value / ConverstionRatio
	case api.Units_mmol_L:
		return value / ConverstionRatio
	case api.Units_mg_dL:
		return value
	case api.Units_mg_dl:
		return value
	default:
		return value
	}
}
