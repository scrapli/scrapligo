package base

import (
	"strings"
	"time"
)

// MultiResponse defines a response object for plural operations -- contains a slice of `Responses`
// for each operation.
type MultiResponse struct {
	Host        string
	StartTime   time.Time
	EndTime     time.Time
	ElapsedTime float64
	Responses   []*Response
}

// NewMultiResponse create a new MultiResponse object.
func NewMultiResponse(host string) *MultiResponse {
	r := &MultiResponse{
		Host:        host,
		StartTime:   time.Now(),
		EndTime:     time.Time{},
		ElapsedTime: 0,
	}

	return r
}

// AppendResponse append a response to the `Responses` attribute of the MultiResponse object.
func (mr *MultiResponse) AppendResponse(r *Response) {
	mr.Responses = append(mr.Responses, r)
}

// JoinedResult convenience method to print out a single unified/joined result -- this is basically
// all of the channel output for each individual response in the MultiResponse joined by newline
// characters.
func (mr *MultiResponse) JoinedResult() string {
	resultsSlice := make([]string, len(mr.Responses))

	for _, resp := range mr.Responses {
		resultsSlice = append(resultsSlice, resp.Result)
	}

	return strings.Join(resultsSlice, "\n")
}

// Failed method indicating if the MultiResponse is failed -- if any Response in the MultiResponse
// is failed, return true.
func (mr *MultiResponse) Failed() bool {
	// can this be made a property or property-like in python? dont see an obvious way to
	//  have this be an attribute of the struct as it doesnt seem like it can take the receiver
	//  that points to the object itself?
	for _, r := range mr.Responses {
		if r.Failed {
			return true
		}
	}

	return false
}
