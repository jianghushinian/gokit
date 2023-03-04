package errors

import (
	"encoding/json"
	"fmt"
	"io"
)

var WrapC = NewAPIError

func NewAPIError(coder APICoder, cause ...error) error {
	var c error
	if len(cause) > 0 {
		c = cause[0]
	}
	return &apiError{
		coder: coder,
		cause: c,
		stack: callers(),
	}
}

type apiError struct {
	coder APICoder
	cause error
	*stack
}

// Error implement interface error
func (a *apiError) Error() string {
	return fmt.Sprintf("[%d] - %s", a.coder.Code(), a.coder.Message())
}

func (a *apiError) Coder() APICoder {
	return a.coder
}

func (a *apiError) Unwrap() error {
	return a.cause
}

func (a *apiError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			str := a.Error()
			if a.Unwrap() != nil {
				str += " " + a.Unwrap().Error()
			}
			_, _ = io.WriteString(s, str)
			a.stack.Format(s, verb)
			return
		}
		if s.Flag('#') {
			cause := ""
			if a.Unwrap() != nil {
				cause = a.Unwrap().Error()
			}
			data, _ := json.Marshal(errorMessage{
				Code:      a.coder.Code(),
				Message:   a.coder.Message(),
				Reference: a.coder.Reference(),
				Cause:     cause,
				Stack:     fmt.Sprintf("%+v", a.stack),
			})
			_, _ = io.WriteString(s, string(data))
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, a.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", a.Error())
	}
}

func (a *apiError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&errorMessage{
		Code:      a.coder.Code(),
		Message:   a.coder.Message(),
		Reference: a.coder.Reference(),
	})
}

func (a *apiError) UnmarshalJSON(data []byte) error {
	e := &errorMessage{}
	if err := json.Unmarshal(data, e); err != nil {
		return err
	}
	a.coder = NewAPICode(e.Code, e.Message, e.Reference)
	return nil
}

type errorMessage struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Reference string `json:"reference,omitempty"`
	Cause     string `json:"cause,omitempty"`
	Stack     string `json:"stack,omitempty"`
}
