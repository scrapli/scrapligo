package cfg

import (
	"strings"
	"time"

	"github.com/scrapli/scrapligo/difflibgo"

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

type sideBySideDiffFunc func(response *DiffResponse) string


// DiffResponse cfg diff response object that gets returned from diff operations.
type DiffResponse struct {
	*Response
	Source          string
	SourceConfig    string
	CandidateConfig string
	DeviceDiff      string
	difflines       []string
	additions       []string
	subtractions    []string
	UnifiedDiff     string
	sideBySideDiff  string
	SideBySideDiff  sideBySideDiffFunc
	colorize        bool
	sideBySideWidth int
}

func (r *DiffResponse) RecordDiff(sourceConfig, candidateConfig, deviceDiff string) {
	r.SourceConfig = sourceConfig
	r.CandidateConfig = candidateConfig
	r.DeviceDiff = deviceDiff

	d := difflibgo.Differ{}
	r.difflines = d.Compare(
		strings.Split(r.SourceConfig, "\n"),
		strings.Split(r.CandidateConfig, "\n"),
	)

	for _, v := range r.difflines {
		if v[:2] == "+ " {
			r.additions = append(r.additions, v[2:])
		} else if v[:2] == "- " {
			r.subtractions = append(r.additions, v[2:])
		}
	}
}

// NewDiffResponse create a new cfg diff response object.
func NewDiffResponse(
	host string,
	source string,
	colorize bool,
	sideBySideWidth int,
) *DiffResponse {
	r := &Response{
		Host:        host,
		Result:      "",
		StartTime:   time.Now(),
		EndTime:     time.Time{},
		ElapsedTime: 0,
		ErrorType:   ErrDiffConfigFailed,
		Failed:      true,
	}

	dr := &DiffResponse{
		Response:        r,
		Source:          source,
		colorize:        colorize,
		sideBySideWidth: sideBySideWidth,
	}

	return dr
}
