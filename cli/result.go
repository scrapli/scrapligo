package cli

import (
	"bytes"
	"math"
	"strings"
	"time"
)

const (
	elapsedTimeMultiplierDivider = 100
	operationDelimiter           = "__libscrapli__"
)

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
	host string,
	port uint16,
	inputs []byte,
	startTime uint64,
	splits []uint64,
	resultsRaw []byte,
	results []byte,
	resultsFailedIndicator []byte,
) *Result {
	inputsS := strings.Split(string(inputs), operationDelimiter)
	resultsS := strings.Split(string(results), operationDelimiter)

	var elapsed float64

	if len(splits) > 0 {
		elapsed = elapsedTime(startTime, splits[len(splits)-1])
	}

	return &Result{
		Host:                   host,
		Port:                   port,
		Inputs:                 inputsS,
		ResultsRaw:             bytes.Split(resultsRaw, []byte(operationDelimiter)),
		Results:                resultsS,
		StartTime:              startTime,
		Splits:                 splits,
		ElapsedTimeSeconds:     elapsed,
		ResultsFailedIndicator: string(resultsFailedIndicator),
	}
}

// Result is a struct returned from all Driver operations.
type Result struct {
	Host                   string
	Port                   uint16
	Inputs                 []string
	ResultsRaw             [][]byte
	Results                []string
	StartTime              uint64
	Splits                 []uint64
	ElapsedTimeSeconds     float64
	ResultsFailedIndicator string
}

func (r *Result) extend(res *Result) error {
	r.Inputs = append(r.Inputs, res.Inputs...)
	r.ResultsRaw = append(r.ResultsRaw, res.ResultsRaw...)
	r.Results = append(r.Results, res.Results...)
	r.Splits = append(r.Splits, res.Splits...)

	if len(res.Splits) > 0 {
		r.ElapsedTimeSeconds = elapsedTime(r.StartTime, res.Splits[len(res.Splits)-1])
	}

	return nil
}

func (r *Result) EndTime() uint64 {
	if len(r.ResultsRaw) == 0 {
		return 0
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

func (r *Result) Failed() bool {
	return len(r.ResultsFailedIndicator) > 0
}

// TextFsmParse parses recorded output w/ a provided textfsm template.
// the argument is interpreted as URL or filesystem path, for example:
// response.TextFsmParse("http://example.com/textfsm.template") or
// response.TextFsmParse("./local/textfsm.template").
func (r *Result) TextFsmParse(path string) ([]map[string]interface{}, error) {
	// TODO
	// return scrapligoutil.TextFsmParse(r.Result, path)
	return nil, nil
}
