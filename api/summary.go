package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SummaryServer struct {
}

// (POST /v1/clinics/{clinicid}/summaries)
func (c *SummaryServer) PostV1ClinicsCliniidSummary(ctx echo.Context, clinicid string) error {
	var summaryRequest SummaryRequest

	if err := ctx.Bind(&summaryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error parsing parameters")
	}

	summaryResponse := SummaryResponse{}
	return ctx.JSON(http.StatusOK, &summaryResponse)
}

// (POST /v1/users/{userid}/summary)
func (c *SummaryServer) PostV1UsersUseridSummary(ctx echo.Context, userid string) error {
	var summaryRequest SummaryRequest

	if err := ctx.Bind(&summaryRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "error parsing parameters")
	}

	summaryResponse := SummaryResponse{}
	return ctx.JSON(http.StatusOK, &summaryResponse)
}
