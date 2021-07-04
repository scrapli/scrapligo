package cfg

import (
	"time"

	"github.com/scrapli/scrapligo/driver/base"
)

// Response cfg response object that gets returned from cfg operations.
type Response struct {
	Host             string
	Result           string
	StartTime        time.Time
	EndTime          time.Time
	ElapsedTime      float64
	ScrapliResponses []*base.Response
	ErrorType        error
	Failed           bool
}

// NewResponse create a new cfg response object.
func NewResponse(
	host string,
	raiseError error,
) *Response {
	r := &Response{
		Host:        host,
		Result:      "",
		StartTime:   time.Now(),
		EndTime:     time.Time{},
		ElapsedTime: 0,
		ErrorType:   raiseError,
		Failed:      true,
	}

	return r
}

func (r *Response) Record(scrapliResponses []*base.Response, result string) {
	r.EndTime = time.Now()
	r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()

	r.Result = result
	r.Failed = false
	r.ScrapliResponses = scrapliResponses

	for _, response := range r.ScrapliResponses {
		if response.Failed {
			r.Failed = true
			break
		}
	}
}

// DiffResponse cfg diff response object that gets returned from diff operations.
type DiffResponse struct {
}
