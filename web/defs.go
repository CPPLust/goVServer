package main

// ApiBody 前端 AJAX 发来的统一网关请求格式
type ApiBody struct {
	Url     string `json:"url"`
	Method  string `json:"method"`
	ReqBody string `json:"req_body"`
}

// 错误响应
var (
	ErrorRequestNotRecognized   = map[string]string{"error": "request not recognized"}
	ErrorRequestBodyParseFailed = map[string]string{"error": "request body parse failed"}
	ErrorInternalFaults         = map[string]string{"error": "internal faults"}
)
