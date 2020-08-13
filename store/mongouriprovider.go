package store

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// MongoURIProvider parameters for Mongo database
type MongoURIProvider struct {
	Scheme    string `envconfig:"TIDEPOOL_STORE_SCHEME" default:"mongodb"`
	Hosts     string `envconfig:"TIDEPOOL_STORE_ADDRESSES" required:"true"`
	User      string `envconfig:"TIDEPOOL_STORE_USERNAME" required:"true"`
	Password  string `envconfig:"TIDEPOOL_STORE_PASSWORD" required:"true"`
	OptParams string `envconfig:"TIDEPOOL_STORE_OPT_PARAMS" default:""`
	Ssl       string `envconfig:"TIDEPOOL_STORE_TLS" default:"true"`
}

// URIProvider provides a URI
type URIProvider interface {
	URI() string
}

// NewMongoURIProviderFromEnv creates a URI provider from environment variables
func NewMongoURIProviderFromEnv() MongoURIProvider {
	var mongo MongoURIProvider
	err := envconfig.Process("", &mongo)
	if err != nil {
		log.Fatal(err.Error())
	}
	return mongo
}

var _ URIProvider = MongoURIProvider{}

//URI provide URI to Mongo
func (m MongoURIProvider) URI() string {

	var cs string
	if m.Scheme != "" {
		cs = m.Scheme + "://"
	} else {
		cs = "mongodb://"
	}

	if m.User != "" {
		cs += m.User
		if m.Password != "" {
			cs += ":"
			cs += m.Password
		}
		cs += "@"
	}

	if m.Hosts != "" {
		cs += m.Hosts
		cs += "/"
	} else {
		cs += "localhost/"
	}

	if m.Ssl == "true" {
		cs += "?ssl=true"
	} else {
		cs += "?ssl=false"
	}

	if m.OptParams != "" {
		cs += "&"
		cs += m.OptParams
	}
	return cs
}
