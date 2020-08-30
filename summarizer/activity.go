package summarizer

import (
	"fmt"
	"reflect"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
)

const (
	//Layout is how time is represented
	Layout = "2006-01-02T15:04:05Z"
)

//ActivitySummarizer accumulates data on which devices and clients reported activity
type ActivitySummarizer struct {
	Usage         []api.UploadActivity
	ActivityMap   map[string]int
	DeviceGlucose []*GlucoseSummarizer
	Request       api.SummaryRequest
	Periods       []api.SummaryPeriod
}

//NewActivitySummarizer  constructor for ActivitySummarizer
func NewActivitySummarizer(request api.SummaryRequest, periods []api.SummaryPeriod) *ActivitySummarizer {
	return &ActivitySummarizer{
		Usage:         make([]api.UploadActivity, 0),
		ActivityMap:   make(map[string]int),
		DeviceGlucose: make([]*GlucoseSummarizer, 0),
		Request:       request,
		Periods:       periods,
	}
}

//UploadIDToIndex translates an uploadid into an index into the Usage slice
func (a *ActivitySummarizer) UploadIDToIndex(uploadID *string) (int, error) {
	if uploadID == nil {
		return a.Intern(nil, nil), nil
	}
	offset, ok := a.ActivityMap[*uploadID]
	if !ok {
		return 0, fmt.Errorf("illegal offset for %v activity map %v", *uploadID, a.ActivityMap)
	}
	return offset, nil
}

//ProcessBG updates the time
func (a *ActivitySummarizer) ProcessBG(bg *data.Blood) error {
	if bg.Time == nil {
		return fmt.Errorf("no time provided %v", bg.ID)
	}

	offset, err := a.UploadIDToIndex(bg.UploadID)
	if err != nil {
		return err
	}
	if offset >= len(a.DeviceGlucose) {
		a.DeviceGlucose = append(a.DeviceGlucose, NewGlucoseSummarizer(a.Request, a.Periods))
	}
	//fmt.Printf("%d,%0.05f,%v,%v\n", offset, *bg.Value, *bg.Base.Time, *bg.Base.UploadID)

	return a.DeviceGlucose[offset].Process(bg)
}

//DeviceClientForUpload extracts a the device and client for an upload
func (a *ActivitySummarizer) DeviceClientForUpload(upload *data.Upload) (*api.Device, *api.Client) {
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
	return &device, &client
}

//Intern adds a canonical entry to the device/client table and returns the index to the canonical entry
func (a *ActivitySummarizer) Intern(device *api.Device, client *api.Client) int {
	for i, u := range a.Usage {
		if reflect.DeepEqual(u.Device, device) && reflect.DeepEqual(u.Client, client) {
			return i
		}
	}
	record := api.UploadActivity{
		Client: client,
		Device: device,
		Event:  &api.UpdateEvent{},
	}
	a.Usage = append(a.Usage, record)
	return len(a.Usage) - 1
}

//ProcessUpload intern device/client, add upload id to canonical entry to map
func (a *ActivitySummarizer) ProcessUpload(upload *data.Upload) {
	device, client := a.DeviceClientForUpload(upload)
	offset := a.Intern(device, client)
	if upload.Base.UploadID != nil {
		a.ActivityMap[*upload.Base.UploadID] = offset
	}
}

//Summary returns an activity summary
func (a *ActivitySummarizer) Summary() []api.UploadActivity {
	for i := range a.DeviceGlucose {
		a.Usage[i].Glucose = a.DeviceGlucose[i].Summary()
		a.Usage[i].Event.Time = a.DeviceGlucose[i].BGEnd
		if a.DeviceGlucose[i].Type == nil {
			t := "none"
			a.DeviceGlucose[i].Type = &t
		} else {
			a.Usage[i].Event.Type = *a.DeviceGlucose[i].Type
		}
	}
	return a.Usage
}
