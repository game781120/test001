package utils

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
type CommonResponse struct {
	BaseResponse
	CurrentPage int `json:"currentPage"`
	Total       int `json:"total"`
}
