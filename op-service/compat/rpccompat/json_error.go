package rpccompat

import "fmt"

// JsonError is a compatibility shim for rpc.JsonError which is unexported
// in older geth versions (as jsonError).
type JsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *JsonError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("json-rpc error %d", e.Code)
	}
	return e.Message
}

func (e *JsonError) ErrorCode() int {
	return e.Code
}

func (e *JsonError) ErrorData() any {
	return e.Data
}
