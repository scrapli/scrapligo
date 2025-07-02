package util //nolint: revive

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	"github.com/sirikothe/gotextfsm"
)

// ResolveAtFileOrURL returns the bytes from `path` where path is either a filepath or URL.
func ResolveAtFileOrURL(path string) ([]byte, error) {
	var b []byte

	switch {
	case strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://"):
		resp, err := http.Get(path) //nolint:gosec,noctx
		if err != nil {
			return nil, scrapligoerrors.NewUtilError(
				fmt.Sprintf("failed downloading file at path %q", path),
				err,
			)
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, scrapligoerrors.NewUtilError(
				fmt.Sprintf("failed reading downloaded file at path %q", path),
				err,
			)
		}

	default: // fall-through to local filesystem
		var err error

		b, err = os.ReadFile(path) //nolint:gosec
		if err != nil {
			return nil, scrapligoerrors.NewUtilError(
				fmt.Sprintf("failed opening provided file at path %q", path),
				err,
			)
		}
	}

	return b, nil
}

// TextFsmParse parses recorded output w/ a provided textfsm template.
// the argument is interpreted as URL or filesystem path, for example:
// response.TextFsmParse("http://example.com/textfsm.template") or
// response.TextFsmParse("./local/textfsm.template").
func TextFsmParse(s, path string) ([]map[string]interface{}, error) {
	t, err := ResolveAtFileOrURL(path)
	if err != nil {
		return []map[string]interface{}{}, err
	}

	fsm := gotextfsm.TextFSM{}

	err = fsm.ParseString(string(t))
	if err != nil {
		return nil, scrapligoerrors.NewUtilError("failed parsing provided template", err)
	}

	parser := gotextfsm.ParserOutput{}

	err = parser.ParseTextString(s, fsm, true)
	if err != nil {
		return nil, scrapligoerrors.NewUtilError("failed parsing device output", err)
	}

	return parser.Dict, nil
}
