package netconf

import "strings"

// Result is a struct returned from all Cli operations.
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
	r := &Result{
		Host:      host,
		Port:      port,
		Input:     input,
		ResultRaw: resultRaw,
		Result:    result,
		StartTime: startTime,
		EndTime:   endTime,
		Warnings:  strings.Split(string(warnings), "\n"),
		Errors:    strings.Split(string(errors), "\n"),
	}

	if len(errors) > 0 {
		// only errors == failure, warnings are just... warnings
		r.Failed = true
	}

	return r
}
