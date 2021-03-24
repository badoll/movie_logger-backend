package api

// 响应码
const (
	Succ     = 0
	ParamErr = 10000
	DBErr    = 10001
	HTTPErr  = 20000
)

// NilStruct 返回空结构体
var NilStruct = struct{}{}

// Response 标准返回
type Response struct {
	Base struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"base"`
	Result interface{} `json:"result"`
}

// NewResp 构建标准返回结构体
func NewResp(status int, msg string, resp interface{}) Response {
	return Response{
		Base: struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		}{
			status,
			msg,
		},
		Result: resp,
	}
}
