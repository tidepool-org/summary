package kafkasource

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/kelseyhightower/envconfig"
	"github.com/tidepool-org/summary/debezium"
)

type EventSource interface {
	Run(ctx context.Context, ch chan<- *debezium.MongoDBEvent)
}

// KafkaSource defines a kafka source
type KafkaSource struct {
	Consumer *kafka.Consumer
	Topic    string
}

//Config is the input configuration
type Config struct {
	Brokers     string `envconfig:"TIDEPOOL_KAFKA_BROKERS" required:"true"`
	Prefix      string `envconfig:"TIDEPOOL_KAFKA_TOPIC_PREFIX" default:""`
	Topic       string `envconfig:"TIDEPOOL_KAFKA_DEVICEDATA_TOPIC" required:"true"`
	ServiceAuth string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SECRET" required:"true"`
	Address     string `envconfig:"TIDEPOOL_SUMMARY_SERVICE_SERVER_ADDRESS" default:":8080"`
}

//NewKafkaConfigFromEnv create a Kafka config
func NewKafkaConfigFromEnv() *Config {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &config
}

var _ EventSource = &KafkaSource{} // ensures that contract is met

//NewKafkaSource creates a new kafka message source
func NewKafkaSource(config Config) (*KafkaSource, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Brokers,
		"group.id":          "summary",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	return &KafkaSource{
		Consumer: c,
		Topic:    config.Prefix + config.Topic,
	}, nil
}

// Run extracts messages from Kafka topic, deserialized them, and passes them to the given channel
func (s *KafkaSource) Run(ctx context.Context, ch chan<- *debezium.MongoDBEvent) {

	s.Consumer.SubscribeTopics([]string{s.Topic}, nil)

	run := true

	for run == true {
		select {
		case <-ctx.Done():
			close(ch)
			run = false
		default:
			ev := s.Consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
				var rec debezium.MongoDBEvent
				if err := json.Unmarshal(e.Value, &rec); err != nil {
					fmt.Println(s.Topic, "Error Unmarshalling", err)
				} else {
					ch <- &rec
				}

			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}
