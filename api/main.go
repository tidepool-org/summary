package api


import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Host = "localhost"
	Port = 8080
	ServerString = fmt.Sprintf("%s:%d", Host, Port)
	ServerTimeoutAmount = 20
	_ = openapi3filter.Options{}

)


func MainLoop() {
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
	e.Logger.Printf("Starting Server at: %s\n", ServerString)
	go func() {
		if err := e.Start(ServerString); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ServerTimeoutAmount) * time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}


