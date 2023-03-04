package errors

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNewAPICode(t *testing.T) {
	type args struct {
		code      int
		message   string
		reference []string
	}
	tests := []struct {
		name string
		args args
		want APICoder
	}{
		{
			name: "normal",
			args: args{
				code:    40000001,
				message: "用户输入参数错误",
			},
			want: &apiCode{
				code: 40000001,
				msg:  "用户输入参数错误",
			},
		},
		{
			name: "test_exist_code",
			args: args{
				code:    40000000,
				message: "请求不合法",
			},
			want: CodeBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAPICode(tt.args.code, tt.args.message, tt.args.reference...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAPICode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiCode_Code(t *testing.T) {
	type fields struct {
		code int
		msg  string
		ref  string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "40100000",
			fields: fields{
				code: 40100000,
				msg:  "认证失败",
				ref:  "https://github.com/jianghushinian/gokit/blob/main/README.md",
			},
			want: 40100000,
		},
		{
			name: "50000000",
			fields: fields{
				code: CodeUnknownError.Code(),
				msg:  CodeUnknownError.Message(),
			},
			want: CodeUnknownError.Code(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiCode{
				code: tt.fields.code,
				msg:  tt.fields.msg,
				ref:  tt.fields.ref,
			}
			if got := a.Code(); got != tt.want {
				t.Errorf("Code() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiCode_Message(t *testing.T) {
	type fields struct {
		code int
		msg  string
		ref  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "40100000",
			fields: fields{
				code: 40100000,
				msg:  "认证失败",
				ref:  "https://github.com/jianghushinian/gokit/blob/main/README.md",
			},
			want: "认证失败",
		},
		{
			name: "50000000",
			fields: fields{
				code: CodeUnknownError.Code(),
				msg:  CodeUnknownError.Message(),
			},
			want: "系统错误",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiCode{
				code: tt.fields.code,
				msg:  tt.fields.msg,
				ref:  tt.fields.ref,
			}
			if got := a.Message(); got != tt.want {
				t.Errorf("Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiCode_Reference(t *testing.T) {
	type fields struct {
		code int
		msg  string
		ref  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "40100000",
			fields: fields{
				code: 40100000,
				msg:  "认证失败",
				ref:  "https://github.com/jianghushinian/gokit/blob/main/README.md",
			},
			want: "https://github.com/jianghushinian/gokit/blob/main/README.md",
		},
		{
			name: "50000000",
			fields: fields{
				code: CodeUnknownError.Code(),
				msg:  CodeUnknownError.Message(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiCode{
				code: tt.fields.code,
				msg:  tt.fields.msg,
				ref:  tt.fields.ref,
			}
			if got := a.Reference(); got != tt.want {
				t.Errorf("Reference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiCode_HTTPStatus(t *testing.T) {
	type fields struct {
		code int
		msg  string
		ref  string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "40100000",
			fields: fields{
				code: 40100000,
				msg:  "认证失败",
				ref:  "https://github.com/jianghushinian/gokit/blob/main/README.md",
			},
			want: 401,
		},
		{
			name: "50000000",
			fields: fields{
				code: CodeUnknownError.Code(),
				msg:  CodeUnknownError.Message(),
			},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiCode{
				code: tt.fields.code,
				msg:  tt.fields.msg,
				ref:  tt.fields.ref,
			}
			if got := a.HTTPStatus(); got != tt.want {
				t.Errorf("HTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCoder(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want APICoder
	}{
		{
			name: "forbidden_error",
			args: args{
				err: NewAPIError(CodeForbidden),
			},
			want: CodeForbidden,
		},
		{
			name: "unknown_error",
			args: args{
				err: errors.New("unknown"),
			},
			want: CodeUnknownError,
		},
		{
			name: "warp_unauthorized_error",
			args: args{
				err: fmt.Errorf("%w", NewAPIError(CodeUnauthorized)),
			},
			want: CodeUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCoder(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCode(t *testing.T) {
	type args struct {
		err   error
		coder APICoder
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal",
			args: args{
				err:   NewAPIError(CodeNotFound),
				coder: CodeNotFound,
			},
			want: true,
		},
		{
			name: "warp_code",
			args: args{
				err:   NewAPIError(CodeNotFound, NewAPIError(CodeBadRequest)),
				coder: CodeBadRequest,
			},
			want: true,
		},
		{
			name: "CodeNotFound_and_CodeUnknownError_no_match",
			args: args{
				err:   NewAPIError(CodeNotFound),
				coder: CodeUnknownError,
			},
			want: false,
		},
		{
			name: "errors.New_and_CodeNotFound_no_match",
			args: args{
				err:   errors.New("unknown"),
				coder: CodeNotFound,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCode(tt.args.err, tt.args.coder); got != tt.want {
				t.Errorf("IsCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
