package netconf

import (
	"math"
	"strings"
	"time"
	"unsafe"

	scrapligoffi "github.com/scrapli/scrapligo/v2/ffi"
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
	Result             string
	StartTime          time.Time
	EndTime            time.Time
	ElapsedTimeSeconds float64
	Failed             bool
	Warnings           []string
	Errors             []string

	resultRawJournal []byte
	resultRaw        []byte
}

// NewResult prepares a new Result object from ffi integration pointers (the pointers we pass to
// zig for it to populate the values of stuff).
func newResult(
	input,
	host string,
	port uint16,
	startTime uint64,
	endTime uint64,
	resultRawJournal []byte,
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
		Result:             result,
		StartTime:          start,
		EndTime:            end,
		ElapsedTimeSeconds: elapsed,
		Warnings:           strings.Split(string(warnings), "\n"),
		Errors:             strings.Split(string(errors), "\n"),
		resultRawJournal:   resultRawJournal,
	}

	if len(errors) > 0 {
		// only errors == failure, warnings are just... warnings
		r.Failed = true
	}

	return r
}

// ResultRaw returns the raw unprocessed result from the netconf server -- this includes framing
// and the like -- this can error due to the reconstruction process because we do not store the
// raw content, only a journal that we can use to rebuild it.
func (r *Result) ResultRaw() ([]byte, error) {
	if r.resultRaw != nil {
		// already reconstructed, ship it.
		return r.resultRaw, nil
	}

	if len(r.resultRawJournal) == 0 {
		// should never happen because we'll always have some framing in netconf land, but if
		// the journal is len 0 then that means processed/raw is same/same and we can just return
		// the processed and skip any extra work
		return []byte(r.Result), nil
	}

	ffiMap, err := scrapligoffi.GetMapping()
	if err != nil {
		return nil, err
	}

	var resultRawSize uintptr

	// libscrapli wont mutate this and this lets us *not* have to have another copy (rather a view
	// into the slice) of result (since we have to pass in a *[]byte but we store (for sanity
	// reasons) the result as a string) - so, unsafe sorta, but less allocations for no reason.
	resultB := unsafe.Slice(unsafe.StringData(r.Result), len(r.Result)) //nolint: gosec

	err = ffiMap.Netconf.GetReconstructedResultRawSize(
		&resultB,
		&r.resultRawJournal,
		&resultRawSize,
	)
	if err != nil {
		return nil, err
	}

	reconstructedRaw := make([]byte, resultRawSize)

	err = ffiMap.Netconf.GetReconstructedResultRaw(
		&resultB,
		&r.resultRawJournal,
		&reconstructedRaw,
	)
	if err != nil {
		return nil, err
	}

	r.resultRaw = reconstructedRaw

	return r.resultRaw, nil
}
