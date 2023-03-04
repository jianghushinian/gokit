package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func TestNewAPIError(t *testing.T) {
	type args struct {
		coder APICoder
		cause []error
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "test_coder",
			args: args{
				coder: CodeUnauthorized,
			},
			want: &apiError{
				coder: CodeUnauthorized,
			},
		},
		{
			name: "test_coder_and_cause",
			args: args{
				coder: CodeBadRequest,
				cause: []error{errors.New("cause")},
			},
			want: &apiError{
				coder: CodeBadRequest,
				cause: errors.New("cause"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NewAPIError(tt.args.coder, tt.args.cause...); err.Error() != tt.want.Error() {
				t.Errorf("NewAPIError() = %v, wantErr %v", err.Error(), tt.want.Error())
			}
			if err := WrapC(tt.args.coder, tt.args.cause...); err.Error() != tt.want.Error() {
				t.Errorf("WrapC() = %v, wantErr %v", err.Error(), tt.want.Error())
			}
		})
	}
}

func Test_apiError_Coder(t *testing.T) {
	type fields struct {
		coder APICoder
		cause error
		stack *stack
	}
	tests := []struct {
		name   string
		fields fields
		want   APICoder
	}{
		{
			name: "normal",
			fields: fields{
				coder: CodeForbidden,
				cause: NewAPIError(CodeBadRequest),
			},
			want: CodeForbidden,
		},
		{
			name: "NewAPICode",
			fields: fields{
				coder: NewAPICode(40100000, "Unauthorized"),
				cause: NewAPIError(CodeBadRequest),
			},
			want: NewAPICode(40100000, "Unauthorized"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiError{
				coder: tt.fields.coder,
				cause: tt.fields.cause,
				stack: tt.fields.stack,
			}
			if got := a.Coder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Coder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiError_Error(t *testing.T) {
	type fields struct {
		coder APICoder
		cause error
		stack *stack
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				coder: CodeForbidden,
			},
			want: "[40300000] - 授权失败",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiError{
				coder: tt.fields.coder,
				cause: tt.fields.cause,
				stack: tt.fields.stack,
			}
			if got := a.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiError_Unwrap(t *testing.T) {
	type fields struct {
		coder APICoder
		cause error
		stack *stack
	}
	tests := []struct {
		name   string
		fields fields
		want   error
	}{
		{
			name: "errors.New",
			fields: fields{
				coder: CodeUnknownError,
				cause: errors.New("cause"),
			},
			want: errors.New("cause"),
		},
		{
			name: "NewAPIError",
			fields: fields{
				coder: CodeUnknownError,
				cause: NewAPIError(CodeUnknownError),
			},
			want: NewAPIError(CodeUnknownError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiError{
				coder: tt.fields.coder,
				cause: tt.fields.cause,
				stack: tt.fields.stack,
			}
			if err := a.Unwrap(); err.Error() != tt.want.Error() {
				t.Errorf("Unwrap() error = %v, want %v", err.Error(), tt.want.Error())
			}
		})
	}
}

func Test_apiError_Format(t *testing.T) {
	tests := []struct {
		name string
		err  error
		verb string
		want string
	}{
		{
			name: "test_%v",
			err:  NewAPIError(CodeBadRequest),
			verb: "v",
			want: `^\[40000000\] - 请求不合法$`,
		},
		{
			name: "test_%+v",
			err:  NewAPIError(CodeBadRequest),
			verb: "+v",
			want: `^\[40000000\] - 请求不合法` + `[\s\S]*gokit/errors\.Test_apiError_Format`,
		},
		{
			name: "test_%+v_and_cause",
			err:  NewAPIError(CodeBadRequest, errors.New("cause")),
			verb: "+v",
			want: `^\[40000000\] - 请求不合法` + `[\s\S]*cause` + `[\s\S]*gokit/errors\.Test_apiError_Format`,
		},
		{
			name: "test_%#v_and_cause",
			err:  NewAPIError(CodeBadRequest, errors.New("cause")),
			verb: "#v",
			want: `^{"code":40000000,"message":"请求不合法","cause":"cause","stack":"[\S]+gokit/errors.Test_apiError_Format.*}$`,
		},
		{
			name: "test_%q",
			err:  NewAPIError(CodeBadRequest),
			verb: "q",
			want: `^"\[40000000\] - 请求不合法"$`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf("%"+tt.verb, tt.err); !regexp.MustCompile(tt.want).Match([]byte(got)) {
				t.Errorf("got = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_apiError_Format_by_fakeFmtState(t *testing.T) {
	type fields struct {
		coder APICoder
		cause error
		stack *stack
	}
	type args struct {
		s    fmt.State
		verb rune
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		{
			name: "test_%s",
			fields: fields{
				coder: CodeNotFound,
				cause: errors.New("cause"),
			},
			args: args{
				s:    fakeFmtState{},
				verb: 's',
			},
			want: []byte("[40400000] - 资源未找到"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiError{
				coder: tt.fields.coder,
				cause: tt.fields.cause,
				stack: tt.fields.stack,
			}
			a.Format(tt.args.s, tt.args.verb)
			if string(buf) != string(tt.want) {
				t.Errorf("got = %s, want %s", buf, tt.want)
			}
		})
	}
}

func Test_apiError_MarshalJSON(t *testing.T) {
	type fields struct {
		coder APICoder
		cause error
		stack *stack
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				coder: CodeBadRequest,
				cause: errors.New("cause"),
			},
			want:    []byte(`{"code":40000000,"message":"请求不合法"}`),
			wantErr: false,
		},
		{
			name: "with_reference",
			fields: fields{
				coder: NewAPICode(
					CodeBadRequest.Code(),
					CodeBadRequest.Message(),
					"https://github.com/jianghushinian/gokit/blob/main/README.md",
				),
			},
			want:    []byte(`{"code":40000000,"message":"请求不合法","reference":"https://github.com/jianghushinian/gokit/blob/main/README.md"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiError{
				coder: tt.fields.coder,
				cause: tt.fields.cause,
				stack: tt.fields.stack,
			}
			got, err := a.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_apiError_UnmarshalJSON(t *testing.T) {
	type fields struct {
		coder APICoder
		cause error
		stack *stack
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			fields: fields{
				coder: CodeBadRequest,
				cause: errors.New("cause"),
			},
			args: args{
				data: []byte(`{"code":40000000,"message":"请求不合法"}`),
			},
			wantErr: false,
		},
		{
			name: "with_reference",
			fields: fields{
				coder: NewAPICode(
					CodeBadRequest.Code(),
					CodeBadRequest.Message(),
					"https://github.com/jianghushinian/gokit/blob/main/README.md",
				),
				cause: errors.New("cause"),
			},
			args: args{
				data: []byte(`{"code":40000000,"message":"请求不合法","reference":"https://github.com/jianghushinian/gokit/blob/main/README.md"}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &apiError{
				coder: tt.fields.coder,
				cause: tt.fields.cause,
				stack: tt.fields.stack,
			}
			if err := a.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if v, err := json.Marshal(a); err != nil || string(v) != string(tt.args.data) {
				t.Errorf("UnmarshalJSON() got = %s, want %s", v, tt.args.data)
			}
		})
	}
}

// ref: https://medium.com/@virup/testing-in-go-making-use-of-duck-typing-6927feb125c6
type fakeFmtState struct{}

var buf []byte

func (fakeFmtState) Write(b []byte) (n int, err error) {
	buf = b
	return len(string(b)), nil
}

func (fakeFmtState) Width() (wid int, ok bool) {
	return -1, false
}

func (fakeFmtState) Precision() (prec int, ok bool) {
	return -1, false
}

func (fakeFmtState) Flag(c int) bool {
	return c == '+'
}
