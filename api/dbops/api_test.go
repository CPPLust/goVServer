package dbops

import (
	"fmt"
	"testing"
)

func clearTables() {
	_, err := dbConn.Exec("truncate users")
	if err != nil {
		fmt.Println(err)
	}
	_, err = dbConn.Exec("truncate video_info")
	if err != nil {
		fmt.Println(err)
	}
}

func TestMain(m *testing.M) {

	m.Run()

}

func TestUserWorkFlow(t *testing.T) {
	clearTables()
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Delete", testDeleteUser)
	t.Run("ReGet", testReGetUser)
	clearTables()
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("long", "123")
	if err != nil {
		t.Error(err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("long")
	if err != nil {
		t.Errorf("testGetUser err:%v", err)
	}

	if pwd != "123" {
		t.Errorf("testGetUser pwd:%v", pwd)
	}
}

func testDeleteUser(t *testing.T) {
	//err := DeleteUser("lon", "123")
	//if err == nil {
	//	t.Errorf("testDeleteUser 应该不为空")
	//}
	err = DeleteUser("long", "123")
	if err != nil {
		t.Errorf("testDeleteUser err:%v", err)
	}
}

func testReGetUser(t *testing.T) {
	pwd, err := GetUserCredential("long")
	if err == nil {
		t.Errorf("testGetUser err:%v", err)
	}

	if pwd != "" {
		t.Errorf("testGetUser pwd:%v", pwd)
	}
}

func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("AddVideo", testAddVideoInfo)
}

var tempvid string

func testAddVideoInfo(t *testing.T) {
	vi, err := AddNewVideo(1, "my-video")
	if err != nil {
		t.Errorf("testAddVideoInfo err:%v", err)
	}
	tempvid = vi.Id
}
