package cli

import (
	"bytes"
	"context"
	"math"
	"strings"
	"time"
	"unsafe"

	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
	scrapligoutil "github.com/scrapli/scrapligo/v2/util"
)

const (
	elapsedTimeMultiplierDivider = 100
)

// Result is a struct returned from all Cli operations.
type Result struct {
	Host string
	Port uint16

	StartTime              time.Time
	Splits                 []time.Time
	ElapsedTimeSeconds     float64
	ResultsFailedIndicator string

	inputs    []byte
	inputLens []uint64

	results    []byte
	resultLens []uint64

	resultRawJournals    []byte
	resultRawJournalLens []uint64
	resultsRaw           []byte
}

// newResult prepares a new Result object from ffi integration pointers (the pointers we pass to
// zig for it to populate the values of stuff).
func newResult(
	host string,
	port uint16,
	startTime uint64,
	splits []uint64,
	inputs []byte,
	inputLens []uint64,
	resultRawJournals []byte,
	resultRawJournalLens []uint64,
	results []byte,
	resultLens []uint64,
	resultsFailedIndicator []byte,
) *Result {
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
		Host:                   host,
		Port:                   port,
		StartTime:              start,
		Splits:                 splitTimes,
		ElapsedTimeSeconds:     elapsed,
		ResultsFailedIndicator: string(resultsFailedIndicator),

		inputs:    inputs,
		inputLens: inputLens,

		results:    results,
		resultLens: resultLens,

		resultRawJournals:    resultRawJournals,
		resultRawJournalLens: resultRawJournalLens,
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
	if len(r.results) == 0 {
		return ""
	}

	var outSize int

	for _, size := range r.resultLens {
		outSize += int(size) //nolint: gosec
	}

	// plus newlines we will add between results
	outSize += len(r.resultLens) - 1

	var out strings.Builder

	out.Grow(outSize)

	var cur uint64

	for idx, resultLen := range r.resultLens {
		out.Write(r.results[cur : cur+resultLen])

		cur += resultLen

		if idx < len(r.resultLens)-1 {
			out.WriteString("\n")
		}
	}

	return out.String()
}

// ResultAtIndex returns the result at the given index (rather than Result which returns all
// results as one joined string).
func (r *Result) ResultAtIndex(index int) (string, error) {
	if index >= len(r.resultLens) {
		return "", scrapligoerrors.NewFfiError(
			"index error, result",
			nil,
		)
	}

	outSize := int(r.resultLens[index]) //nolint: gosec

	var out strings.Builder

	out.Grow(outSize)

	var startPos int

	for _, resultLen := range r.resultLens[0:index] {
		startPos += int(resultLen) //nolint: gosec
	}

	out.Write(r.results[startPos : startPos+outSize])

	return out.String(), nil
}

// ResultRaw returns all raw results joined on newline chars. Can error as we rebuild the raw
// from the raw journal via libscrapli which means we can fail to get the mapping... that should
// never happen, but it *could*.
func (r *Result) ResultRaw() ([]byte, error) {
	if r.resultsRaw != nil {
		// must be before the len check because read with callbacks will be manually setting
		// this value in due to the way the read any bits work
		return r.resultsRaw, nil
	}

	if len(r.resultRawJournals) == 0 {
		// unlike netconf this *can* and likely will happen, in this case we have no journal entries
		// meaning there was no ansi/ascii stripping bits happening, so result raw is the same as
		// result so we just return the result so we dont have to make any more calls to libscrapli
		return []byte(r.Result()), nil
	}

	reconstructedSegments := make([][]byte, len(r.resultRawJournalLens))

	for index := range r.resultRawJournalLens {
		reconstructed, err := r.ResultRawAtIndex(index)
		if err != nil {
			return nil, err
		}

		reconstructedSegments[index] = reconstructed
	}

	r.resultsRaw = bytes.Join(reconstructedSegments, []byte("\n"))

	return r.resultsRaw, nil
}

// ResultRawAtIndex returns the (raw) result at the given index (rather than RawResult which returns
// all results as one joined byte slice). Same note regarding erroring as ResultRaw.
func (r *Result) ResultRawAtIndex(index int) ([]byte, error) {
	if index >= len(r.resultRawJournalLens) {
		return nil, scrapligoerrors.NewFfiError(
			"index error, raw journal",
			nil,
		)
	}

	outSize := int(r.resultRawJournalLens[index]) //nolint: gosec

	var startPos int

	for _, journalLen := range r.resultRawJournalLens[0:index] {
		startPos += int(journalLen) //nolint: gosec
	}

	journalEntry := r.resultRawJournals[startPos : startPos+outSize]

	result, err := r.ResultAtIndex(index)
	if err != nil {
		return nil, err
	}

	if outSize == 0 {
		return []byte(result), nil
	}

	ffiMap, err := scrapligoffi.GetMapping()
	if err != nil {
		return nil, err
	}

	resultB := unsafe.Slice(unsafe.StringData(result), len(result)) //nolint: gosec

	var reconstructedSize uintptr

	err = ffiMap.Cli.GetReconstructedResultRawSize(
		&resultB,
		&journalEntry,
		&reconstructedSize,
	)
	if err != nil {
		return nil, err
	}

	out := make([]byte, reconstructedSize)

	err = ffiMap.Cli.GetReconstructedResultRaw(
		&resultB,
		&journalEntry,
		&out,
	)
	if err != nil {
		return nil, err
	}

	return out, nil
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
