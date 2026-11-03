package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"goVServer/api/defs"

	"github.com/gin-gonic/gin"
)

// const VIDEO_DIR = "./videos/"
const VIDEO_DIR = "d:/gowork/code/3.code/goVServer/streamserver/videos/"
const MAX_UPLOAD_SIZE = 50 << 20 // 50MB

func streamHanler(c *gin.Context) {
	vid := c.Param("vid-id")
	vl := VIDEO_DIR + vid

	video, err := os.Open(vl)
	if err != nil {
		fmt.Printf("Error when try to open file: %v", err)
		sendErrorResponse(c, defs.ErrorInternalFaults)
		return
	}
	defer video.Close()

	// 获取文件信息（用于 Content-Length 和 Last-Modified）
	stat, err := video.Stat()
	if err != nil {
		log.Printf("Error when get file stat: %v", err)
		sendErrorResponse(c, defs.ErrorInternalFaults)
		return
	}

	// 根据文件扩展名获取 Content-Type
	contentType := mime.TypeByExtension(filepath.Ext(vl))
	if contentType == "" {
		// 扩展名识别不了，读取文件头 512 字节来检测
		buffer := make([]byte, 512)
		_, _ = video.Read(buffer)
		contentType = http.DetectContentType(buffer)
		video.Seek(0, 0) // 重置文件指针
	}

	c.Header("Content-Type", contentType)
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size()))

	// http.ServeContent 自动处理 Range / If-Modified-Since 等条件请求
	// 底层用 io.Copy 分段传输，不会一次性加载整个文件到内存
	http.ServeContent(c.Writer, c.Request, "", time.Now(), video)
}

func uploadHanler(c *gin.Context) {
	// 限制请求体大小
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MAX_UPLOAD_SIZE)
	if err := c.Request.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Printf("[uploadHanler] File is too big: %v", err)
		sendErrorResponse(c, defs.ErrorBadRequest)
		return
	}

	// 获取上传的文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("[uploadHanler] Error when try to get file: %v", err)
		sendErrorResponse(c, defs.ErrorInternalFaults)
		return
	}
	defer file.Close()

	// 读取文件内容
	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[uploadHanler] Read file error: %v", err)
		sendErrorResponse(c, defs.ErrorInternalFaults)
		return
	}

	// 获取视频 ID 并保存文件
	fn := c.Param("vid-id")
	err = os.WriteFile(VIDEO_DIR+fn, data, 0666)
	if err != nil {
		log.Printf("[uploadHanler] Write file error: %v", err)
		sendErrorResponse(c, defs.ErrorInternalFaults)
		return
	}

	fmt.Printf("[uploadHanler] Uploaded successfully: %s\n", fn)
	c.String(http.StatusCreated, "Uploaded successfully")
}

func testPageHandler(c *gin.Context) {
	t, err := template.ParseFiles(VIDEO_DIR + "upload.html")
	if err != nil {
		fmt.Printf("Parse template error: %v", err)
		sendErrorResponse(c, defs.ErrorInternalFaults)
		return
	}
	t.Execute(c.Writer, nil)
}
