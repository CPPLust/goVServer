package main

import (
	"encoding/json"
	"fmt"
	"goVServer/api/dbops"
	"goVServer/api/defs"
	"goVServer/api/session"
	"goVServer/api/utils"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	fmt.Println("CreateUser")
	res, _ := io.ReadAll(c.Request.Body)
	ubody := &defs.UserCredential{}

	fmt.Println(string(res))
	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(c, defs.ErrorRequestBodyParseFailed)
		return
	}
	fmt.Printf("recv name: %v, recv pwd:%v", ubody.Username, ubody.Pwd)

	// 先检查用户名是否已存在
	exists, err := dbops.UserExists(ubody.Username)
	if err != nil {
		sendErrorResponse(c, defs.ErrorDBError)
		return
	}
	if exists {
		fmt.Printf("%v +++++User already exists\n", ubody.Username)
		sendErrorResponse(c, defs.ErrorUserAlreadyExist)
		return
	}

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		// 兑底：并发场景下 UserExists 检查通过但 INSERT 时已被其他请求抢先插入
		if strings.Contains(err.Error(), "Duplicate entry") {
			sendErrorResponse(c, defs.ErrorUserAlreadyExist)
			return
		}
		sendErrorResponse(c, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUp{Success: true, Session: id}

	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(c, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(c, string(resp), http.StatusCreated)
	}
}

func Login(c *gin.Context) {
	// 读取请求体
	res, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("%s\n", res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		fmt.Printf("%s\n", err)
		sendErrorResponse(c, defs.ErrorRequestBodyParseFailed)
		return
	}

	// 校验 URL 参数用户名与请求体用户名是否一致
	uname := c.Param("username")
	fmt.Printf("Login url name: %s\n", uname)
	fmt.Printf("Login body name: %s\n", ubody.Username)
	if uname != ubody.Username {
		sendErrorResponse(c, defs.ErrorNotAuthUser)
		return
	}

	// 从数据库取密码并比对
	fmt.Printf("%s\n", ubody.Username)
	pwd, err := dbops.GetUserCredential(ubody.Username)
	fmt.Printf("Login pwd: %s\n", pwd)
	fmt.Printf("Login body pwd: %s\n", ubody.Pwd)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		sendErrorResponse(c, defs.ErrorNotAuthUser)
		return
	}

	// 密码正确，生成会话并返回
	id := session.GenerateNewSessionId(ubody.Username)
	si := &defs.SignedIn{Success: true, SessionId: id}
	if resp, err := json.Marshal(si); err != nil {
		sendErrorResponse(c, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(c, string(resp), http.StatusOK)
	}
}

func ListAllVideos(c *gin.Context) {
	if !ValidateUser(c.Request) {
		sendErrorResponse(c, defs.ErrorNotAuthUser)
		return
	}

	uname := c.Param("username")
	vs, err := dbops.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		fmt.Printf("Error in ListAllVideos: %s\n", err)
		sendErrorResponse(c, defs.ErrorDBError)
		return
	}

	vsi := &defs.VideosInfo{Videos: vs}
	// 确保 videos 不为 null，返回空数组
	if vsi.Videos == nil {
		fmt.Println("+++++++++++++++No videos found")
		vsi.Videos = []*defs.VideoInfo{}
	}
	if resp, err := json.Marshal(vsi); err != nil {
		sendErrorResponse(c, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(c, string(resp), http.StatusOK)

	}
}

func AddNewVideo(c *gin.Context) {
	if !ValidateUser(c.Request) {
		fmt.Println("[AddNewVideo] Unauthorized user")
		sendErrorResponse(c, defs.ErrorNotAuthUser)
		return
	}

	res, _ := io.ReadAll(c.Request.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal(res, nvbody); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(c, defs.ErrorRequestBodyParseFailed)
		return
	}

	vi, err := dbops.AddNewVideo(nvbody.AuthorId, nvbody.Name)
	log.Printf("[AddNewVideo ]Author id : %d, name: %s \n", nvbody.AuthorId, nvbody.Name)
	if err != nil {
		log.Printf("[AddNewVideo] Error in AddNewVideo: %s", err)
		sendErrorResponse(c, defs.ErrorDBError)
		return
	}

	if resp, err := json.Marshal(vi); err != nil {
		sendErrorResponse(c, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(c, string(resp), http.StatusOK)
	}
}

func GetUserInfo(c *gin.Context) {
	if !ValidateUser(c.Request) {
		fmt.Println("[GetUserInfo] Unauthorized user")
		sendErrorResponse(c, defs.ErrorNotAuthUser)
		return
	}

	uname := c.Param("username")
	fmt.Printf("[GetUserInfo] username=%s\n", uname)

	u, err := dbops.GetUser(uname)
	if err != nil {
		fmt.Printf("[GetUserInfo] Error: %v\n", err)
		sendErrorResponse(c, defs.ErrorDBError)
		return
	}

	ui := &defs.UserInfo{Id: u.Id}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(c, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(c, string(resp), http.StatusOK)
	}
}

func DeleteVideo(c *gin.Context) {
	fmt.Println("[DeleteVideo] delete video")
	if !ValidateUser(c.Request) {
		fmt.Println("[DeleteVideo] Unauthorized user")
		sendErrorResponse(c, defs.ErrorNotAuthUser)
		return
	}

	vid := c.Param("vid-id")
	fmt.Printf("[DeleteVideo] delete  vid-id %v\n", vid)
	err := dbops.DeleteVideoInfo(vid)
	if err != nil {
		log.Printf("[DeleteVideo] Error in DeletVideo: %s", err)
		sendErrorResponse(c, defs.ErrorDBError)
		return
	}

	//给删除的任务程序

	sendNormalResponse(c, "", http.StatusOK)
}

func PostComment(c *gin.Context) {
	fmt.Println("[PostComment] ====")
	if !ValidateUser(c.Request) {
		fmt.Println("[PostComment] Unauthorized user")
		sendErrorResponse(c, defs.ErrorNotAuthUser)
		return
	}

	res, _ := io.ReadAll(c.Request.Body)
	cbody := &defs.NewComment{}
	if err := json.Unmarshal(res, cbody); err != nil {
		log.Printf("[PostComment] %s", err)
		sendErrorResponse(c, defs.ErrorRequestBodyParseFailed)
		return
	}

	vid := c.Param("vid-id")
	fmt.Printf("[PostComment] show comment  vid-id %v\n", vid)
	if err := dbops.AddNewComments(vid, cbody.AuthorId, cbody.Content); err != nil {
		log.Printf("[PostComment] Error in PostComment: %s", err)
		sendErrorResponse(c, defs.ErrorDBError)
	} else {
		sendNormalResponse(c, "ok", http.StatusOK)
	}
}

func ShowComments(c *gin.Context) {
	fmt.Println("[ShowComments] ====")
	if !ValidateUser(c.Request) {
		fmt.Println("[ShowComments] Unauthorized user")
		sendErrorResponse(c, defs.ErrorNotAuthUser)
		return
	}
	vid := c.Param("vid-id")
	fmt.Printf("[ShowComments] show comment  vid-id %v\n", vid)
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ShowComments: %s", err)
		sendErrorResponse(c, defs.ErrorDBError)
		return
	}

	cms := &defs.Comments{Comments: cm}
	if cms.Comments == nil {
		fmt.Println("[ShowComments]+++++++++++++++No comment found")
		cms.Comments = []*defs.Comment{}
	}
	if resp, err := json.Marshal(cms); err != nil {
		sendErrorResponse(c, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(c, string(resp), http.StatusOK)
	}
}
