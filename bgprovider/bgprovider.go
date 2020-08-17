package bgprovider

import (
	"context"
	"time"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
)

//BG is a data record, usual cbg, bg, or upload
type BG interface{}

//BGProvider provides a sequence of blood glucose readings
type BGProvider interface {
	Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BG, continuous bool)
}

// MockProvider provides a static sequence of BG values
type MockProvider struct {
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
func (b *MockProvider) Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BG, continuous bool) {
	ch <- data.Upload{
		Base: data.Base{
			Active:   true,
			DeviceID: PStr("foo"),
			Time:     PStr("2020-07-11T08:29:02Z"),
			Type:     "upload",
			UploadID: PStr("xyz"),
			UserID:   PStr("foo"),
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
	ch <- data.Blood{
		Base: data.Base{
			Active:   true,
			DeviceID: PStr("foo"),
			Time:     PStr("2020-07-11T08:29:02Z"),
			Type:     "cbg",
			UploadID: PStr("xyz"),
			UserID:   PStr("foo"),
		},
		Units: PStr("mg/dl"),
		Value: PF64(130.0),
	}
	close(ch)
}
