package summarizer

import (
	"log"
	"time"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
)

//UserSummarizer summarizes use activity
type UserSummarizer struct {
	Activity *ActivitySummarizer
}

//NewPeriods creates sample periods
func NewPeriods(request api.SummaryRequest) []api.SummaryPeriod {
	periods := make([]api.SummaryPeriod, request.Period.NumPeriods)
	now := time.Now()
	ending := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
	var nDays int
	switch request.Period.Length {
	case "day":
		nDays = 1
	case "week":
		nDays = 7
	}

	for i := range periods {
		periods[i].End = ending.UTC()
		ending = ending.AddDate(0, 0, -nDays)
		periods[i].Start = ending.UTC()
		periods[i].Length = request.Period.Length
		periods[i].Updated = now.UTC()
	}
	return periods
}

//LastThousandDays creates a single period covering the last 1000 days
func LastThousandDays(request api.SummaryRequest) []api.SummaryPeriod {
	periods := make([]api.SummaryPeriod, 1)
	now := time.Now()
	ending := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)

	for i := range periods {
		periods[i].End = ending.UTC()
		ending = ending.AddDate(0, 0, -1000)
		periods[i].Start = ending.UTC()
		periods[i].Length = request.Period.Length
		periods[i].Updated = now.UTC()
	}
	return periods
}

//NewUserSummarizer create a new User Summarizer
func NewUserSummarizer(request api.SummaryRequest, periods []api.SummaryPeriod) *UserSummarizer {
	return &UserSummarizer{
		Activity: NewActivitySummarizer(request, periods),
	}
}

//Summarizer creates summaries of upload activity
type Summarizer struct {
	Request   api.SummaryRequest
	Summaries map[string]*UserSummarizer
	Periods   []api.SummaryPeriod
}

// NewSummarizer creates a Summarizer for the given request
func NewSummarizer(request api.SummaryRequest) *Summarizer {
	return &Summarizer{
		Request:   request,
		Summaries: make(map[string]*UserSummarizer),
		Periods:   NewPeriods(request),
		//Periods: LastThousandDays(request),
	}
}

//DateRange provide the times needed to produce the reports
func (s *Summarizer) DateRange() (from, to time.Time) {
	from = s.Periods[0].Start
	to = s.Periods[len(s.Periods)-1].End
	log.Printf("from %v, to %v", from, to)
	return
}

//SummarizerForUser return summary for given user
func (s *Summarizer) SummarizerForUser(userid string) *UserSummarizer {
	if summary, ok := s.Summaries[userid]; ok {
		return summary
	}
	s.Summaries[userid] = NewUserSummarizer(s.Request, s.Periods)
	return s.Summaries[userid]
}

//Process an event
func (s *Summarizer) Process(rec interface{}) {
	switch v := rec.(type) {
	case data.Upload:
		if v.UserID != nil {
			s.SummarizerForUser(*v.UserID).Activity.ProcessUpload(&v)
		} else {
			log.Printf("upload missing userid : userid  %v uploadid %v", v.Base.UserID, *v.Base.UploadID)
		}
	case data.Blood:
		if v.UserID != nil {
			s.SummarizerForUser(*v.UserID).Activity.ProcessBG(&v)
		} else {
			log.Printf("blood missing userid : userid  %v uploadid %v", v.Base.UserID, *v.Base.UploadID)

		}
	default:
		log.Printf("unexpected data type returned %v", v)
	}
}

//Summary return summary report
func (s *Summarizer) Summary() []*api.SummaryResponse {
	summaries := make([]*api.SummaryResponse, 0)
	for userid, summary := range s.Summaries {
		summaries = append(summaries,
			&api.SummaryResponse{
				Activity: summary.Activity.Summary(),
				Userid:   api.UserId(userid),
			},
		)
	}
	return summaries
}
