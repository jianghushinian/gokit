# config

一个简单的支持加载 YAML、JSON 两种格式文件的配置包。

## 文档

[如何用 Go 实现一个配置包](https://jianghushinian.cn/2023/03/31/how-to-implement-a-config-package-with-go/)

## 特性

- [x] 支持 YAML/JSON 两种配置文件格式
- [x] 支持序列化/反序列化
- [x] 支持通过环境变量/命令行参数指定配置文件

## 使用示例

以 YAML 格式配置为例，JSON 格式同理。

### 加载配置（反序列化）

有如下 `config.yaml` 配置：

```yaml
username: user
password: pass
server:
  endpoint: https://jianghushinian.cn/
```

示例代码：

```go
package main

import (
	"fmt"

	"github.com/jianghushinian/gokit/config/config"
)

type Config struct {
	Username string
	Password string
	Server   struct {
		Endpoint string
	}
}

func main() {
	c := &Config{}
	err := config.LoadOrDumpYAMLConfigFromFlag(c)
	fmt.Println(c, err)
}
```

通过环境变量指定配置文件：

```bash
$ export CONFIG_PATH=./config.yaml
$ go run main.go
```

或者通过命令行参数指定配置文件：

```bash
$ go run main.go -c ./config.yaml
```

控制台输出:

```bash
&{user pass {https://jianghushinian.cn/}} <nil>
```

**注意：** 命令行参数优先级高于环境变量，如果同时使用了环境变量和命令行参数，则以命令行参数为准。

### 将配置写入文件（序列化）

示例代码：

```go
package main

import (
	"fmt"

	"github.com/jianghushinian/gokit/config/config"
)

type Config struct {
	Username string
	Password string
	Server   struct {
		Endpoint string
	}
}

func main() {
	c := &Config{
		Username: "username",
		Password: "password",
		Server: struct {
			Endpoint string
		}{"https://jianghushinian.cn/"},
	}
	err := config.LoadOrDumpYAMLConfigFromFlag(c)
	fmt.Println(c, err)
}
```

通过环境变量指定参数：

```bash
$ export CONFIG_PATH=./config.yaml DUMP_CONFIG=true
$ go run main.go
```

或者通过命令行指定参数：

```bash
$ go run main.go -c ./config.yaml -d=true
```

得到 `config.yaml` 内容如下:

```yaml
username: username
password: password
server:
    endpoint: https://jianghushinian.cn/
```

**注意：** 程序在将配置写入文件后会自动执行 `os.Exit(0)` 退出，不会有任何输出。
