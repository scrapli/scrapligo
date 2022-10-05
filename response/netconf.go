package response

import (
	"bytes"
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
)

type netconfPatterns struct {
	v1dot1Chunk *regexp.Regexp
	rpcErrors   *regexp.Regexp
}

var (
	netconfPatternsInstance     *netconfPatterns //nolint:gochecknoglobals
	netconfPatternsInstanceOnce sync.Once        //nolint:gochecknoglobals
)

func getNetconfPatterns() *netconfPatterns {
	netconfPatternsInstanceOnce.Do(func() {
		netconfPatternsInstance = &netconfPatterns{
			v1dot1Chunk: regexp.MustCompile(`(?ms)(\d+)\n(.*?)#`),
			rpcErrors:   regexp.MustCompile(`(?s)<rpc-errors?>(.*)</rpc-errors?>`),
		}
	})

	return netconfPatternsInstance
}

// NewNetconfResponse prepares a new NetconfResponse object.
func NewNetconfResponse(
	input []byte,
	host string,
	port int,
	version string,
) *NetconfResponse {
	return &NetconfResponse{
		Host:        host,
		Port:        port,
		Input:       input,
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

	r.Result = string(b)
}

func (r *NetconfResponse) validateChunk(i int, b []byte) {
	// does this need more ... "massaging" like scrapli?
	// chunk regex matches the newline before the chunk size or end of message delimiter, so we
	// subtract one for that newline char
	if len(b)-1 != i {
		errMsg := fmt.Sprintf("return element lengh invalid, expted: %d, got %d for element: %s\n",
			i,
			len(b)-1,
			b)

		r.Failed = &OperationError{
			Input:       string(r.Input),
			Output:      r.Result,
			ErrorString: errMsg,
		}
	}
}

func (r *NetconfResponse) record1dot1() {
	patterns := getNetconfPatterns()

	chunkSections := patterns.v1dot1Chunk.FindAllSubmatch(r.RawResult, -1)

	var joined []byte

	for _, chunkSection := range chunkSections {
		chunk := chunkSection[2]

		size, _ := strconv.Atoi(string(chunkSection[1]))

		r.validateChunk(size, chunk)

		joined = append(joined, chunk[:len(chunk)-1]...)
	}

	joined = bytes.TrimPrefix(joined, []byte(xmlHeader))

	r.Result = string(bytes.TrimSpace(joined))
}
