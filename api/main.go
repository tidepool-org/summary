package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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
	ServerTimeoutAmount = 20
	_ = openapi3filter.Options{}
)

//Config is the input configuration
type Config struct {
	Brokers     string `envconfig:"TIDEPOOL_KAFKA_BROKERS" required:"true"`
	Prefix      string `envconfig:"TIDEPOOL_KAFKA_TOPIC_PREFIX" default:""`
	Topic       string `envconfig:"TIDEPOOL_KAFKA_DEVICEDATA_TOPIC" required:"true"`
	ServiceAuth string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SECRET" required:"true"`
	Address     string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SERVER_ADDRESS" default:":8080"`
}

type Debezium struct {
	After string `json:"after"`
	Before string `json:"before"`
	op string `json:"op"`
}

//Base is a subset of the fields common to all datums
type Base struct {
	Active   bool    `json:"-" bson:"_active"` // if false, this object has been effectively deleted
	DeviceID *string `json:"deviceId,omitempty" bson:"deviceId,omitempty"`
	ID       *string `json:"id,omitempty" bson:"id,omitempty"`
	Source   *string `json:"source,omitempty" bson:"source,omitempty"`
	Time     *string `json:"time,omitempty" bson:"time,omitempty"`
	Type     string  `json:"type,omitempty" bson:"type,omitempty"`
	UploadID *string `json:"uploadId,omitempty" bson:"uploadId,omitempty"`
	UserID   *string `json:"-" bson:"_userId,omitempty"`
}

// Blood is the type of a blood value
type Blood struct {
	Base `bson:",inline"`

	Units *string  `json:"units,omitempty" bson:"units,omitempty"`
	Value *float64 `json:"value,omitempty" bson:"value,omitempty"`
}

// ProcessDeviceDataTopic processs messages for device
func ProcessDeviceDataTopic(config *Config) error {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Brokers,
		"group.id":          "summary",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{config.Prefix + config.Topic}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			var rec map[string]interface{}
			if err := json.Unmarshal(msg.Value, &rec); err != nil {
				fmt.Println(config.Topic, "Error Unmarshalling", err)
			} else {
				//source, source_ok := rec["source"]
				afterField, ok := rec["after"]
				if ok {
					var data Blood
					data_string := fmt.Sprintf("%v", afterField)
					if err := json.Unmarshal([]byte(data_string), &data); err != nil {
						//fmt.Println(topic, "Error Unmarshalling after field", err)
					} else {
					}
				}
			}

		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

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

	go ProcessDeviceDataTopic(&config)

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
