package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/dataprovider"
	"github.com/tidepool-org/summary/summarizer"
)

//SummaryServer provides summaries as a service
type SummaryServer struct {
	Provider      dataprovider.BGProvider
	ShareProvider dataprovider.ShareProvider
}

// NewSummaryServer Create a new summary service
func NewSummaryServer(provider dataprovider.BGProvider, shareProvider dataprovider.ShareProvider) *SummaryServer {
	return &SummaryServer{
		Provider:      provider,
		ShareProvider: shareProvider,
	}
}

var _ api.ServerInterface = &SummaryServer{} // confirms that interface is implemented

// PostV1UsersUseridSummaries provides summaries for a given user
// (POST /v1/users/{userid}/summaries)
func (c *SummaryServer) PostV1UsersUseridSummaries(ctx echo.Context, userID string) error {
	var summaryRequest api.SummaryRequest

	if err := ctx.Bind(&summaryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error parsing parameters")
	}

	summarizer := summarizer.NewSummarizer(summaryRequest)
	from, to := DateRange(summaryRequest)

	userids, err := c.ShareProvider.SharerIdsForUser(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	ch := make(chan dataprovider.BG)
	go c.Provider.Get(ctx.Request().Context(), from, to, ch, userids)

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

// PostV1ClinicsCliniidSummary provides summaries for clinicians
// (POST /v1/clinics/{clinicid}/summaries)
func (c *SummaryServer) PostV1ClinicsCliniidSummary(ctx echo.Context, clinicID string) error {
	var summaryRequest api.SummaryRequest

	if err := ctx.Bind(&summaryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error parsing parameters")
	}

	summarizer := summarizer.NewSummarizer(summaryRequest)
	from, to := DateRange(summaryRequest)

	userids, err := c.ShareProvider.SharerIdsForClinic(ctx.Request().Context(), clinicID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}
	ch := make(chan dataprovider.BG)
	go c.Provider.Get(ctx.Request().Context(), from, to, ch, userids)

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
//TODO round ranges
func DateRange(req api.SummaryRequest) (from, to time.Time) {
	var numDays int
	if req.Period.Length == "day" {
		numDays = req.Period.NumPeriods
	} else if req.Period.Length == "week" {
		numDays = 7 * req.Period.NumPeriods
	}
	to = time.Now()
	from = to.AddDate(0, 0, -numDays)
	return
}

// PostV1UsersUseridSummary provides summaries for a given user
// (POST /v1/users/{userid}/summary)
func (c *SummaryServer) PostV1UsersUseridSummary(ctx echo.Context, userID string) error {
	var summaryRequest api.SummaryRequest
	if err := ctx.Bind(&summaryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error parsing parameters")
	}

	summarizer := summarizer.NewSummarizer(summaryRequest)
	from, to := DateRange(summaryRequest)

	userids := []string{userID}
	ch := make(chan dataprovider.BG)
	go c.Provider.Get(ctx.Request().Context(), from, to, ch, userids)

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
