package main

import (
	"goVServer/api/session"
	"net/http"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_UNAME = "X-User-Name"

func validateUserSession(request *http.Request) bool {
	sid := request.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}
	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}
	request.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}
