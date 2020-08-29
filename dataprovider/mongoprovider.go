package dataprovider

import (
	"context"
	"log"
	"time"

	"github.com/tidepool-org/summary/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoProvider provides individual blood glucose values for a list of userids
type MongoProvider struct {
	Client *mongo.Client
}

var _ BGProvider = &MongoProvider{}

//NewMongoProvider creates a new MongoProvider that uses the given Mongo client
func NewMongoProvider(client *mongo.Client) *MongoProvider {
	return &MongoProvider{
		Client: client,
	}
}

//Get provide blood glucose and upload values on a channel, close channel when no more values
// provide uploads BEFORE blood glucose that refers to them
func (b *MongoProvider) Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BG, users []string) {
	b.GetDeviceData(ctx, from, to, ch, users)
	close(ch)
}

//GetUpload returns the upload record with the given uploadID
func (b *MongoProvider) GetUpload(ctx context.Context, deviceData *mongo.Collection, uploadID string) (*data.Upload, error) {
	singleResult := deviceData.FindOne(ctx,
		bson.M{
			"type":     "upload",
			"uploadId": uploadID,
		})
	var val data.Upload
	if err := singleResult.Decode(&val); err != nil {
		return nil, err
	}
	return &val, nil
}

// GetCGMSettings returns array of CGM settings for all users
func (b *MongoProvider) GetCGMSettings(ctx context.Context, start, end time.Time, userIds []string) error {
	deviceData := b.Client.Database("data").Collection("deviceData")
	endTime := end.Format(time.RFC3339)

	cgmFilter := bson.M{
		"_active": true,
		"_userId": bson.M{"$in": userIds},
		"time":    bson.M{"$lt": endTime},
		"type":    "cgmSettings",
	}

	cgmProjection := new(options.FindOptions).SetProjection(bson.M{
		"_userId":         1,
		"time":            1,
		"uploadId":        1,
		"manufacturerers": 1,
		"serialNumber":    1,
	})

	_, err := deviceData.Find(ctx, cgmFilter, cgmProjection)
	return err
}

//GetDeviceData sends device data for given userIds over given time period to given channel
// When a EGV record is retrieved, we retrieve the corresponding upload record if it has not already been retrieved
// We send the upload record before the EGV record.
func (b *MongoProvider) GetDeviceData(ctx context.Context, start, end time.Time, ch chan<- BG, userIds []string) {

	deviceData := b.Client.Database("data").Collection("deviceData")

	projection := new(options.FindOptions).SetProjection(bson.M{
		"_userId":  1,
		"type":     1,
		"value":    1,
		"units":    1,
		"time":     1,
		"uploadId": 1,
	})

	startTime := start.Format(time.RFC3339)
	endTime := end.Format(time.RFC3339)

	log.Printf("startTime %s", startTime)
	log.Printf("endTime %s", endTime)
	log.Printf("Userids %v", userIds)

	filter := bson.M{
		"_active": true,
		"_userId": bson.M{"$in": userIds},
		"time":    bson.M{"$gte": startTime, "$lt": endTime},
		"type":    bson.M{"$in": []string{"cbg", "smbg"}}}

	log.Printf("filter %v", filter)
	log.Printf("projection %v", projection)

	log.Printf("starting Find of BG")
	cursor, err := deviceData.Find(ctx, filter, projection)
	log.Printf("received cursor of BG")

	if err != nil {
		log.Fatal(err)
	}

	seen := make(map[string]bool)
	count := 0
	for cursor.Next(ctx) {
		var bg data.Blood
		if err := cursor.Decode(&bg); err != nil {
			log.Printf("error decoding bg %v", err)
			continue
		}

		if bg.UploadID != nil {
			uploadID := *bg.UploadID
			if !seen[uploadID] {
				seen[uploadID] = true
				upload, err := b.GetUpload(ctx, deviceData, uploadID)
				if err != nil {
					log.Printf("error decoding upload %v: %v", uploadID, err)
					continue
				}
				ch <- *upload
			}
			ch <- bg
			count++
		}
	}
}
