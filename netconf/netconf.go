package netconf

import (
	"context"
	"errors"

	scrapligoerrors "github.com/scrapli/scrapligo/errors"
	scrapligoffi "github.com/scrapli/scrapligo/ffi"
	scrapligointernal "github.com/scrapli/scrapligo/internal"
	scrapligologging "github.com/scrapli/scrapligo/logging"
	scrapligooptions "github.com/scrapli/scrapligo/options"
	"golang.org/x/sys/unix"
)

func newCloseOptions(options ...Option) *closeOptions {
	o := &closeOptions{}

	for _, opt := range options {
		opt(o)
	}

	return o
}

type closeOptions struct {
	force bool
}

// Netconf is an object representing a netconf connection to a device of some sort -- this object
// wraps the underlying zig (netconf) driver (created via libscrapli).
type Netconf struct {
	ptr     uintptr
	pollFd  int
	ffiMap  *scrapligoffi.Mapping
	host    string
	options *scrapligointernal.Options
}

// NewNetconf returns a new instance of Netconf setup with the given options.
func NewNetconf(
	host string,
	opts ...scrapligooptions.Option,
) (*Netconf, error) {
	ffiMap, err := scrapligoffi.GetMapping()
	if err != nil {
		return nil, err
	}

	n := &Netconf{
		ffiMap:  ffiMap,
		host:    host,
		options: scrapligointernal.NewOptions(),
	}

	for _, opt := range opts {
		err = opt(n.options)
		if err != nil {
			return nil, scrapligoerrors.NewOptionsError("failed applying option", err)
		}
	}

	if n.options.Port == 0 {
		n.options.Port = 830
	}

	return n, nil
}

// GetPtr returns the pointer to the zig driver, don't use this unless you know what you are doing,
// this is just exposed so you *can* get to it if you want to.
func (n *Netconf) GetPtr() (uintptr, *scrapligoffi.Mapping) {
	return n.ptr, n.ffiMap
}

// Open opens the driver object. This method spawns the underlying zig driver which the Netconf
// object then holds a pointer to. All Netconf operations operate against this pointer (though
// this is transparent to the user).
func (n *Netconf) Open(ctx context.Context) (*Result, error) {
	// ensure we dealloc if something happens, otherwise users calls to defer close would not be
	// super handy
	cleanup := false

	defer func() {
		if !cleanup {
			return
		}

		n.ffiMap.Shared.Free(n.ptr)
	}()

	n.ptr = n.ffiMap.Netconf.Alloc(
		scrapligologging.LoggerToLoggerCallback(
			n.options.Logger,
			uint8(scrapligologging.IntFromLevel(n.options.LoggerLevel)),
		),
		string(n.options.LoggerLevel),
		n.host,
		n.options.Port,
		string(n.options.TransportKind),
	)

	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("failed to allocate netconf", nil)
	}

	n.pollFd = int(n.ffiMap.Shared.GetPollFd(n.ptr))
	if n.pollFd == 0 {
		return nil, scrapligoerrors.NewFfiError("failed to allocate netconf", nil)
	}

	err := n.options.Apply(n.ptr, n.ffiMap)
	if err != nil {
		return nil, scrapligoerrors.NewFfiError("failed to applying netconf options", err)
	}

	cancel := false

	var operationID uint32

	status := n.ffiMap.Netconf.Open(n.ptr, &operationID, &cancel)
	if status != 0 {
		cleanup = true

		return nil, scrapligoerrors.NewFfiError("failed to submit open operation", nil)
	}

	result, err := n.getResult(ctx, &cancel, operationID)
	if err != nil {
		cleanup = true

		return nil, err
	}

	return result, nil
}

