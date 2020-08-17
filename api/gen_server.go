// Package Summary provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package api

import (
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Retrieve summaries for all patients of a clinic
	// (POST /v1/clinics/{clinicid}/summaries)
	PostV1ClinicsCliniidSummary(ctx echo.Context, clinicid string) error

	// (POST /v1/users/{userid}/summary)
	PostV1UsersUseridSummary(ctx echo.Context, userid string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostV1ClinicsCliniidSummary converts echo context to params.
func (w *ServerInterfaceWrapper) PostV1ClinicsCliniidSummary(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "clinicid" -------------
	var clinicid string

	err = runtime.BindStyledParameter("simple", false, "clinicid", ctx.Param("clinicid"), &clinicid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter clinicid: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostV1ClinicsCliniidSummary(ctx, clinicid)
	return err
}

// PostV1UsersUseridSummary converts echo context to params.
func (w *ServerInterfaceWrapper) PostV1UsersUseridSummary(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userid" -------------
	var userid string

	err = runtime.BindStyledParameter("simple", false, "userid", ctx.Param("userid"), &userid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userid: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostV1UsersUseridSummary(ctx, userid)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST("/v1/clinics/:clinicid/summaries", wrapper.PostV1ClinicsCliniidSummary)
	router.POST("/v1/users/:userid/summary", wrapper.PostV1UsersUseridSummary)

}

