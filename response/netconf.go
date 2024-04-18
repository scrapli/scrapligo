package response

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/scrapli/scrapligo/util"
)

const (
	v1Dot0      = "1.0"
	v1Dot1      = "1.1"
	v1Dot0Delim = "]]>]]>"
	xmlHeader   = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	// from https://datatracker.ietf.org/doc/html/rfc6242#section-4.2
	v1Dot1MaxChunkSize = 4294967295
)

var errNetconf1Dot1Error = errors.New("unable to parse netconf 1.1 response")

func errNetconf1Dot1ParseError(msg string) error {
	return fmt.Errorf("%w: %s", errNetconf1Dot1Error, msg)
}

type netconfPatterns struct {
	rpcErrors             *regexp.Regexp
	v1Dot1MaxChunkSizeLen int
}

var (
	netconfPatternsInstance     *netconfPatterns //nolint:gochecknoglobals
	netconfPatternsInstanceOnce sync.Once        //nolint:gochecknoglobals
)

func getNetconfPatterns() *netconfPatterns {
	netconfPatternsInstanceOnce.Do(func() {
		netconfPatternsInstance = &netconfPatterns{
			rpcErrors:             regexp.MustCompile(`(?s)<rpc-errors?>(.*)</rpc-errors?>`),
			v1Dot1MaxChunkSizeLen: len(strconv.Itoa(v1Dot1MaxChunkSize)),
		}
	})

	return netconfPatternsInstance
}

// NewNetconfResponse prepares a new NetconfResponse object.
func NewNetconfResponse(
	input,
	framedInput []byte,
	host string,
	port int,
	version string,
) *NetconfResponse {
	return &NetconfResponse{
		Host:        host,
		Port:        port,
		Input:       input,
		FramedInput: framedInput,
		Result:      "",
		StartTime:   time.Now(),
		EndTime:     time.Time{},
		ElapsedTime: 0,
		FailedWhenContains: [][]byte{
			[]byte("<rpc-error>"),
			[]byte("<rpc-errors>"),
			[]byte("</rpc-error>"),
			[]byte("</rpc-errors>"),
		},
		NetconfVersion: version,
	}
}

// NetconfResponse is a struct returned from all netconf driver operations.
type NetconfResponse struct {
	Host               string
	Port               int
	Input              []byte
	FramedInput        []byte
	RawResult          []byte
	Result             string
	StartTime          time.Time
	EndTime            time.Time
	ElapsedTime        float64
	FailedWhenContains [][]byte
	Failed             error
	StripNamespaces    bool
	NetconfVersion     string
	ErrorMessages      [][]string
	SubscriptionID     int
}

// Record records the output of a NETCONF operation.
func (r *NetconfResponse) Record(b []byte) {
	r.EndTime = time.Now()
	r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()

	r.RawResult = b

	if util.ByteContainsAny(r.RawResult, r.FailedWhenContains) {
		patterns := getNetconfPatterns()

		r.Failed = &OperationError{
			Input:       string(r.Input),
			Output:      r.Result,
			ErrorString: string(patterns.rpcErrors.Find(r.RawResult)),
		}
	}

	switch r.NetconfVersion {
	case v1Dot0:
		r.record1dot0()
	case v1Dot1:
		r.record1dot1()
	}
}

func (r *NetconfResponse) record1dot0() {
	b := r.RawResult

	b = bytes.TrimPrefix(b, []byte(xmlHeader))
	// trim space before trimming suffix because we usually have a trailing newline!
	b = bytes.TrimSuffix(bytes.TrimSpace(b), []byte(v1Dot0Delim))

	r.Result = string(bytes.TrimSpace(b))
}

func (r *NetconfResponse) record1dot1() {
	joined, err := r.record1dot1Chunks(r.RawResult)
	if err != nil {
		r.Failed = &OperationError{
			Input:       string(r.Input),
			Output:      r.Result,
			ErrorString: err.Error(),
		}
	}

	joined = bytes.TrimPrefix(joined, []byte(xmlHeader))

	r.Result = string(bytes.TrimSpace(joined))
}

func (r *NetconfResponse) record1dot1Chunks(d []byte) ([]byte, error) {
	pattern := getNetconfPatterns()

	cursor := 0

	joined := []byte{}

	for cursor < len(d) {
		// allow for some amount of newlines
		if d[cursor] == byte('\n') {
			cursor++

			continue
		}

		if d[cursor] != byte('#') {
			return nil, errNetconf1Dot1ParseError(fmt.Sprintf(
				"unable to parse netconf response: chunk marker missing, got '%s'",
				string(d[cursor])),
			)
		}

		cursor++

		// found prompt
		if d[cursor] == byte('#') {
			return joined, nil
		}

		// look for end of chunk size
		// allow to match end with \n char
		chunkSizeLen := 0
		for ; chunkSizeLen < pattern.v1Dot1MaxChunkSizeLen+1; chunkSizeLen++ {
			if cursor+chunkSizeLen >= len(d) {
				return nil, errNetconf1Dot1ParseError("chunk size not found before end of data")
			}

			if d[cursor+chunkSizeLen] == byte('\n') {
				break
			}
		}

		chunkSizeStr := string(d[cursor : cursor+chunkSizeLen])
		cursor += chunkSizeLen + 1

		chunkSize, err := strconv.Atoi(chunkSizeStr)
		if err != nil {
			return nil, errNetconf1Dot1ParseError(
				fmt.Sprintf("unable to parse chunk size '%s': %s", chunkSizeStr, err),
			)
		}

		joined = append(joined, d[cursor:cursor+chunkSize]...)
		// last new line of block is not counted
		// since it's considered a delimiter for next chunk
		cursor += chunkSize + 1
	}

	return joined, nil
}
