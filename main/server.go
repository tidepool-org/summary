package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/bgprovider"
	"github.com/tidepool-org/summary/summarizer"
)

//SummaryServer provides summaries as a service
type SummaryServer struct {
	Provider bgprovider.BGProvider
}

var _ api.ServerInterface = &SummaryServer{} // confirms that interface is implemented

// PostV1ClinicsCliniidSummary provides summaries for clinicians
// (POST /v1/clinics/{clinicid}/summaries)
func (c *SummaryServer) PostV1ClinicsCliniidSummary(ctx echo.Context, clinicid string) error {
	var summaryRequest api.SummaryRequest

	if err := ctx.Bind(&summaryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error parsing parameters")
	}

	summarizer := summarizer.NewSummarizer(summaryRequest)
	from, to := DateRange(req api.SummaryRequest)
	ch := make(chan bgprovider.BG)

	s.provider.Get(ctx.Request().Context(), from, to, ch, false)

	for {
		select {
		case <-ctx.Request().Context().Done():
			return ctx.JSON(http.StatusRequestTimeout, nil)
		case bg, ok := <-ch:
			if !ok {
				summary := summarizer.Summary()
				return ctx.JSON(http.StatusOK, &summary)
			}
			summarizer.Process(bg)
		}
	}
}

//DateRange provide the times needed to produce the reports
func DateRange(req api.SummaryRequest) (from, to time.Time) {
	var numDays int
	if summaryRequest.Period.Length == "day" {
		numDays = summaryRequest.Period.NumPeriods
	} else if summaryRequest.Period.Length == "week" {
		numDays = 7 * summaryRequest.Period.NumPeriods
	}
	to := time.Now()
	from := to.AddDate(0, 0, -numDays)
	return from, to
}

// PostV1UsersUseridSummary provides summaries for a given user
// (POST /v1/users/{userid}/summary)
func (c *SummaryServer) PostV1UsersUseridSummary(ctx echo.Context, userid string) error {
	var summaryRequest api.SummaryRequest

	if err := ctx.Bind(&summaryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error parsing parameters")
	}

	summaryResponse := api.SummaryResponse{}
	return ctx.JSON(http.StatusOK, &summaryResponse)
}
