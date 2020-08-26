package dataprovider

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ShareProvider provides a list of userIds of data storage accounts available to a given user
type ShareProvider interface {
	SharerIdsForClinic(ctx context.Context, clinicID string) ([]string, error)
	SharerIdsForUser(ctx context.Context, userID string) ([]string, error)
}

//MongoShareProvider provide accounts shared person to person
type MongoShareProvider struct {
	Client *mongo.Client
}

var _ ShareProvider = &MongoShareProvider{}

//NewMongoShareProvider creates a new MongoProvider that uses the given Mongo client
func NewMongoShareProvider(client *mongo.Client) *MongoShareProvider {
	return &MongoShareProvider{
		Client: client,
	}
}

//SharerIdsForUser returns the user ids of accounts that are shared with the given user
func (b *MongoShareProvider) SharerIdsForUser(ctx context.Context, userID string) ([]string, error) {
	perms := b.Client.Database("gatekeeper").Collection("perms")

	log.Printf("get distinct")
	sharerIds, err := perms.Distinct(ctx, "sharedId", bson.M{"userId": userID})
	if err != nil {
		log.Fatalf("error getting distinct %v", err)
	}
	log.Printf("get distinct completed")
	type Share struct {
		ID string `bson:"sharerId"`
	}
	log.Printf("num ids %v", len(sharerIds))
	ids := make([]string, len(sharerIds))
	for i, id := range sharerIds {
		log.Printf("i %d share %v", i, id)
		ids[i] = id.(Share).ID
	}
	return ids, err
}

//SharerIdsForClinic returns the user ids of accounts that are shared with the given clinic
func (b *MongoShareProvider) SharerIdsForClinic(ctx context.Context, clinicID string) ([]string, error) {
	perms := b.Client.Database("clinic").Collection("clinicPatients")
	log.Printf("get distinct")
	sharerIds, err := perms.Distinct(ctx, "patientId", bson.M{"clinicId": clinicID})
	if err != nil {
		log.Fatalf("error getting distinct %v", err)
	}
	log.Printf("get distinct completed")
	type Share struct {
		ID string `bson:"patientId"`
	}
	ids := make([]string, len(sharerIds))
	for i, id := range sharerIds {
		log.Printf("i %d share %v", i, id)
		ids[i] = id.(Share).ID
	}
	return ids, err
}
