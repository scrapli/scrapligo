package base

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/scrapli/scrapligo/util"

	"github.com/scrapli/scrapligo/logging"
	"github.com/sirikothe/gotextfsm"
)

// ErrFailedOpeningTemplate error for failure to open a textfsm template.
var ErrFailedOpeningTemplate = errors.New("failed opening provided path to textfsm template")

// ErrFailedParsingTemplate error for failure of parsing a textfsm template.
var ErrFailedParsingTemplate = errors.New("failed parsing textfsm template")

// OperationError is an error object returned when a scrapli operation completes "successfully" --
// as in does not have an EOF/timeout or otherwise unrecoverable error -- but contains output in the
// device's response indicating that an input was bad/invalid or device failed to process it at
// that time.
type OperationError struct {
	Input       string
	Output      string
	ErrorString string
}

// Error returns an error string for the OperationError object.
func (e *OperationError) Error() string {
	return fmt.Sprintf(
		"operation error from input '%s'. indicated error '%s'",
		e.Input,
		e.ErrorString,
	)
}

// Response is a response object that gets returned from scrapli send operations.
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
	// Failed returns an error if any of the `FailedWhenContains` substrings are seen in the output
	// returned from the device. This error indicates that the operation has completed successfully,
	// but that an input was bad/invalid or device failed to process it at that time
	Failed error
}

// NewResponse creates a new response object.
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

	s := util.StrContainsAnySubStr(r.Result, r.FailedWhenContains)
	if len(s) > 0 {
		r.Failed = &OperationError{
			Input:       r.ChannelInput,
			Output:      r.Result,
			ErrorString: s,
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
