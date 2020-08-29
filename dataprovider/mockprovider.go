package dataprovider

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
)

const (
	//Layout is how time is represented in the API
	Layout = "2006-01-02T15:04:05Z"
)

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
