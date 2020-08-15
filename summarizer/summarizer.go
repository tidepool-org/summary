package summarizer

import (
	"log"

	"github.com/tidepool-org/summary/api"
	"github.com/tidepool-org/summary/data"
)

type UserSummary struct {
	Glucose  *GlucoseSummarizer
	Activity ActivitySummarizer
}

//Summarizer creates summaries of upload activity
type Summarizer struct {
	Request   api.SummaryRequest
	Summaries map[string]*UserSummary
}

// NewSummarizer creates a Summarizer for the given request
func NewSummarizer(request api.SummaryRequest) *Summarizer {
	return &Summarizer{
		Request: request,
	}
}

func (s *Summarizer) SummaryForUser(userid string) *UserSummary {
	if summary, ok := s.Summaries[userid]; ok {
		return summary
	}
	s.Summaries[userid] = new(UserSummary)
	return s.Summaries[userid]
}

//Process an event
func (s *Summarizer) Process(rec interface{}) {
	switch v := rec.(type) {
	case data.Upload:
		s.SummaryForUser(*v.UserID).Activity.Process(&v)
	case data.Blood:
		s.SummaryForUser(*v.UserID).Glucose.Process(&v)
	default:
		log.Printf("skipping  %v\n", v)
	}
}

//Summary return summary report
func (s *Summarizer) Summary() []*api.SummaryResponse {
	summaries := make([]*api.SummaryResponse, 0)
	for userid, summary := range s.Summaries {
		summaries = append(summaries,
			&api.SummaryResponse{
				Activity: summary.Activity.Usage,
				Reports:  summary.Glucose.Summary(),
				Userid:   api.UserId(userid),
			},
		)
	}
	return summaries
}
