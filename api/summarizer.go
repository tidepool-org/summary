package api

import (
	"log"
	"time"

	"github.com/tidepool-org/summary/bgprovider"
)

//Summarizer creates summaries
type Summarizer struct {
	Histograms []*Histogramer // histograms for each time range
	Periods    []SummaryPeriod
	Normalizer UnitNormalizer
	Request    SummaryRequest
}

// NewSummarizer creates a Summarizer for the given request
func NewSummarizer(request SummaryRequest) *Summarizer {
	histograms := make([]*Histogramer, request.Period.NumPeriods)
	quantiles := make([]QuantileInfo, len(request.Quantiles))

	for i := range quantiles {
		quantiles[i].Name = request.Quantiles[i].Name
		quantiles[i].Threshold = float64(request.Quantiles[i].Threshold)
	}

	periods := make([]SummaryPeriod, request.Period.NumPeriods)

	now := time.Now()
	ending := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
	var duration time.Duration
	switch request.Period.Length {
	case "day":
		duration = 24 * time.Hour
	case "week":
		duration = 7 * 24 * time.Hour
	}

	for i := range histograms {
		histograms[i] = NewHistogramer(quantiles)
		periods[i].End = ending
		ending = ending.Add(-1 * duration)
		periods[i].Start = ending
		periods[i].Length = request.Period.Length
	}

	return &Summarizer{
		Histograms: histograms,
		Request:    request,
		Normalizer: &BloodGlucoseNormalizer{},
		Periods:    periods,
	}
}

//Process an event
func (s *Summarizer) Process(rec *bgprovider.BG) {

	now := time.Now()
	blood := rec.Blood
	if blood.Type == "cbg" || blood.Type == "smbg" {

		if blood.Value == nil || blood.Units == nil {
			log.Printf("skipping entry with missing value or units %v\n", blood)
			return
		}
		layout := "2006-01-02T15:04:05Z"
		t, err := time.Parse(layout, *blood.Time)

		if err != nil {
			log.Printf("skipping entry with bad date %v\n", blood)
			return
		}

		standardized := s.Normalizer.ToStandard(float32(*blood.Value), *blood.Units)

		if blood.Active {
			for i, p := range s.Periods {
				if (!t.Before(p.Start)) && p.End.After(t) {
					s.Histograms[i].Add(float64(standardized))
					p.Updated = now
				}
			}
		}
	} 

	upload := rec.Upload.
	
		log.Printf("skipping type %v\n", blood.Type)
	}
}



type Device struct {

	// An array of string tags indicating the manufacturer(s) of the device.
	//
	// In order to avoid confusion resulting from referring to a single manufacturer with more than one name—for example, using both 'Minimed' and 'Medtronic' interchangeably—we restrict the set of strings used to refer to manufacturers to the set listed above and enforce *exact* string matches (including casing).
	//
	// `deviceManufacturers` is an array of one or more string "tags" because there are devices resulting from a collaboration between more than one manufacturer, such as the Tandem G4 insulin pump with CGM integration (a collaboration between `Tandem` and `Dexcom`).
	DeviceManufacturers *[]string `json:"deviceManufacturers,omitempty"`

	// A string identifying the model of the device.
	//
	// The `deviceModel` is a non-empty string that encodes the model of device. We endeavor to match each manufacturer's standard for how they represent model name in terms of casing, whether parts of the name are represented as one word or two, etc.
	DeviceModel *string `json:"deviceModel,omitempty"`

	// A string encoding the device's serial number.
	//
	// The `deviceSerialNumber` is a string that encodes the serial number of the device. Note that even if a manufacturer only uses digits in its serial numbers, the SN should be stored as a string regardless.
	//
	// Uniquely of string fields in the Tidepool device data models, `deviceSerialNumber` *may* be an empty string. This is essentially a compromise: having the device serial number is extremely important (especially for e.g., clinical studies) but in 2016 we came across our first case where we *cannot* recover the serial number of the device that generated the data: Dexcom G5 data uploaded to Tidepool through Apple iOS's HealthKit integration.
	DeviceSerialNumber *string `json:"deviceSerialNumber,omitempty"`

	// An array of string tags indicating the function(s) of the device.
	//
	// The `deviceTags` array should be fairly self-explanatory as an array of tags indicating the function(s) of a particular device. For example, the Insulet OmniPod insulin delivery system has the tags `bgm` and `insulin-pump` since the PDM is both an insulin pump controller and includes a built-in blood glucose monitor.
	DeviceTags *[]interface{} `json:"deviceTags,omitempty"`
}

/*
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
*/
