package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	apierr "github.com/jianghushinian/gokit/errors"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrDatabase        = errors.New("database error")
)

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func AccountOne(id int) (*Account, error) {
	for _, v := range accounts {
		if id == v.ID {
			return &v, nil
		}
	}

	// 模拟返回数据库错误
	if id == 500 {
		return nil, ErrDatabase
	}
	return nil, ErrAccountNotFound
}

var accounts = []Account{
	{ID: 1, Name: "account_1"},
	{ID: 2, Name: "account_2"},
	{ID: 3, Name: "account_3"},
}

func ShowAccount(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		// 将 errors 包装成 APIError 返回
		ResponseError(c, apierr.WrapC(apierr.CodeBadRequest, err))
		return
	}

	account, err := AccountOne(aid)
	if err != nil {
		switch {
		case errors.Is(err, ErrAccountNotFound):
			err = apierr.NewAPIError(apierr.CodeNotFound, err)
		case errors.Is(err, ErrDatabase):
			err = apierr.NewAPIError(apierr.CodeUnknownError, fmt.Errorf("account %d: %w", aid, err))
		}
		ResponseError(c, err)
		return
	}
	ResponseOK(c, account)
}

func main() {
	r := gin.Default()

	r.GET("/accounts/:id", ShowAccount)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
