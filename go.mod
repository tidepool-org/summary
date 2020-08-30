module github.com/tidepool-org/summary

replace github.com/tidepool-org/tdigest => /Users/derrickburns/go/src/github.com/tidepool-org/tdigest

require (
	github.com/confluentinc/confluent-kafka-go v1.4.2
	github.com/deepmap/oapi-codegen v1.3.12
	github.com/getkin/kin-openapi v0.20.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo/v4 v4.1.17
	github.com/tidepool-org/tdigest v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.4.0
	go.uber.org/fx v1.13.1
)

go 1.15
