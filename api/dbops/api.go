package dbops

import (
	"database/sql"
	"fmt"
	"goVServer/api/defs"
	"goVServer/api/utils"
	"log"
	"time"
)

func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (login_name,pwd) VALUES (?,?)")
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
func GetUserCredential(loginName string) (string, error) {
	stmOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name=?")
	if err != nil {
		return "", err
	}
	defer stmOut.Close()

	var pwd string
	err = stmOut.QueryRow(loginName).Scan(&pwd)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return pwd, nil
}

// UserExists 判断用户名是否已存在
func UserExists(loginName string) (bool, error) {
	stmtOut, err := dbConn.Prepare("SELECT 1 FROM users WHERE login_name=?")
	if err != nil {
		return false, err
	}
	defer stmtOut.Close()

	var exists int
	err = stmtOut.QueryRow(loginName).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmOut, err := dbConn.Prepare("DELETE FROM users WHERE login_name=? AND pwd =?")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmOut.Close()
	_, err = stmOut.Exec(loginName, pwd)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// GetUser 根据用户名获取用户信息
func GetUser(loginName string) (*defs.User, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	var id int
	var pwd string

	err = stmtOut.QueryRow(loginName).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	res := &defs.User{Id: id, LoginName: loginName, Pwd: pwd}

	defer stmtOut.Close()

	return res, nil
}

func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	//create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	t := time.Now()
	ctime := t.Format(time.DateTime)
	stmIns, err := dbConn.Prepare(`INSERT INTO video_info (id,author_id,name,display_ctime,create_time) VALUES(?,?,?,?,?)`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer stmIns.Close()
	_, err = stmIns.Exec(vid, aid, name, ctime, t)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}
	return res, nil
}

func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=?")

	var aid int
	var dct string
	var name string

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}

	return res, nil
}

func ListVideoInfo(uname string, from, to int) ([]*defs.VideoInfo, error) {
	fmt.Printf("[ListVideoInfo] uname=%s, from=%d, to=%d\n", uname, from, to)

	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.author_id, video_info.name, video_info.display_ctime FROM video_info 
		INNER JOIN users ON video_info.author_id = users.id
		WHERE users.login_name = ? AND video_info.create_time >= FROM_UNIXTIME(?) AND video_info.create_time <= FROM_UNIXTIME(?) 
		ORDER BY video_info.create_time DESC`)

	var res []*defs.VideoInfo

	if err != nil {
		fmt.Printf("[ListVideoInfo] Prepare error: %v\n", err)
		return res, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(uname, from, to)
	if err != nil {
		fmt.Printf("[ListVideoInfo] Query error: %v\n", err)
		return res, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, name, ctime string
		var aid int
		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			fmt.Printf("[ListVideoInfo] Scan error: %v\n", err)
			return res, err
		}
		fmt.Printf("[ListVideoInfo] Found video: id=%s, aid=%d, name=%s, ctime=%s\n", id, aid, name, ctime)

		vi := &defs.VideoInfo{Id: id, AuthorId: aid, Name: name, DisplayCtime: ctime}
		res = append(res, vi)
		count++
	}

	fmt.Printf("[ListVideoInfo] Total videos found: %d\n", count)
	return res, nil
}

func DeleteVideoInfo(vid string) error {
	stmOut, err := dbConn.Prepare("DELETE FROM video_info WHERE id=?")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmOut.Close()
	_, err = stmOut.Exec(vid)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		fmt.Println(err)
		return err
	}

	/*
		CREATE TABLE IF NOT EXISTS comments (
		  id VARCHAR(64) PRIMARY KEY NOT NULL,
		  video_id VARCHAR(64),
		  author_id INT,
		  content TEXT,
		  time DATETIME
		);
	*/
	stmIns, err := dbConn.Prepare("INSERT INTO comments (id,video_id,author_id,content,time) VALUES(?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmIns.Close()
	_, err = stmIns.Exec(id, vid, aid, content, time.Now())
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(` SELECT comments.id, users.login_name, comments.content FROM comments
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.video_id = ? AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)`)

	var res []*defs.Comment

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}

		c := &defs.Comment{Id: id, VideoId: vid, Author: name, Content: content}
		res = append(res, c)
	}
	defer stmtOut.Close()

	return res, nil
	/*
				INNER JOIN 是 SQL 中的内连接操作，用于根据两个表之间的关联条件，返回两个表中都匹配的行
				comments 表:
		| id | video_id | author_id | content |
		|----|----------|-----------|---------|
		| 1  | v001     | 100       | 好看！   |
		| 2  | v001     | 101       | 不错~   |
		| 3  | v001     | 999       | 一般    |  ← 这个用户ID在users表中不存在

		users 表:
		| id  | login_name |
		|-----|------------|
		| 100 | 张三       |
		| 101 | 李四       |
		使用 INNER JOIN
		SELECT comments.id, users.login_name, comments.content
		FROM comments
		INNER JOIN users ON comments.author_id = users.id
		结果：
		| id | login_name | content |
		|----|------------|---------|
		| 1  | 张三       | 好看！   |
		| 2  | 李四       | 不错~   |
		⚠️ 注意：评论ID=3 的记录被过滤掉了，因为 author_id=999 在 users 表中没有对应的用户。
	*/
}
