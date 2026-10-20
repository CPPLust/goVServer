package session

import (
	"goVServer/api/dbops"
	"goVServer/api/defs"
	"goVServer/api/utils"
	"sync"
	"time"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1e6
}
func deleteExpiredSessions(sid string) {
	sessionMap.Delete(sid)
	_ = dbops.DeleteSession(sid)
}

func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}
	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})
}

func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000
	ss := &defs.SimpleSession{TTL: ttl, Username: un}
	sessionMap.Store(id, ss)
	_ = dbops.InsertSession(id, ttl, un)
	return id
}

func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := nowInMilli()
		if ss.(*defs.SimpleSession).TTL < ct {
			deleteExpiredSessions(sid)
			return "", true
		}
		return ss.(*defs.SimpleSession).Username, false
	}

	return "", true
}
