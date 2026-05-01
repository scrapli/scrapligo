package cli

import (
	"bytes"
	"context"
	"math"
	"strings"
	"time"

	scrapligoconstants "github.com/scrapli/scrapligo/v2/constants"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

const (
	elapsedTimeMultiplierDivider = 100
)

// Result is a struct returned from all Cli operations.
type Result struct {
	Host                   string
	Port                   uint16
	Inputs                 []string
	ResultsRaw             [][]byte
	Results                []string
	StartTime              time.Time
	Splits                 []time.Time
	ElapsedTimeSeconds     float64
	ResultsFailedIndicator string
}

// NewResult prepares a new Result object from ffi integration pointers (the pointers we pass to
// zig for it to populate the values of stuff).
func NewResult(
	host string,
	port uint16,
	inputs []byte,
	startTime uint64,
	splits []uint64,
	resultsRaw []byte,
	results []byte,
	resultsFailedIndicator []byte,
) *Result {
	inputsS := strings.Split(string(inputs), scrapligoconstants.LibScrapliDelimiter)
	resultsS := strings.Split(string(results), scrapligoconstants.LibScrapliDelimiter)

	start := time.Unix(0, scrapligoutil.SafeUint64ToInt64(startTime))
	splitTimes := make([]time.Time, len(splits))

	for i, s := range splits {
		splitTimes[i] = time.Unix(0, scrapligoutil.SafeUint64ToInt64(s))
	}

	var elapsed float64

	if len(splitTimes) > 0 {
		elapsed = math.Round(
			splitTimes[len(splitTimes)-1].Sub(start).Seconds()*elapsedTimeMultiplierDivider,
		) / elapsedTimeMultiplierDivider
	}

	return &Result{
		Host:   host,
		Port:   port,
		Inputs: inputsS,
		ResultsRaw: bytes.Split(
			resultsRaw,
			[]byte(scrapligoconstants.LibScrapliDelimiter),
		),
		Results:                resultsS,
		StartTime:              start,
		Splits:                 splitTimes,
		ElapsedTimeSeconds:     elapsed,
		ResultsFailedIndicator: string(resultsFailedIndicator),
	}
}

// EndTime returns the end time of the Result. If there are no split times, it returns the start
// time.
func (r *Result) EndTime() time.Time {
	if len(r.Splits) == 0 {
		return r.StartTime
	}

	return r.Splits[len(r.Splits)-1]
}

// Result returns all results joined on newline chars.
func (r *Result) Result() string {
	if len(r.Results) == 0 {
		return ""
	}

	return strings.Join(r.Results, "\n")
}

// ResultRaw returns all raw results joined on newline chars.
func (r *Result) ResultRaw() []byte {
	if len(r.Results) == 0 {
		return nil
	}

	return bytes.Join(r.ResultsRaw, []byte("\n"))
}

// Failed returns true if any result has any failed indicator present.
func (r *Result) Failed() bool {
	return r.ResultsFailedIndicator != ""
}

// TextFsmParse parses recorded output w/ a provided textfsm template. The argument is interpreted
// as URL or filesystem path, for example,
// response.TextFsmParse("http://example.com/textfsm.template") or
// response.TextFsmParse("./local/textfsm.template"). Note that the content passed to textfsm is
// the content of the Result() method -- meaning, if there are multiple inputs, the full output
// contained in this Result object will be passed. If you have a Result object with multiple inputs
// and would like to only parse one of the results, simply invoke scrapligoutil.TextFsmParse
// directly with the content you wish to parse.
func (r *Result) TextFsmParse(ctx context.Context, path string) ([]map[string]any, error) {
	return scrapligoutil.TextFsmParse(ctx, r.Result(), path)
}
