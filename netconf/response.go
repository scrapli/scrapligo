package netconf

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/scrapli/scrapligo/logging"
)

// Netconf11ChunkPatternCompiled pre compiled pattern to match netconf 1.1 chunk pattern.
var Netconf11ChunkPatternCompiled = regexp.MustCompile(Version11ChunkPattern)

// Response the netconf response object.
type Response struct {
	Host               string
	Port               int
	ChannelInput       []byte
	XMLInput           interface{}
	RawResult          []byte
	Result             string
	StartTime          time.Time
	EndTime            time.Time
	ElapsedTime        float64
	FailedWhenContains [][]byte
	Failed             bool
	StripNamespaces    bool
	NetconfVersion     string
	ErrorMessages      [][]string
}

// NewNetconfResponse return a new netconf response object.
func NewNetconfResponse(
	host, netconfVersion string,
	port int,
	channelInput []byte,
	xmlInput interface{},
	stripNamespaces bool,
) *Response {
	failedWhenContains := make([][]byte, 0)
	failedWhenContains = append(
		failedWhenContains,
		[]byte("<rpc-error>"),
		[]byte("<rpc-errors>"),
		[]byte("</rpc-error>"),
		[]byte("</rpc-errors>"),
	)

	r := &Response{
		Host:               host,
		Port:               port,
		ChannelInput:       channelInput,
		XMLInput:           xmlInput,
		RawResult:          nil,
		Result:             "",
		StartTime:          time.Now(),
		EndTime:            time.Time{},
		ElapsedTime:        0,
		Failed:             true,
		NetconfVersion:     netconfVersion,
		StripNamespaces:    stripNamespaces,
		FailedWhenContains: failedWhenContains,
	}

	return r
}

// Record record a netconf response.
func (r *Response) Record(rawResult []byte) {
	r.EndTime = time.Now()
	r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()

	r.RawResult = rawResult

	// at this point the command has completed, so only thing that can "fail" it is there being
	// some bad output in the string matching a FailedWhenContains substr
	r.Failed = false

	for _, failedStr := range r.FailedWhenContains {
		if bytes.Contains(r.RawResult, failedStr) {
			r.Failed = true
			break
		}
	}

	if r.NetconfVersion == Version10 {
		r.recordResponse10()
	} else if r.NetconfVersion == Version11 {
		r.recordResponse11()
	}
}

func (r *Response) recordResponse10() {
	tmpResult := make([]byte, 0)
	tmpResult = bytes.TrimPrefix(tmpResult, []byte(XMLHeader))
	tmpResult = bytes.TrimSuffix(tmpResult, []byte(Version10DelimiterPattern))
	r.Result = string(tmpResult)
}

func (r *Response) validateChunkSize(chunkSize int, chunk []byte) {
	// does this need more ... "massaging" like scrapli?
	// chunk regex matches the newline before the chunk size or end of message delimiter, so we
	// subtract one for that newline char
	if len(chunk)-1 != chunkSize {
		r.Failed = true

		logging.LogError(
			logging.FormatLogMessage(
				"info",
				r.Host,
				r.Port,
				fmt.Sprintf("return element lengh invalid, expted: %d, got %d for element: %s\n",
					chunkSize,
					len(chunk)-1,
					chunk),
			),
		)
	}
}

func (r *Response) recordResponse11() {
	resultSectionLens := Netconf11ChunkPatternCompiled.FindAllSubmatch(r.RawResult, -1)

	joinedResult := make([]byte, 0)

	for _, resultSectionMatch := range resultSectionLens {
		chunkSize, _ := strconv.Atoi(string(resultSectionMatch[1]))
		chunk := resultSectionMatch[2]

		r.validateChunkSize(chunkSize, chunk)

		joinedResult = append(joinedResult, chunk[:len(chunk)-1]...)
	}

	joinedResult = bytes.TrimPrefix(joinedResult, []byte(XMLHeader))
	r.Result = string(joinedResult)
}
