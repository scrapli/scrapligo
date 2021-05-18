package base

import (
	"strings"
	"time"
)

// Response response object that gets returned from scrapli send operations.
type Response struct {
	Host               string
	Port               int
	ChannelInput       string
	RawResult          []byte
	Result             string
	StartTime          time.Time
	EndTime            time.Time
	ElapsedTime        float64
	FailedWhenContains []string
	Failed             bool
}

// NewResponse create a new response object.
func NewResponse(
	host string,
	port int,
	channelInput string,
	failedWhenContains []string,
) *Response {
	r := &Response{
		Host:               host,
		Port:               port,
		ChannelInput:       channelInput,
		RawResult:          nil,
		Result:             "",
		StartTime:          time.Now(),
		EndTime:            time.Time{},
		ElapsedTime:        0,
		Failed:             true,
		FailedWhenContains: failedWhenContains,
	}

	return r
}

// Record records a response from an operation.
func (r *Response) Record(rawResult []byte, result string) {
	r.EndTime = time.Now()
	r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()

	r.RawResult = rawResult
	r.Result = result

	// at this point the command has completed, so only thing that can "fail" it is there being
	// some bad output in the string matching a FailedWhenContains substr
	r.Failed = false

	for _, failedStr := range r.FailedWhenContains {
		if strings.Contains(r.Result, failedStr) {
			r.Failed = true
			break
		}
	}
}
