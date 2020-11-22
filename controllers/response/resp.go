package response

// ResponseBase ...
// 通用数据返回，只有要 body ，这两个字段一定有
type ResponseBase struct {
	// http code
	Code int `json:"RetCode"`
	// error msg
	Error string `json:"Message,omitempty"`
}
