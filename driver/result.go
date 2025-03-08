package driver

import (
	"math"
	"strings"
	"time"

	scrapligoutil "github.com/scrapli/scrapligo/util"
)

const elapsedTimeMultiplierDivider = 100

func elapsedTime(startTime, endTime uint64) float64 {
	elapsed := endTime - startTime

	if elapsed >= math.MaxInt64 {
		return math.MaxInt64
	}

	return math.Round(
		time.Duration(elapsed).Seconds()*elapsedTimeMultiplierDivider,
	) / elapsedTimeMultiplierDivider
}

// NewResult prepares a new Result object from ffi integration pointers (the pointers we pass to
// zig for it to populate the values of stuff).
func NewResult(
	input,
	host string,
	port uint16,
	startTime uint64,
	endTime uint64,
	resultRaw []byte,
	result string,
	resultFailedWhenIndicator []byte,
) *Result {
	return &Result{
		Host:               host,
		Port:               port,
		Input:              input,
		ResultRaw:          resultRaw,
		Result:             result,
		StartTime:          startTime,
		EndTime:            endTime,
		ElapsedTimeSeconds: elapsedTime(startTime, endTime),
		Failed:             len(resultFailedWhenIndicator) > 0,
		FailedIndicator:    string(resultFailedWhenIndicator),
	}
}

// NewMultiResult creates a new MultiResult object.
func NewMultiResult(
	host string,
	port uint16,
) *MultiResult {
	return &MultiResult{
		Host: host,
		Port: port,
	}
}

// Result is a struct returned from all Driver operations.
type Result struct {
	Host               string
	Port               uint16
	Input              string
	ResultRaw          []byte
	Result             string
	StartTime          uint64
	EndTime            uint64
	ElapsedTimeSeconds float64
	Failed             bool
	FailedIndicator    string
}

// TextFsmParse parses recorded output w/ a provided textfsm template.
// the argument is interpreted as URL or filesystem path, for example:
// response.TextFsmParse("http://example.com/textfsm.template") or
// response.TextFsmParse("./local/textfsm.template").
func (r *Result) TextFsmParse(path string) ([]map[string]interface{}, error) {
	return scrapligoutil.TextFsmParse(r.Result, path)
}

// MultiResult defines a response object for plural operations -- contains a slice of `Result`
// for each operation.
type MultiResult struct {
	Host               string
	Port               uint16
	Results            []*Result
	StartTime          uint64
	EndTime            uint64
	ElapsedTimeSeconds float64
}

// AppendResult appends a `Result` to the `MultiResult`.
func (mr *MultiResult) AppendResult(r *Result) {
	if mr.StartTime == 0 {
		mr.StartTime = r.StartTime
	}

	mr.EndTime = r.EndTime
	mr.ElapsedTimeSeconds = elapsedTime(r.StartTime, r.EndTime)

	mr.Results = append(mr.Results, r)
}

// JoinedResult is a convenience method to print out a single unified/joined result joined by
// newlines.
func (mr *MultiResult) JoinedResult() string {
	resultsSlice := make([]string, len(mr.Results))

	for idx, resp := range mr.Results {
		resultsSlice[idx] = resp.Result
	}

	return strings.Join(resultsSlice, "\n")
}