// Close closes the netconf object. This also deallocates the underlying (zig) netconf object.
func (n *Netconf) Close(ctx context.Context, options ...Option) (*Result, error) {
	if n.ptr == 0 {
		return nil, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	cancel := false

	var operationID uint32

	loadedOptions := newCloseOptions(options...)

	status := n.ffiMap.Netconf.Close(n.ptr, &operationID, &cancel, loadedOptions.force)
	if status != 0 {
		return nil, scrapligoerrors.NewFfiError("failed to submit close operation", nil)
	}

	result, err := n.getResult(ctx, &cancel, operationID)

	n.ffiMap.Shared.Free(n.ptr)

	return result, err
}

// GetSessionID returns the session-id as parsed during the capabilities exchange -- if we for some
// reason didn't parse the session-id during capabilities exchange this will return an error.
func (n *Netconf) GetSessionID() (uint64, error) {
	if n.ptr == 0 {
		return 0, scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	var sessionID uint64

	status := n.ffiMap.Netconf.GetSessionID(n.ptr, &sessionID)
	if status != 0 {
		return 0, scrapligoerrors.NewFfiError("session-id not set", nil)
	}

	return sessionID, nil
}

// GetSubscriptionID attempts to parse a subscription id from a netconf rpc reply message. It can
// parse a response from content like:
// <?xml version="1.0" encoding="UTF-8"?>
// <rpc-reply message-id="101">
// <subscription-result>ok</subscription-result>
// <subscription-id>2147483737</subscription-id>
// </rpc-reply>
// In case of failure, a ffi error will wrap an ErrSubscriptionId error.
func (n *Netconf) GetSubscriptionID(message string) (uint64, error) {
	var subscriptionID uint64

	status := n.ffiMap.Netconf.GetSubscriptionID(message, &subscriptionID)
	if status != 0 {
		return 0, scrapligoerrors.NewFfiError(
			"failed parsing subscription id",
			scrapligoerrors.ErrSubscriptionID,
		)
	}

	return subscriptionID, nil
}

// GetNextNotification returns the next notification type message, if any. If there are no messages,
// a ErrNoMessages will be returned.
func (n *Netconf) GetNextNotification() (string, error) {
	if n.ptr == 0 {
		return "", scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	var notifSize uint64

	n.ffiMap.Netconf.GetNextNotificationSize(n.ptr, &notifSize)

	if notifSize == 0 {
		return "", scrapligoerrors.NewMessagesError()
	}

	notif := make([]byte, notifSize)

	rc := n.ffiMap.Netconf.GetNextNotification(n.ptr, &notif)
	if rc != 0 {
		return "", scrapligoerrors.NewFfiError("get next notification failed", nil)
	}

	return string(notif), nil
}

// GetNextSubscription returns the next subscription type message for the given subscription id,
// if any. If there are no messages, a ErrNoMessages will be returned.
func (n *Netconf) GetNextSubscription(subscriptionID uint64) (string, error) {
	if n.ptr == 0 {
		return "", scrapligoerrors.NewFfiError("driver pointer nil", nil)
	}

	var subSize uint64

	n.ffiMap.Netconf.GetNextSubscriptionSize(n.ptr, subscriptionID, &subSize)

	if subSize == 0 {
		return "", scrapligoerrors.NewMessagesError()
	}

	sub := make([]byte, subSize)

	rc := n.ffiMap.Netconf.GetNextSubscription(n.ptr, subscriptionID, &sub)
	if rc != 0 {
		return "", scrapligoerrors.NewFfiError("get next subscription failed", nil)
	}

	return string(sub), nil
}

func (n *Netconf) getResult(
	ctx context.Context,
	cancel *bool,
	operationID uint32,
) (*Result, error) {
	done := make(chan struct{}, 1)
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			*cancel = true

			return
		case <-done:
			return
		}
	}()

	pollFd := &unix.FdSet{}
	pollFd.Set(n.pollFd)

	var _n int

	for {
		var err error

		_n, err = unix.Select(n.pollFd+1, pollFd, &unix.FdSet{}, &unix.FdSet{}, nil)
		if err != nil {
			if errors.Is(err, unix.EINTR) {
				// python automagically handles interrupts i guess go doesnt, so just act like
				// we do on the python side when polling the wakeup fd
				continue
			}

			return nil, scrapligoerrors.NewFfiError("waiting on operation ready signal", err)
		}

		break
	}

	// if the context wasn't cancelled the goroutine will still be running, this will stop it
	done <- struct{}{}

	out := make([]byte, _n)

	_, _ = unix.Read(n.pollFd, out)

	var inputSize, resultRawSize, resultSize, rpcWarningsSize, rpcErrorsSize, errSize uint64

	rc := n.ffiMap.Netconf.FetchOperationSizes(
		n.ptr,
		operationID,
		&inputSize,
		&resultRawSize,
		&resultSize,
		&rpcWarningsSize,
		&rpcErrorsSize,
		&errSize,
	)
	if rc != 0 {
		return nil, scrapligoerrors.NewFfiError("poll operation failed", nil)
	}

	var resultStartTime, resultEndTime uint64

	input := make([]byte, inputSize)

	resultRaw := make([]byte, resultRawSize)

	result := make([]byte, resultSize)

	rpcWarnings := make([]byte, rpcWarningsSize)

	rpcErrors := make([]byte, rpcErrorsSize)

	errString := make([]byte, errSize)

	rc = n.ffiMap.Netconf.FetchOperation(
		n.ptr,
		operationID,
		&resultStartTime,
		&resultEndTime,
		&input,
		&resultRaw,
		&result,
		&rpcWarnings,
		&rpcErrors,
		&errString,
	)
	if rc != 0 {
		return nil, scrapligoerrors.NewFfiError("fetch operation result failed", nil)
	}

	if errSize != 0 {
		return nil, scrapligoerrors.NewFfiError(string(errString), nil)
	}

	return NewResult(
		string(input),
		n.host,
		n.options.Port,
		resultStartTime,
		resultEndTime,
		resultRaw,
		string(result),
		rpcWarnings,
		rpcErrors,
	), nil
}
