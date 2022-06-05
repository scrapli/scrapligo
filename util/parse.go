package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sirikothe/gotextfsm"
)

// TextFsmParse parses recorded output w/ a provided textfsm template.
// the argument is interpreted as URL or filesystem path, for example:
// response.TextFsmParse("http://example.com/textfsm.template") or
// response.TextFsmParse("./local/textfsm.template").
func TextFsmParse(s, path string) ([]map[string]interface{}, error) {
	var t []byte

	switch {
	case strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://"):
		resp, err := http.Get(path) // nolint:gosec,noctx
		if err != nil {
			return []map[string]interface{}{}, fmt.Errorf(
				"%w: failed downloading template, error: %s",
				ErrParseError,
				err,
			)
		}

		defer resp.Body.Close() //nolint:errcheck

		t, err = io.ReadAll(resp.Body)
		if err != nil {
			return []map[string]interface{}{}, fmt.Errorf(
				"%w: failed reading downloaded template, error: %s",
				ErrParseError,
				err,
			)
		}

	default: // fall-through to local filesystem
		var err error

		t, err = os.ReadFile(path) //nolint:gosec
		if err != nil {
			return []map[string]interface{}{}, fmt.Errorf(
				"%w: failed opening provided template, error: %s",
				ErrParseError,
				err,
			)
		}
	}

	fsm := gotextfsm.TextFSM{}

	err := fsm.ParseString(string(t))
	if err != nil {
		return []map[string]interface{}{}, fmt.Errorf(
			"%w: failed parsing provided template, gotextfsm error: %s",
			ErrParseError,
			err,
		)
	}

	parser := gotextfsm.ParserOutput{}

	err = parser.ParseTextString(s, fsm, true)
	if err != nil {
		return []map[string]interface{}{}, fmt.Errorf(
			"%w: failed parsing device output, gotextfsm error: %s",
			ErrParseError,
			err,
		)
	}

	return parser.Dict, nil
}
