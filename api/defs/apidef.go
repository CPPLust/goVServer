package defs

type SignedUp struct {
	Success bool   `json:"success"`
	Session string `json:"session_id"`
}

// request
type UserCredential struct {
	Username string `json:"username"`
	Pwd      string `json:"passwd"`
}

type VideoInfo struct {
	Id           string
	AuthorId     int
	Name         string
	DisplayCtime string
}

type Comment struct {
	Id      string
	VideoId string
	Author  string
	Content string
}

type SimpleSession struct {
	Username string
	TTL      int64
}
