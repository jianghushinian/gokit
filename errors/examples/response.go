package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	apierr "github.com/jianghushinian/gokit/errors"
)

func ResponseOK(c *gin.Context, spec interface{}) {
	if spec == nil {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, spec)
}

func ResponseError(c *gin.Context, err error) {
	log(err)
	e := apierr.ParseCoder(err)
	httpStatus := e.HTTPStatus()
	if httpStatus >= 500 {
		// send error msg to email/feishu/sentry...
		go fakeSendErrorEmail(err)
	}
	c.AbortWithStatusJSON(httpStatus, err)
}

// log 打印错误日志，输出堆栈
func log(err error) {
	fmt.Println("========== log start ==========")
	fmt.Printf("%+v\n", err)
	fmt.Println("========== log end ==========")
}

// fakeSendErrorEmail 模拟将错误信息发送到邮件，JSON 格式
func fakeSendErrorEmail(err error) {
	fmt.Println("========== error start ==========")
	fmt.Printf("%#v\n", err)
	fmt.Println("========== error end ==========")
}
