package summarizer

import (
	"log"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
)

//UserSummarizer summarizes use activity
type UserSummarizer struct {
	Glucose  *GlucoseSummarizer
	Activity ActivitySummarizer
}

//NewUserSummarizer create a new User Summarizer
func NewUserSummarizer(request api.SummaryRequest) *UserSummarizer {
	return &UserSummarizer{
		Glucose: NewGlucoseSummarizer(request),
	}
}

//Summarizer creates summaries of upload activity
type Summarizer struct {
	Request   api.SummaryRequest
	Summaries map[string]*UserSummarizer
}

// NewSummarizer creates a Summarizer for the given request
func NewSummarizer(request api.SummaryRequest) *Summarizer {
	return &Summarizer{
		Request:   request,
		Summaries: make(map[string]*UserSummarizer),
	}
}

//SummaryForUser return summary for given user
func (s *Summarizer) SummaryForUser(userid string) *UserSummarizer {
	if summary, ok := s.Summaries[userid]; ok {
		return summary
	}
	s.Summaries[userid] = NewUserSummarizer(s.Request)
	return s.Summaries[userid]
}

//Process an event
func (s *Summarizer) Process(rec interface{}) {
	switch v := rec.(type) {
	case data.Upload:
		if v.UserID != nil {
			s.SummaryForUser(*v.UserID).Activity.Process(&v)
		} else {
			log.Printf("upload missing userid : userid  %v uploadid %v", v.Base.UserID, *v.Base.UploadID)
		}
	case data.Blood:
		if v.UserID != nil {
			s.SummaryForUser(*v.UserID).Glucose.Process(&v)
		} else {
			log.Printf("blood missing userid : userid  %v uploadid %v", v.Base.UserID, *v.Base.UploadID)

		}
	default:
		log.Printf("skipping  %v \n", v)
	}
}

//Summary return summary report
func (s *Summarizer) Summary() []*api.SummaryResponse {
	summaries := make([]*api.SummaryResponse, 0)
	for userid, summary := range s.Summaries {
		summaries = append(summaries,
			&api.SummaryResponse{
				Activity: summary.Activity.Usage,
				Glucose:  summary.Glucose.Summary(),
				Userid:   api.UserId(userid),
			},
		)
	}
	return summaries
}
