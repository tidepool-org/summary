package summarizer

import (
	"log"
	"reflect"
	"time"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
)

const (
	//Layout is how time is represented
	Layout = "2006-01-02T15:04:05Z"
)

//ActivitySummarizer accumulates data on which devices and clients reported activity
type ActivitySummarizer struct {
	Usage       []api.UploadActivity
	ActivityMap map[string]int
}

//NewActivitySummarizer  constructor for ActivitySummarizer
func NewActivitySummarizer() *ActivitySummarizer {
	return &ActivitySummarizer{
		Usage:       make([]api.UploadActivity, 0),
		ActivityMap: make(map[string]int),
	}
}

//ProcessBG updates the time
func (a *ActivitySummarizer) ProcessBG(bg *data.Blood) {
	if bg.UploadID != nil {
		offset := a.ActivityMap[*bg.UploadID]

		var uploadTime time.Time
		var err error
		if bg.Time != nil {
			uploadTime, err = time.Parse(Layout, *bg.Time)
			if err != nil {
				log.Printf("cannot parse time %v", bg.Time)
				return
			}
		}

		if bg.Time != nil && a.Usage[offset].Event.Time.Before(uploadTime) {
			a.Usage[offset].Event.Time = uploadTime
		}
	}
}

//ProcessUpload the device and client used in the upload
func (a *ActivitySummarizer) ProcessUpload(upload *data.Upload) {
	device := api.Device{
		DeviceManufacturers: upload.DeviceManufacturers,
		DeviceModel:         upload.DeviceModel,
		DeviceSerialNumber:  upload.DeviceSerialNumber,
	}

	var client api.Client
	if upload.Client != nil {
		client = api.Client{
			Name:     upload.Client.Name,
			Platform: upload.Client.Platform,
			Version:  upload.Client.Version,
		}
	}

	found := false
	for _, u := range a.Usage {
		if reflect.DeepEqual(u.Device, &device) &&
			reflect.DeepEqual(u.Client, &client) &&
			u.Event.Type == upload.Type {
			found = true
			break
		}
	}
	if !found {
		if upload.ID != nil {
			a.ActivityMap[*upload.ID] = len(a.Usage)
		}
		a.Usage = append(a.Usage,
			api.UploadActivity{
				Client: &client,
				Device: &device,
				Event:  &api.UpdateEvent{Type: upload.Type},
			})
	}
}
