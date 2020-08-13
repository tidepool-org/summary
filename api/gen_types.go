// Package Summary provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"time"
)

// Quantile defines model for Quantile.
type Quantile struct {

	// The number of values below the threshold value.
	Count *int `json:"count,omitempty"`

	// The name of the quantile.
	Name string `json:"name"`

	// The percentage of samples below the threshold value.
	Percentage float32 `json:"percentage"`

	// The threshold of the quantile.
	Threshold float32 `json:"threshold"`
}

// SummaryPeriod defines model for SummaryPeriod.
type SummaryPeriod struct {

	// The length of the period.
	Length string `json:"length"`

	// The start date of the period being reported.
	Start time.Time `json:"start"`

	// The time that these statistics were last updated for the given period.
	Updated time.Time `json:"updated"`
}

// SummaryReport defines model for SummaryReport.
type SummaryReport struct {

	// Summary period
	Period SummaryPeriod `json:"period"`

	// Summary diabetes statistics
	Stats SummaryStatistics `json:"stats"`
}

// SummaryRequest defines model for SummaryRequest.
type SummaryRequest struct {
	Period struct {
		Length     string `json:"length"`
		NumPeriods int    `json:"numPeriods"`
	} `json:"period"`
	Quantiles []struct {

		// The name of the quantile.
		Name string `json:"name"`

		// The threshold value for the quantiles.  All samples must be below the threshold to be included in the count.
		Threshold float32 `json:"threshold"`
	} `json:"quantiles"`
	Units SummaryRequestUnits `json:"units"`
}

// SummaryRequestUnits defines model for SummaryRequest-units.
type SummaryRequestUnits string

// List of SummaryRequestUnits
const (
	SummaryRequestUnits_mg_dL  SummaryRequestUnits = "mg/dL"
	SummaryRequestUnits_mg_dl  SummaryRequestUnits = "mg/dl"
	SummaryRequestUnits_mmol_L SummaryRequestUnits = "mmol/L"
	SummaryRequestUnits_mmol_l SummaryRequestUnits = "mmol/l"
)

// SummaryResponse defines model for SummaryResponse.
type SummaryResponse struct {

	// A summary of which devices were used and when to upload diabetes data
	Activity *[]struct {

		// The client software that provided diabetes data
		Client *struct {

			// The name of the client software used to extract the data
			Name *string `json:"name,omitempty"`

			// The software platform on which the client software was run
			Platform *string                 `json:"platform,omitempty"`
			Private  *map[string]interface{} `json:"private,omitempty"`

			// The version of the client software used to extract the data
			Version *string `json:"version,omitempty"`
		} `json:"client,omitempty"`

		// The device used to provide diabetes data
		Device *struct {

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
		} `json:"device,omitempty"`

		// The time that that the device was last used to provide diabetes data
		Event *UpdateEvent `json:"event,omitempty"`
	} `json:"activity,omitempty"`

	// Summary of recent glucose information.
	Reports []SummaryReport `json:"reports"`

	// String representation of a Tidepool User ID
	Userid UserId `json:"userid"`
}

// SummaryStatistics defines model for SummaryStatistics.
type SummaryStatistics struct {

	// Total number of samples in period.
	Count int `json:"count"`

	// Mean glucose over samples in period
	Mean float32 `json:"mean"`

	// An array of quantile measurements.
	Quantiles []Quantile          `json:"quantiles"`
	Units     SummaryRequestUnits `json:"units"`
}

// UpdateEvent defines model for UpdateEvent.
type UpdateEvent struct {

	// The time of the most recent upload.
	Time time.Time `json:"time"`

	// The data type that was uploaded.
	Type string `json:"type"`
}

// UserId defines model for userId.
type UserId string

// PostV1ClinicsCliniidSummaryJSONBody defines parameters for PostV1ClinicsCliniidSummary.
type PostV1ClinicsCliniidSummaryJSONBody SummaryRequest

// PostV1UsersUseridSummaryJSONBody defines parameters for PostV1UsersUseridSummary.
type PostV1UsersUseridSummaryJSONBody SummaryRequest

// PostV1ClinicsCliniidSummaryRequestBody defines body for PostV1ClinicsCliniidSummary for application/json ContentType.
type PostV1ClinicsCliniidSummaryJSONRequestBody PostV1ClinicsCliniidSummaryJSONBody

// PostV1UsersUseridSummaryRequestBody defines body for PostV1UsersUseridSummary for application/json ContentType.
type PostV1UsersUseridSummaryJSONRequestBody PostV1UsersUseridSummaryJSONBody

