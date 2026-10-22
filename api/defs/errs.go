package defs

import "net/http"

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}
type ErrorResponse struct {
	HttpSC int
	Error  Err
}

var (
	ErrorRequestBodyParseFailed = ErrorResponse{HttpSC: 400, Error: Err{Error: "request body parse failed", ErrorCode: "001"}}
	ErrorNotAuthUser            = ErrorResponse{HttpSC: 401, Error: Err{Error: "user not authentic", ErrorCode: "002"}}
	ErrorDBError                = ErrorResponse{HttpSC: 500, Error: Err{Error: "db error", ErrorCode: "003"}}
	ErrorInternalFaults         = ErrorResponse{HttpSC: 500, Error: Err{Error: "internal service error", ErrorCode: "004"}}
	ErrorServerReachLimit       = ErrorResponse{HttpSC: http.StatusServiceUnavailable, Error: Err{Error: "connection limit reached", ErrorCode: "005"}}
	ErrorBadRequest             = ErrorResponse{HttpSC: http.StatusBadRequest, Error: Err{Error: "File is too big", ErrorCode: "006"}}
)
