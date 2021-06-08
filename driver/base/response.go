package base

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/scrapli/scrapligo/logging"
	"github.com/sirikothe/gotextfsm"
)

// ErrFailedOpeningTemplate error for failure to open a textfsm template.
var ErrFailedOpeningTemplate = errors.New("failed opening provided path to textfsm template")

// ErrFailedParsingTemplate error for failure of parsing a textfsm template.
var ErrFailedParsingTemplate = errors.New("failed parsing textfsm template")

// Response response object that gets returned from scrapli send operations.
type Response struct {
	Host               string
	Port               int
	ChannelInput       string
	RawResult          []byte
	Result             string
	StartTime          time.Time
	EndTime            time.Time
	ElapsedTime        float64
	FailedWhenContains []string
	Failed             bool
}

// NewResponse create a new response object.
func NewResponse(
	host string,
	port int,
	channelInput string,
	failedWhenContains []string,
) *Response {
	r := &Response{
		Host:               host,
		Port:               port,
		ChannelInput:       channelInput,
		RawResult:          nil,
		Result:             "",
		StartTime:          time.Now(),
		EndTime:            time.Time{},
		ElapsedTime:        0,
		Failed:             true,
		FailedWhenContains: failedWhenContains,
	}

	return r
}

// Record records a response from an operation.
func (r *Response) Record(rawResult []byte, result string) {
	r.EndTime = time.Now()
	r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()

	r.RawResult = rawResult
	r.Result = result

	// at this point the command has completed, so only thing that can "fail" it is there being
	// some bad output in the string matching a FailedWhenContains substr
	r.Failed = false

	for _, failedStr := range r.FailedWhenContains {
		if strings.Contains(r.Result, failedStr) {
			r.Failed = true
			break
		}
	}
}

// TextFsmParse parses recorded output w/ a provided textfsm template.
func (r *Response) TextFsmParse(template string) ([]map[string]interface{}, error) {
	t, err := os.ReadFile(template)
	if err != nil {
		logging.LogError(
			r.FormatLogMessage(
				"error",
				fmt.Sprintf("Failed opening provided template, error: %s\n", err.Error()),
			),
		)

		return []map[string]interface{}{}, ErrFailedOpeningTemplate
	}

	fsm := gotextfsm.TextFSM{}

	err = fsm.ParseString(string(t))
	if err != nil {
		logging.LogError(
			r.FormatLogMessage(
				"error",
				fmt.Sprintf("Failed parsing provided template, gotextfsm error: %s\n", err.Error()),
			),
		)

		return []map[string]interface{}{}, ErrFailedParsingTemplate
	}

	parser := gotextfsm.ParserOutput{}

	err = parser.ParseTextString(r.Result, fsm, true)
	if err != nil {
		logging.LogError(
			r.FormatLogMessage(
				"error",
				fmt.Sprintf(
					"Error while parsing device output, gotextfsm error: %s\n",
					err.Error(),
				),
			),
		)

		return []map[string]interface{}{}, err
	}

	return parser.Dict, nil
}

// FormatLogMessage formats log message payload, adding contextual info about the host.
func (r *Response) FormatLogMessage(level, msg string) string {
	return logging.FormatLogMessage(level, r.Host, r.Port, msg)
}
