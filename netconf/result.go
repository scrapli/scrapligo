package netconf

import (
	"math"
	"strings"
	"time"

	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

const (
	elapsedTimeMultiplierDivider = 100
)

// Result is a struct returned from all Cli operations.
type Result struct {
	Host               string
	Port               uint16
	Input              string
	ResultRaw          []byte
	Result             string
	StartTime          time.Time
	EndTime            time.Time
	ElapsedTimeSeconds float64
	Failed             bool
	Warnings           []string
	Errors             []string
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
	warnings []byte,
	errors []byte,
) *Result {
	start := time.Unix(0, scrapligoutil.SafeUint64ToInt64(startTime))
	end := time.Unix(0, scrapligoutil.SafeUint64ToInt64(endTime))

	elapsed := math.Round(
		end.Sub(start).Seconds()*elapsedTimeMultiplierDivider,
	) / elapsedTimeMultiplierDivider

	r := &Result{
		Host:               host,
		Port:               port,
		Input:              input,
		ResultRaw:          resultRaw,
		Result:             result,
		StartTime:          start,
		EndTime:            end,
		ElapsedTimeSeconds: elapsed,
		Warnings:           strings.Split(string(warnings), "\n"),
		Errors:             strings.Split(string(errors), "\n"),
	}

	if len(errors) > 0 {
		// only errors == failure, warnings are just... warnings
		r.Failed = true
	}

	return r
}
