package cfg

import (
	"fmt"
	"strings"
	"time"

	"github.com/scrapli/scrapligo/driver/base"

	"golang.org/x/term"

	"github.com/carlmontanari/difflibgo/difflibgo"
)

const (
	subtraction = "- "
	addition    = "+ "
	unknown     = "? "
)

// Response cfg response object that gets returned from cfg operations.
type Response struct {
	Host             string
	OperationType    string
	Result           string
	StartTime        time.Time
	EndTime          time.Time
	ElapsedTime      float64
	ScrapliResponses []*base.Response
	ErrorType        error
	Failed           error
}

// NewResponse create a new cfg response object.
func NewResponse(
	host,
	opType string,
	raiseError error,
) *Response {
	r := &Response{
		Host:          host,
		OperationType: opType,
		Result:        "",
		StartTime:     time.Now(),
		EndTime:       time.Time{},
		ElapsedTime:   0,
		ErrorType:     raiseError,
	}

	return r
}

func (r *Response) Record(scrapliResponses []*base.Response, result string) {
	r.EndTime = time.Now()
	r.ElapsedTime = r.EndTime.Sub(r.StartTime).Seconds()

	r.Result = result
	r.ScrapliResponses = scrapliResponses

	for _, resp := range r.ScrapliResponses {
		if resp.Failed != nil {
			r.Failed = r.ErrorType
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
	unifiedDiff     string
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
		Host:          host,
		OperationType: "DiffConfig",
		Result:        "",
		StartTime:     time.Now(),
		EndTime:       time.Time{},
		ElapsedTime:   0,
		ErrorType:     ErrDiffConfigFailed,
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
		if v[:2] == addition {
			r.additions = append(r.additions, v[2:])
		} else if v[:2] == subtraction {
			r.subtractions = append(r.subtractions, v[2:])
		}
	}
}

func (r *DiffResponse) generateColors() (unknown, subtraction, addition, end string) {
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

	sideBySideDiffLines := make([]string, 0)

	for _, line := range r.difflines {
		var current, candidate string

		trimLen := diffSideWidth
		difflineLen := len(line)

		if difflineLen-1 <= trimLen {
			trimLen = difflineLen - 2
		}

		switch line[:2] {
		case unknown:
			current = yellow + fmt.Sprintf(
				"%-*s",
				halfTermWidth,
				strings.TrimRight(line[2:][:trimLen], " "),
			) + end
			candidate = yellow + strings.TrimRight(line[2:][:trimLen], " ") + end
		case subtraction:
			current = red + fmt.Sprintf(
				"%-*s",
				halfTermWidth,
				strings.TrimRight(line[2:][:trimLen], " "),
			) + end
			candidate = ""
		case addition:
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

func (r *DiffResponse) UnifiedDiff() string {
	if len(r.unifiedDiff) > 0 {
		return r.unifiedDiff
	}

	yellow, red, green, end := r.generateColors()

	unifiedDiffLines := make([]string, 0)

	for _, line := range r.difflines {
		var diffLine string

		switch line[:2] {
		case unknown:
			diffLine = yellow + line[2:] + end
		case subtraction:
			diffLine = red + line[2:] + end
		case addition:
			diffLine = green + line[2:] + end
		default:
			diffLine = line[2:]
		}

		unifiedDiffLines = append(unifiedDiffLines, diffLine)
	}

	r.unifiedDiff = strings.Join(unifiedDiffLines, "\n")

	return r.unifiedDiff
}
