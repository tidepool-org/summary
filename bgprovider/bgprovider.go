package bgprovider

import (
	"context"
	"time"

	"github.com/tidepool-org/summary/data"
)

//BG combines a BG reading and the upload that it was a part of
type BG struct {
	Blood  *data.Blood
	Upload *data.Upload
}

//BGProvider provides a sequence of blood glucose readings
type BGProvider interface {
	Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BG)
}

// MockProvider provides a static sequence of BG values
type MockProvider struct {
}

var _ BGProvider = &MockProvider{}

//Get provide BG values on a channel, close channel when no more values
func (b *MockProvider) Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BG) {

	blood1 := data.Blood{}
	upload1 := data.Upload{}
	ch <- BG{&blood1, &upload1}
}
