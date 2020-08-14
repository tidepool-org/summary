package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/debezium"
)

//SummaryServer provides summaries as a service
type SummaryServer struct {
	Config Config
}

var _ api.ServerInterface = &SummaryServer{} // confirms that interface is implemented

// PostV1ClinicsCliniidSummary provides summaries for clinicians
// (POST /v1/clinics/{clinicid}/summaries)
func (c *SummaryServer) PostV1ClinicsCliniidSummary(ctx echo.Context, clinicid string) error {
	var summaryRequest SummaryRequest

	if err := ctx.Bind(&summaryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error parsing parameters")
	}

	// create (kafka) source
	// feed source to summarizer
	// when source closes, summarizer sends report on channel
	// await reports from summarizer

	mongoEventCh := make(chan *debezium.MongoDBEvent)
	summaryCh := make(chan SummaryResponse)

	source, _ := NewKafkaSource(c.Config)
	summarizer := NewSummarizer(summaryRequest)

	source.Run(ctx.Request().Context(), mongoEventCh)
	summarizer.Run(ctx.Request().Context(), mongoEventCh, summaryCh)

	for {
		select {
		case <-ctx.Request().Context().Done():
			return ctx.JSON(http.StatusRequestTimeout, nil)
		case s := <-summaryCh:
			return ctx.JSON(http.StatusOK, &s)
		}
	}
}

// PostV1UsersUseridSummary provides summaries for a given user
// (POST /v1/users/{userid}/summary)
func (c *SummaryServer) PostV1UsersUseridSummary(ctx echo.Context, userid string) error {
	var summaryRequest SummaryRequest

	if err := ctx.Bind(&summaryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error parsing parameters")
	}

	summaryResponse := SummaryResponse{}
	return ctx.JSON(http.StatusOK, &summaryResponse)
}
