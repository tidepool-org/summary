package main

import (
	"context"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/dataprovider"
	"github.com/tidepool-org/summary/server"
	"github.com/tidepool-org/summary/store"

	"net/http"
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
func NewServiceConfigFromEnv() (*ServiceConfig, error) {
	var config ServiceConfig
	err := envconfig.Process("", &config)
	return &config, err
}

// NewBGProvider providers a BGProvider
func NewBGProvider(provider *dataprovider.MongoProvider) dataprovider.BGProvider {
	return provider
}

// NewShareProvider providers a ShareProvider
func NewShareProvider(provider *dataprovider.MongoShareProvider) dataprovider.ShareProvider {
	return provider
}

//main is the main loop
func main() {
	fx.New(
		fx.Provide(NewServiceConfigFromEnv),
		fx.Provide(api.GetSwagger),
		fx.Provide(store.NewMongoURIProviderFromEnv),
		fx.Provide(ProvideMongoClient),
		fx.Provide(dataprovider.NewMongoProvider),
		fx.Provide(NewBGProvider),
		fx.Provide(NewShareProvider),
		fx.Provide(dataprovider.NewMongoShareProvider),
		fx.Provide(server.NewSummaryServer),
		fx.Provide(ProvideEchoServer),
		fx.Invoke(invokeHooks),
	).Run()
}

// ProvideMongoClient provides a mongo client
func ProvideMongoClient(uriProvider store.MongoURIProvider) (client *mongo.Client, err error) {
	client, err = mongo.NewClient(mongoOptions.Client().ApplyURI(uriProvider.URI()))
	return
}

//ProvideEchoServer creates an Echo server with a default status endpoint and swagger validation of requests and responses
func ProvideEchoServer(config *ServiceConfig, swagger *openapi3.Swagger) *echo.Echo {
	e := echo.New()

	// Middleware
	//authClient := AuthClient{store: dbstore}
	//filterOptions := openapi3filter.Options{AuthenticationFunc: authClient.AuthenticationFunc}
	//options := Options{Options: filterOptions}

	e.GET("/status", hello)

	loggerConfig := middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/status")
		},
	}
	e.Use(middleware.LoggerWithConfig(loggerConfig))
	e.Use(middleware.Recover())

	options := api.Options{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Path(), "/status")
		},
	}
	e.Use(api.OapiRequestValidatorWithOptions(swagger, &options))
	return e
}

func invokeHooks(lifecycle fx.Lifecycle, e *echo.Echo, config *ServiceConfig, summaryServer *server.SummaryServer, client *mongo.Client) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				// Register Handler
				if err := client.Connect(context.Background()); err != nil {
					return err
				}

				time.Sleep(1 * time.Second)
				if err := client.Ping(context.Background(), nil); err != nil {
					return err
				}

				api.RegisterHandlers(e, summaryServer)

				go func() {
					// Start server
					e.Logger.Printf("Starting Server at: %s\n", config.Address)
					if err := e.Start(config.Address); err != nil {
						e.Logger.Info("shutting down the server")
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				if err := e.Shutdown(ctx); err != nil {
					e.Logger.Fatal(err)
					return err
				}
				return nil
			},
		},
	)
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
