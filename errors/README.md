# errors

一个支持业务错误码的错误包，适用于 Web API 错误处理。

## 文档

[如何规范 RESTful API 的业务错误处理](https://jianghushinian.cn/2023/03/04/how-to-standardize-the-handling-of-restful-api-business-errors/)

## 错误码

业务错误码由 8 位纯数字组成，类型为 `int`。

业务错误码格式：`40000000`

错误码说明：

- 1-3 位: HTTP 状态码
- 4-5 位: 组件
- 6-8 位: 组件内部错误码

示例：

| 组件  | 组件编号 | HTTP 状态码 | 错误码      | 说明    |
|-----|------|----------|----------|-------|
| 用户  | 01   | 400      | 40001001 | 请求不合法 |
| 用户  | 01   | 401      | 40101001 | 认证失败  |
| 用户  | 01   | 403      | 40301001 | 授权失败  |
| 用户  | 01   | 404      | 40401001 | 资源未找到 |
| 用户  | 01   | 500      | 50001001 | 系统错误  |

## 错误包

错误包支持功能：

- `Wrap/Unwrap`
- `MarshalJSON/UnmarshalJSON`
- `Format`: `%s`、`%v`、`%+v`、`%#v`、`%q`
- stack

## 使用示例

### 创建错误码

```go
var (
	CodeBadRequest   = NewAPICode(40000000, "请求不合法")
	CodeUnauthorized = NewAPICode(40100000, "认证失败")
	CodeForbidden    = NewAPICode(40300000, "授权失败")
	CodeNotFound     = NewAPICode(40400000, "资源未找到")
	CodeUnknownError = NewAPICode(50000000, "系统错误")
)
```

### 创建错误对象

```go
var (
	err = errors.New("new error")
	apierr1 = NewAPIError(CodeNotFound)
	apierr2 = NewAPIError(CodeForbidden, err)
	apierr3 = WrapC(CodeUnknownError, err)
)
```

### 在 Gin 中使用

```go
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
```

### 错误返回结果

```json
{
	"code": 40400000,
	"message": "资源未找到",
	"reference": "https://jianghushinian.cn"
}
```

上述错误返回结果中 `code` 表示错误码，`message` 表示错误信息，`reference` 为可选的解决错误的文档地址。

更多使用详情请参考 [examples](./examples)。
