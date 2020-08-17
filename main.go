package main

import (
	"context"
	"log"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/bgprovider"
	"github.com/tidepool-org/summary/server"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	// ServerTimeoutAmount is the amount of time before we time out the server
	ServerTimeoutAmount = 20
	_                   = openapi3filter.Options{}
)

//ServiceConfig the configuration for the summary service
type ServiceConfig struct {
	ServiceAuth string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SECRET" required:"false"`
	Address     string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SERVER_ADDRESS" default:":8080"`
}

//NewServiceConfigFromEnv creates a service config
func NewServiceConfigFromEnv() *ServiceConfig {
	var config ServiceConfig
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &config
}

//main is the main loop
func main() {

	config := NewServiceConfigFromEnv()

	// Echo instance
	e := echo.New()
	e.Logger.Print("Starting Main Loop")
	swagger, err := api.GetSwagger()
	if err != nil {
		e.Logger.Fatal("Cound not get spec")
	}

	// Middleware
	//authClient := AuthClient{store: dbstore}
	//filterOptions := openapi3filter.Options{AuthenticationFunc: authClient.AuthenticationFunc}
	//options := Options{Options: filterOptions}
	options := api.Options{}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(api.OapiRequestValidator(swagger, &options))

	// Routes
	e.GET("/status", hello)

	data, err := json.MarshalIndent(e.Routes(), "", "  ")
        if err != nil {
	        e.Logger.Printf("cannot list routes")
        } else {
		e.Logger.Printf("routes %v", data)
	}

	// Register Handler
	api.RegisterHandlers(e, &server.SummaryServer{
		Provider: &bgprovider.MockProvider{},
	})

	// Start server
	e.Logger.Printf("Starting Server at: %s\n", config.Address)
	go func() {
		if err := e.Start(config.Address); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ServerTimeoutAmount)*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
