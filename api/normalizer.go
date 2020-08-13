package api

// ConverstionRatio is the mmol to mg/dL ratio
const ConverstionRatio = 18.0182

// UnitNormalizer converts values to/from standard units
type UnitNormalizer interface {
	ToStandard(value float32, unit string) float32
	FromStandard(value float32, unit string) float32
}

// BloodGlucoseNormalizer converts to/from the internal representatation
type BloodGlucoseNormalizer struct {
}

var _ UnitNormalizer = &BloodGlucoseNormalizer{}

// ToStandard Converts the given value to the standard units
func (*BloodGlucoseNormalizer) ToStandard(value float32, unit string) float32 {
	switch SummaryRequestUnits(unit) {
	case SummaryRequestUnits_mmol_l:
		return value * ConverstionRatio
	case SummaryRequestUnits_mmol_L:
		return value * ConverstionRatio
	case SummaryRequestUnits_mg_dL:
		return value
	case SummaryRequestUnits_mg_dl:
		return value
	default:
		return value
	}
}

// FromStandard Converts the given value from the standard units into the units given
func (*BloodGlucoseNormalizer) FromStandard(value float32, unit string) float32 {
	switch SummaryRequestUnits(unit) {
	case SummaryRequestUnits_mmol_l:
		return value / ConverstionRatio
	case SummaryRequestUnits_mmol_L:
		return value / ConverstionRatio
	case SummaryRequestUnits_mg_dL:
		return value
	case SummaryRequestUnits_mg_dl:
		return value
	default:
		return value
	}
}
