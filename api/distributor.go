// ChangeNotificationEvent is a notification of a change given event time
type ChangeNotificationEvent struct {
	Date   time.Time `json:"date"`   // event time of original record
	UserID string    `json:"userid"` // userid
	Kind   string    `kind:"kind"`   // enum: cbg, smbg, profile
}

//Distributor distributes change notifications to all listeners
func (s *Distributor) Run(ctx context.Context, in <-chan *ChangeNotificationEvent) {
	defer close(out)

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-in:
			s.Process(msg)
		}
	}
}

type Listener interface {
	Channel() chan *ChangeNotificationEvent
	Filter(*ChangeNotificationEvent) bool
}

type Disributor interface {
   Process( event *ChangeNotificationEvent)
   RegisterUserIdListener( userID string )
   RegisterDateListener( start, end time.Time ) 
}

type SimplerDistributor struct {
}

var _ Distributor = &SimplerDistributor{}

func (s *Distributor) Process(msg *ChangeNotificationEvent) {
   // send to all registered listeners
   // listeners can filter on userid or kind or daterange.

}