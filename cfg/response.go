package cfg

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/scrapli/scrapligo/difflibgo"
	"github.com/scrapli/scrapligo/driver/base"
)

// Response cfg response object that gets returned from cfg operations.
type Response struct {
	Host             string
	Result           string
	StartTime        time.Time
	EndTime          time.Time
	ElapsedTime      float64
	ScrapliResponses []*base.Response
	ErrorType        error
	Failed           bool
}

// NewResponse create a new cfg response object.
func NewResponse(
	host string,
	raiseError error,
) *Response {
	r := &Response{
		Host:        host,
		Result:      "",
		StartTime:   time.Now(),
		EndTime:     time.Time{},
		ElapsedTime: 0,
		ErrorType:   raiseError,
		Failed:      true,
	}

	return r
}

func (r *Response) Record(scrapliResponses []*base.Response, result string) {
	r.EndTime = time.Now()
	r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()

	r.Result = result
	r.Failed = false
	r.ScrapliResponses = scrapliResponses

	for _, response := range r.ScrapliResponses {
		if response.Failed {
			r.Failed = true
			break
		}
	}
}

// DiffResponse cfg diff response object that gets returned from diff operations.
type DiffResponse struct {
	*Response
	Source          string
	SourceConfig    string
	CandidateConfig string
	DeviceDiff      string
	difflines       []string
	additions       []string
	subtractions    []string
	sideBySideDiff  string
	colorize        bool
	sideBySideWidth int
}

// NewDiffResponse create a new cfg diff response object.
func NewDiffResponse(
	host string,
	source string,
	colorize bool,
	sideBySideWidth int,
) *DiffResponse {
	r := &Response{
		Host:        host,
		Result:      "",
		StartTime:   time.Now(),
		EndTime:     time.Time{},
		ElapsedTime: 0,
		ErrorType:   ErrDiffConfigFailed,
		Failed:      true,
	}

	dr := &DiffResponse{
		Response:        r,
		Source:          source,
		colorize:        colorize,
		sideBySideWidth: sideBySideWidth,
	}

	return dr
}

func (r *DiffResponse) RecordDiff(sourceConfig, candidateConfig, deviceDiff string) {
	r.SourceConfig = sourceConfig
	r.CandidateConfig = candidateConfig
	r.DeviceDiff = deviceDiff

	d := difflibgo.Differ{}
	r.difflines = d.Compare(
		strings.Split(r.SourceConfig, "\n"),
		strings.Split(r.CandidateConfig, "\n"),
	)

	for _, v := range r.difflines {
		if v[:2] == "+ " {
			r.additions = append(r.additions, v[2:])
		} else if v[:2] == "- " {
			r.subtractions = append(r.additions, v[2:])
		}
	}
}

func (r *DiffResponse) generateColors() (string, string, string, string) {
	if !r.colorize {
		return "? ", "- ", "+ ", ""
	}

	return "\033[93m", "\033[91m", "\033[92m", "\033[0m"
}

func (r *DiffResponse) terminalWidth() int {
	if r.sideBySideWidth != 0 {
		return r.sideBySideWidth
	}

	width, _, _ := term.GetSize(0)

	return width
}

func (r *DiffResponse) SideBySideDiff() string {
	if len(r.sideBySideDiff) > 0 {
		return r.sideBySideDiff
	}

	yellow, red, green, end := r.generateColors()

	termWidth := r.terminalWidth()

	halfTermWidth := termWidth / 2
	diffSideWidth := halfTermWidth - 5

	var sideBySideDiffLines []string

	for _, line := range r.difflines {
		var current, candidate string

		trimLen := diffSideWidth
		difflineLen := len(line)
		if difflineLen < trimLen {
			trimLen = difflineLen - 2
		}

		switch line[:2] {
		case " ?":
			current = yellow + fmt.Sprintf(
				"%-*s",
				halfTermWidth,
				strings.TrimRight(line[2:][:trimLen], " "),
			) + end
			candidate = yellow + strings.TrimRight(line[2:][:trimLen], " ") + end
		case "- ":
			current = red + fmt.Sprintf(
				"%-*s",
				halfTermWidth,
				strings.TrimRight(line[2:][:trimLen], " "),
			) + end
			candidate = ""
		case "+ ":
			current = strings.Repeat(" ", halfTermWidth)
			candidate = green + strings.TrimRight(line[2:][:trimLen], " ") + end
		default:
			current = fmt.Sprintf(
				"%-*s",
				halfTermWidth,
				strings.TrimRight(line[2:][:trimLen], " "),
			)
			candidate = strings.TrimRight(line[2:][:trimLen], " ")
		}

		sideBySideDiffLines = append(sideBySideDiffLines, current+candidate)
	}

	r.sideBySideDiff = strings.Join(sideBySideDiffLines, "\n")

	return r.sideBySideDiff
}
