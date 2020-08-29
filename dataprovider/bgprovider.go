package dataprovider

import (
	"context"
	"time"
)

//BG is a data record, usual cbg, bg, or upload
type BG interface{}

//BGProvider provides a sequence of blood glucose readings
type BGProvider interface {
	Get(ctx context.Context, from time.Time, to time.Time, ch chan<- BG, users []string)
}
