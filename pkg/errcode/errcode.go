package errcode

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Details []string `json:"details"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{Code: code, Msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码：%d, 错误信息：%s", e.ErrCode(), e.ErrMsg())
}

func (e *Error) ErrCode() int {
	return e.Code
}

func (e *Error) ErrMsg() string {
	return e.Msg
}

func (e *Error) ErrMsgf(args []interface{}) string {
	return fmt.Sprintf(e.Msg, args...)
}

func (e *Error) ErrDetails() []string {
	return e.Details
}

func (e *Error) ErrWithDetails(details ...string) *Error {
	newError := *e
	newError.Details = []string{}
	for _, d := range details {
		newError.Details = append(newError.Details, d)
	}

	return &newError
}

func (e *Error) StatusCode() int {
	switch e.ErrCode() {
	case Success.ErrCode():
		// 200
		return http.StatusOK
	case ServerError.ErrCode():
		// 500
		return http.StatusInternalServerError
	case InvalidParams.ErrCode():
		// 400
		return http.StatusBadRequest
	case NotFound.ErrCode():
		// 404
		return http.StatusNotFound
	case UnauthorizedAuthNotExist.ErrCode():
		// 使用fallthrough强制执行后面的case代码。
		// fallthrough不能用在switch的最后一个分支。
		fallthrough
	case UnauthorizedTokenError.ErrCode():
		fallthrough
	case UnauthorizedTokenGenerate.ErrCode():
		fallthrough
	case UnauthorizedTokenTimeout.ErrCode():
		// 401
		return http.StatusUnauthorized
	case TooManyRequests.ErrCode():
		// 429
		return http.StatusTooManyRequests
	}

	return http.StatusInternalServerError
}
