package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

// getLbAddr 返回后端服务地址（对应原项目 config.GetLbAddr()）
func getLbAddr() string {
	return "127.0.0.1"
}

func request(b *ApiBody, c *gin.Context) {
	var resp *http.Response
	var err error

	// 解析前端传来的目标 URL，将 host 替换为后端服务地址，保留原端口
	u, _ := url.Parse(b.Url)
	u.Host = getLbAddr() + ":" + u.Port()
	newUrl := u.String()

	switch b.Method {
	case http.MethodGet:
		// GET 请求：无请求体，直接转发
		req, _ := http.NewRequest("GET", newUrl, nil)
		req.Header = c.Request.Header // 复制原始请求的所有 Header（含 X-Session-Id）
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(c, resp)

	case http.MethodPost:
		// POST 请求：携带 ReqBody 作为请求体
		req, _ := http.NewRequest("POST", newUrl, bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = c.Request.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(c, resp)

	case http.MethodDelete:
		// DELETE 请求：无请求体，直接转发
		req, _ := http.NewRequest("DELETE", newUrl, nil)
		req.Header = c.Request.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		normalResponse(c, resp)

	default:
		// 不支持的 HTTP 方法
		c.String(http.StatusBadRequest, "Bad api request")
		return
	}
}

func normalResponse(c *gin.Context, r *http.Response) {
	// 读取后端服务返回的响应体
	res, err := io.ReadAll(r.Body)
	if err != nil {
		// 读取失败，返回 500 内部错误
		re, _ := json.Marshal(ErrorInternalFaults)
		c.String(http.StatusInternalServerError, string(re))
		return
	}

	// 将后端的状态码和响应体原样返回给前端
	c.String(r.StatusCode, string(res))
}