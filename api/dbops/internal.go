package dbops

import (
	"fmt"
	"goVServer/api/defs"
	"strconv"
	"sync"
)

func InsertSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmtIns, err := dbConn.Prepare("INSERT INTO sessions (session_id, TTL, login_name) VALUES (?, ?,?)")
	if err != nil {
		fmt.Printf("stmtIns err:%v", err)
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(sid, ttlstr, uname)
	if err != nil {
		fmt.Printf("stmtIns err:%v", err)
		return err
	}
	return nil
}
func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmtOut, err := dbConn.Prepare(`SELECT TTL, login_name FROM sessions WHERE session_id = ?`)
	if err != nil {
		fmt.Printf("stmtOut err:%v", err)
		return nil, err
	}
	defer stmtOut.Close()
	var ttl string
	var uname string
	err = stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil {
		fmt.Printf("stmtOut err:%v", err)
		return nil, err

	}
	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = res
		ss.Username = uname
	} else {
		fmt.Printf("stmtOut err:%v", err)
		return nil, err
	}

	return ss, nil
}

func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare("SELECT * FROM sessions")
	if err != nil {
		fmt.Printf("stmtOut err:%v", err)
		return nil, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query()
	if err != nil {
		fmt.Printf("stmtOut err:%v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sid string
		var ttl string
		var uname string
		if err := rows.Scan(&sid, &ttl, &uname); err != nil {
			fmt.Printf("rows.Scan err:%v", err)
			break
		}
		if ttl, err := strconv.ParseInt(ttl, 10, 64); err == nil {

			ss := &defs.SimpleSession{Username: uname, TTL: ttl}
			m.Store(sid, ss)
			fmt.Printf("session:%v, ttl:%v", ss, ttl)
		}
	}
	return m, nil
}

func DeleteSession(sid string) error {
	stmOut, err := dbConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		fmt.Printf("stmtOut err:%v", err)
		return err
	}
	defer stmOut.Close()
	_, err = stmOut.Exec(sid)
	if err != nil {
		fmt.Printf("stmtOut err:%v", err)
		return err
	}
	return nil
}
