package distributor

import (
	"context"
	"time"
)

// ChangeNotificationEvent is a notification of a change given event time
type ChangeNotificationEvent struct {
	Date   time.Time `json:"date"`   // event time of original record
	UserID string    `json:"userid"` // userid
	Kind   string    `kind:"kind"`   // enum: cbg, smbg, profile
}

//Distributor distributes
type Distributor interface {
	Process(event *ChangeNotificationEvent)
	RegisterUserIDListener(userID string)
	RegisterDateListener(start, end time.Time)
}

//SimpleDistributor distributes
type SimpleDistributor struct {
}

var _ Distributor = &SimpleDistributor{}

//Run distributes change notifications to all listeners
func (s *SimpleDistributor) Run(ctx context.Context, in <-chan *ChangeNotificationEvent) {

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-in:
			s.Process(msg)
		}
	}
}

//Process processes
func (s *SimpleDistributor) Process(msg *ChangeNotificationEvent) {
	// send to all registered listeners
	// listeners can filter on userid or kind or daterange.
}

//RegisterUserIDListener processes
func (s *SimpleDistributor) RegisterUserIDListener(userID string) {
	// send to all registered listeners
	// listeners can filter on userid or kind or daterange.
}

//RegisterDateListener processes
func (s *SimpleDistributor) RegisterDateListener(start, end time.Time) {
	// send to all registered listeners
	// listeners can filter on userid or kind or daterange.
}

//Listener listens
type Listener interface {
	Channel() chan *ChangeNotificationEvent
	Filter(*ChangeNotificationEvent) bool
}
