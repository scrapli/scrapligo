package internal

import (
	"context"
	"encoding/binary"
	"errors"
	"sync"

	scrapligoconstants "github.com/scrapli/scrapligo/v2/constants"
	scrapligoerrors "github.com/scrapli/scrapligo/v2/errors"
	"golang.org/x/sys/unix"
)

// GetResultWaitWakeup shared func to wait for wakeup signals on result fetching for Cli and
// Netconf objects.
func GetResultWaitWakeup( //nolint: gocyclo
	ctx context.Context,
	pollFd int,
	cancelLock *sync.Mutex,
	cancel *bool,
	operationID uint32,
) error {
	var n int

	pollFds := []unix.PollFd{{Fd: int32(pollFd), Events: unix.POLLIN}} //nolint: gosec

	for {
		if ctx.Err() != nil {
			cancelLock.Lock()

			*cancel = true

			cancelLock.Unlock()

			return ctx.Err()
		}

		pollFds[0].Revents = 0

		var err error

		n, err = unix.Poll(pollFds, scrapligoconstants.ReadyFDPollTimeoutMs)
		if err != nil {
			if errors.Is(err, unix.EINTR) {
				// python automagically handles interrupts i guess go doesnt, so just act like
				// we do on the python side when polling the wakeup fd
				continue
			}

			return scrapligoerrors.NewFfiError("waiting on operation ready signal", err)
		}

		if n > 0 {
			if pollFds[0].Revents&unix.POLLNVAL != 0 {
				return scrapligoerrors.NewFfiError(
					"waiting on operation ready signal",
					unix.EBADF,
				)
			}

			break
		}
	}

	var out [4]byte

	for {
		if ctx.Err() != nil {
			cancelLock.Lock()

			*cancel = true

			cancelLock.Unlock()

			return ctx.Err()
		}

		n, err := unix.Read(pollFd, out[:])
		if err != nil {
			if errors.Is(err, unix.EINTR) {
				// same as loop above -- retry on interrupts
				continue
			}

			return scrapligoerrors.NewFfiError("draining operation ready signal", err)
		}

		switch n {
		case 4: //nolint: mnd
			// happy path, nthoing to do
		case 0:
			return scrapligoerrors.NewFfiError("wakeup signal fd closed", nil)
		default:
			return scrapligoerrors.NewFfiError("wakeup signal read invalid read size", nil)
		}

		waitWakeupOperationID := binary.LittleEndian.Uint32(out[:])
		switch {
		case waitWakeupOperationID == operationID:
			return nil
		case waitWakeupOperationID < operationID:
			// we snagged a wakeup for a stale (probably cancelled op)
			return GetResultWaitWakeup(ctx, pollFd, cancelLock, cancel, operationID)
		case waitWakeupOperationID > operationID:
			return scrapligoerrors.NewFfiError(
				"wakeup signal received for operation id greater than requested, "+
					"this should not happen",
				nil,
			)
		}
	}
}
