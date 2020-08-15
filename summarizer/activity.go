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
	Usage []api.UploadActivity
}

//Process the device and client used in the upload
func (a *ActivitySummarizer) Process(upload *data.Upload) {
	device := api.Device{
		DeviceManufacturers: upload.DeviceManufacturers,
		DeviceModel:         upload.DeviceModel,
		DeviceSerialNumber:  upload.DeviceSerialNumber,
	}

	uploadTime, err := time.Parse(Layout, *upload.Time)
	if err != nil {
		log.Printf("cannot parse time %v", upload.Time)
		return
	}

	found := false
	for _, u := range a.Usage {
		if reflect.DeepEqual(u.Device, device) &&
			reflect.DeepEqual(u.Client, upload.Client) &&
			u.Event.Type == upload.Type {
			if u.Event.Time.Before(uploadTime) {
				u.Event.Time = uploadTime
				found = true
				break
			}
		}
	}
	if !found {
		a.Usage = append(a.Usage,
			api.UploadActivity{
				upload.Client,
				&device,
				&api.UpdateEvent{Time: uploadTime, Type: upload.Type},
			})
	}
}
