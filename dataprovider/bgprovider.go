package dataprovider

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	//Layout is how time is represented in the API
	Layout = "2006-01-02T15:04:05Z"
)

//BG is a data record, usual cbg, bg, or upload
type BG interface{}

//BGProvider provides a sequence of blood glucose readings
type BGProvider interface {
	Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BG, users []string)
}

// MockProvider provides a static sequence of BG values
type MockProvider struct {
	NumUsers int
}

var _ BGProvider = &MockProvider{}

// PStr moves a value from the stack to the heap and returns pointer to it
func PStr(x string) *string {
	return &x
}

// PF64 moves a value from the stack to the heap and returns pointer to it
func PF64(x float64) *float64 {
	return &x
}

//Get provide blood glucose and upload values on a channel, close channel when no more values
// provide uploads BEFORE blood glucose that refers to them
func (b *MockProvider) Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BG, users []string) {

	duration := int(to.Sub(from).Minutes())

	for j := 0; j < len(users); j++ {

		ch <- data.Upload{
			Base: data.Base{
				Active:   true,
				DeviceID: PStr(fmt.Sprintf("device-for-user-%d", j)),
				Time:     PStr("2020-08-18T08:29:02Z"),
				Type:     "upload",
				UploadID: PStr("xyz"),
				UserID:   &users[j],
			},
			Client: &api.Client{
				Name:     PStr("Tidepool Mobile 99.3"),
				Platform: PStr("windows"),
			},
			Device: api.Device{
				DeviceManufacturers: &[]string{"dexcom"},
				DeviceModel:         PStr("G6"),
				DeviceSerialNumber:  PStr("0xfeedbeef"),
			},
		}

		for i := 0; i < 1000; i++ {
			t := rand.Intn(duration)
			d := from.Add(time.Duration(t) * time.Minute)
			bg := rand.Float64()*250.0 + 30.0
			ch <- data.Blood{
				Base: data.Base{
					Active:   true,
					DeviceID: PStr("foo"),
					Time:     PStr(d.Format(Layout)),
					Type:     "cbg",
					UploadID: PStr("xyz"),
					UserID:   &users[j],
				},
				Units: PStr("mg/dl"),
				Value: PF64(bg),
			}
		}
	}

	close(ch)
}

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

//GetDeviceData sends device data for given userIds over given time period to given channel
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
