package netconf

// NewResult prepares a new Result object from ffi integration pointers (the pointers we pass to
// zig for it to populate the values of stuff).
func NewResult(
	input,
	host string,
	port uint16,
	startTime uint64,
	endTime uint64,
	resultRaw []byte,
	result string,
) *Result {
	return &Result{
		Host:      host,
		Port:      port,
		Input:     input,
		ResultRaw: resultRaw,
		Result:    result,
		StartTime: startTime,
		EndTime:   endTime,
	}
}

// Result is a struct returned from all Driver operations.
type Result struct {
	Host               string
	Port               uint16
	Input              string
	ResultRaw          []byte
	Result             string
	StartTime          uint64
	EndTime            uint64
	ElapsedTimeSeconds float64
	Failed             bool
}
