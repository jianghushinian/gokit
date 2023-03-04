package errors

import "errors"

// 错误码说明
// 1-3 位: HTTP 状态码
// 4-5 位: 组件
// 6-8 位: 组件内部错误码
var (
	CodeBadRequest   = NewAPICode(40000000, "请求不合法")
	CodeUnauthorized = NewAPICode(40100000, "认证失败")
	CodeForbidden    = NewAPICode(40300000, "授权失败")
	CodeNotFound     = NewAPICode(40400000, "资源未找到")
	CodeUnknownError = NewAPICode(50000000, "系统错误", "https://github.com/jianghushinian/gokit/tree/main/errors")
)

type APICoder interface {
	Code() int
	Message() string
	Reference() string
	HTTPStatus() int
}

func NewAPICode(code int, message string, reference ...string) APICoder {
	ref := ""
	if len(reference) > 0 {
		ref = reference[0]
	}
	return &apiCode{
		code: code,
		msg:  message,
		ref:  ref,
	}
}

type apiCode struct {
	code int
	msg  string
	ref  string
}

func (a *apiCode) Code() int {
	return a.code
}

func (a *apiCode) Message() string {
	return a.msg
}

func (a *apiCode) Reference() string {
	return a.ref
}

func (a *apiCode) HTTPStatus() int {
	v := a.Code()
	for v >= 1000 {
		v /= 10
	}
	return v
}

func ParseCoder(err error) APICoder {
	for {
		if e, ok := err.(interface {
			Coder() APICoder
		}); ok {
			return e.Coder()
		}
		if errors.Unwrap(err) == nil {
			return CodeUnknownError
		}
		err = errors.Unwrap(err)
	}
}

func IsCode(err error, coder APICoder) bool {
	if err == nil {
		return false
	}

	for {
		if e, ok := err.(interface {
			Coder() APICoder
		}); ok {
			if e.Coder().Code() == coder.Code() {
				return true
			}
		}

		if errors.Unwrap(err) == nil {
			return false
		}
		err = errors.Unwrap(err)
	}
}
