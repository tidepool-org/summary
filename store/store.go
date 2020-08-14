package store

import (
	// Built-in Golang packages
	"context" // manage multiple requests
	"errors"
	"log"     // os.Exit(1) on Error
	"reflect" // get an object type

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// Official 'mongo-go-driver' packages
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	//DefaultPagingParams default paging parameters
	DefaultPagingParams = MongoPagingParams{Offset: 0, Limit: 10}
)

// StorageInterface is the interface to the persistent storage
type StorageInterface interface {
	InsertOne(ctx context.Context, collection string, document interface{}) (*string, error)
	FindOne(ctx context.Context, collection string, filter interface{}, data interface{}) error
	Find(ctx context.Context, collection string, filter interface{}, pagingParams *MongoPagingParams, data interface{}) error
	UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) error
	Update(ctx context.Context, collection string, filter interface{}, update interface{}) error
	Aggregate(ctx context.Context, collection string, pipeline []bson.D, data interface{}) error
}

//MongoStoreClient is the Mongo client
type MongoStoreClient struct {
	Client       *mongo.Client
	DatabaseName string
}

// MongoPagingParams are the paging parameters
type MongoPagingParams struct {
	Offset int64
	Limit  int64
}

// NewMongoStoreClient creates a client to a mongo store
func NewMongoStoreClient(provider URIProvider, databaseName string) *MongoStoreClient {

	log.Println("NewMongoStoreClient: Creating Mongo Store")
	client, err := mongo.NewClient(options.Client().ApplyURI(provider.URI()))
	if err != nil {
		log.Fatalln("NewMongoStoreClient: cannot create client:", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatalln("NewMongoStoreClient: cannot connect to client:", err)
	}
	log.Println("NewMongoStoreClient: Created Mongo Store Successfully")

	return &MongoStoreClient{
		Client:       client,
		DatabaseName: databaseName,
	}
}

//Ping checks to see if the client is alive
func (d MongoStoreClient) Ping() error {
	return d.Client.Ping(context.Background(), nil)
}

//InsertOne inserts a record into a collection, return pointer to database id
func (d MongoStoreClient) InsertOne(ctx context.Context, collection string, document interface{}) (*string, error) {
	// InsertOne() method Returns mongo.InsertOneResult
	// Access a MongoDB collection through a database
	col := d.Client.Database(d.DatabaseName).Collection(collection)

	result, insertErr := col.InsertOne(ctx, document)
	if insertErr != nil {
		log.Println("InsertOne ERROR:", insertErr)
		return nil, insertErr
	}

	log.Println("InsertOne: result type: ", reflect.TypeOf(result))
	log.Println("InsertOne: API result:", result)

	// get the inserted ID string
	newID := result.InsertedID
	log.Println("InsertOne: newID:", newID)
	log.Println("InsertOne: newID type:", reflect.TypeOf(newID))

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		newID := oid.Hex()
		log.Printf("Inserted: %v\n", newID)
		return &newID, nil
	}
	return nil, errors.New("can not decode database ID")
}

//FindOne finds a record in colletion
func (d MongoStoreClient) FindOne(ctx context.Context, collection string, filter interface{}, data interface{}) error {
	col := d.Client.Database(d.DatabaseName).Collection(collection)

	if err := col.FindOne(ctx, filter).Decode(data); err != nil {
		log.Println("FindOne: error ", err)
		return err
	}
	log.Printf("FindOne: Found: %v\n", data)
	return nil
}

//Find finds records matching filter criteria in a collection
func (d MongoStoreClient) Find(ctx context.Context, collection string, filter interface{}, pagingParams *MongoPagingParams, data interface{}) error {
	findOptions := options.Find()
	findOptions.SetLimit(pagingParams.Limit)
	findOptions.SetSkip(pagingParams.Offset)

	if pagingParams == nil {
		pagingParams = &DefaultPagingParams
	}
	log.Println("Find: print options: ", *findOptions.Limit, *findOptions.Skip)
	log.Println("Find: filter: ", filter)

	col := d.Client.Database(d.DatabaseName).Collection(collection)

	cursor, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return err
	}
	if err = cursor.All(ctx, data); err != nil {
		return err
	}

	return nil
}

//Update updates
func (d MongoStoreClient) Update(ctx context.Context, collection string, filter interface{}, update interface{}) error {
	col := d.Client.Database(d.DatabaseName).Collection(collection)

	_, err := col.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println("Update: error on update many", err)
		return err
	}
	return nil
}

//UpdateOne updates a single record
func (d MongoStoreClient) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}) error {
	log.Println("UpdateOne")

	col := d.Client.Database(d.DatabaseName).Collection(collection)

	_, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("UpdateOne: error on update one", err)
		return err
	}
	log.Println("Updated")
	return nil
}

//Aggregate runs an aggregation
func (d MongoStoreClient) Aggregate(ctx context.Context, collection string, pipeline []bson.D, data interface{}) error {
	col := d.Client.Database(d.DatabaseName).Collection(collection)

	cursor, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	if err = cursor.All(ctx, data); err != nil {
		return err
	}
	log.Println("Aggregate:", data)
	return nil
}
