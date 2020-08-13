package api

import (
	"context"
	"log"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

//Config is the input configuration
type Config struct {
	Brokers     string `envconfig:"TIDEPOOL_KAFKA_BROKERS" required:"true"`
	Prefix      string `envconfig:"TIDEPOOL_KAFKA_TOPIC_PREFIX" default:""`
	Topic       string `envconfig:"TIDEPOOL_KAFKA_DEVICEDATA_TOPIC" required:"true"`
	ServiceAuth string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SECRET" required:"true"`
	Address     string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SERVER_ADDRESS" default:":8080"`
}

//MainLoop is the main loop
func MainLoop() {

	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Echo instance
	e := echo.New()
	e.Logger.Print("Starting Main Loop")
	swagger, err := GetSwagger()
	if err != nil {
		e.Logger.Fatal("Cound not get spec")
	}

	// Middleware
	//authClient := AuthClient{store: dbstore}
	//filterOptions := openapi3filter.Options{AuthenticationFunc: authClient.AuthenticationFunc}
	//options := Options{Options: filterOptions}
	options := Options{}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(OapiRequestValidator(swagger, &options))

	// Routes
	e.GET("/", hello)

	// Register Handler
	RegisterHandlers(e, &SummaryServer{})

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
