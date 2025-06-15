package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/ebitengine/purego"
	scrapligoconstants "github.com/scrapli/scrapligo/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/errors"
)

// ReadWithCallbacks optionally sends an initial "input" to the device and then continually reads
// from the session, checking new session output against the provided callbacks (in the order
// provided). When a callback is triggered, the callback is executed, if the callback is marked as
// "completes" then the parent function exits, otherwise this continues forever.
func (c *Cli) ReadWithCallbacks(
	ctx context.Context,
	options ...OperationOption,
) (*Result, error) {
	if c.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	names := strings.Join([]string{"foo", "bar"}, scrapligoconstants.LibScrapliDelimiter)
	contains := strings.Join([]string{"x86_64", "x86_64"}, scrapligoconstants.LibScrapliDelimiter)
	containsPatterns := strings.Join([]string{"", ""}, scrapligoconstants.LibScrapliDelimiter)
	notContains := strings.Join([]string{"", ""}, scrapligoconstants.LibScrapliDelimiter)
	once := strings.Join([]string{"true", "true"}, scrapligoconstants.LibScrapliDelimiter)
	resetTimer := strings.Join([]string{"", ""}, scrapligoconstants.LibScrapliDelimiter)
	completes := strings.Join([]string{"false", "true"}, scrapligoconstants.LibScrapliDelimiter)

	status := c.ffiMap.Cli.ReadWithCallbacks(
		c.ptr,
		&operationID,
		&cancel,
		"show version",
		names,
		[]uintptr{
			purego.NewCallback(func() uint8 {
				fmt.Println("FOOOOOP!?")
				_, err := c.SendInput(ctx, "show version")
				fmt.Println("SEND IT")
				if err != nil {
					return 1
				}

				return 0
			}),
			purego.NewCallback(func() uint8 {
				fmt.Println("CB 1!")
				return 0
			}),
		},
		contains,
		containsPatterns,
		notContains,
		once,
		resetTimer,
		completes,
	)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit readWithCallbacks operation", nil)
	}

	return c.getResult(ctx, &cancel, operationID)
}
