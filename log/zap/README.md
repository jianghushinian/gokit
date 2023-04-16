# zap

基于 [zap](https://github.com/uber-go/zap) 开发的日志包，致力于提升 zap 使用体验。

## 文档

- [Go 第三方 log 库之 zap 使用](https://jianghushinian.cn/2023/03/19/use-of-zap-in-go-third-party-log-library/)

- [如何基于 zap 封装一个更好用的日志库](https://jianghushinian.cn/2023/04/16/how-to-wrap-a-more-user-friendly-logging-package-based-on-zap/)

## 特性

- [x] 类似 log 标准库的 API 设计
- [x] 动态修改日志级别
- [x] 可以设置不同日志级别输出到不同位置
- [x] 日志轮转，支持按时间/日志大小

## 使用示例

### 开箱即用

```go
package main

import (
	"os"
	"time"

	log "github.com/jianghushinian/gokit/log/zap"
)

func main() {
	defer log.Sync()
	log.Info("failed to fetch URL", log.String("url", "https://jianghushinian.cn/"))
	log.Warn("Warn msg", log.Int("attempt", 3))
	log.Error("Error msg", log.Duration("backoff", time.Second))

	// 修改日志级别
	log.SetLevel(log.ErrorLevel)
	log.Info("Info msg")
	log.Warn("Warn msg")
	log.Error("Error msg")

	// 替换默认 Logger
	file, _ := os.OpenFile("custom.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	logger := log.New(file, log.InfoLevel)
	log.ReplaceDefault(logger)
	log.Info("Info msg in replace default logger after")
}
```

控制台输出:

```log
{"level":"info","ts":"2023-03-19T21:57:59+08:00","msg":"failed to fetch URL","url":"https://jianghushinian.cn/"}
{"level":"warn","ts":"2023-03-19T21:57:59+08:00","msg":"Warn msg","attempt":3}
{"level":"error","ts":"2023-03-19T21:57:59+08:00","msg":"Error msg","backoff":1}
{"level":"error","ts":"2023-03-19T21:57:59+08:00","msg":"Error msg"}
```

`custom.log` 输出:

```log
{"level":"info","ts":"2023-03-19T21:57:59+08:00","msg":"Info msg in replace default logger after"}
```

### 选项

支持 [zap 选项](https://pkg.go.dev/go.uber.org/zap#Option)

```go
package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	log "github.com/jianghushinian/gokit/log/zap"
)

func main() {
	file, _ := os.OpenFile("test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	opts := []log.Option{
		// 附加日志调用信息
		log.WithCaller(true),
		log.AddCallerSkip(1),
		// Warn 级别日志 Hook
		log.Hooks(func(entry zapcore.Entry) error {
			if entry.Level == log.WarnLevel {
				fmt.Printf("Warn Hook: msg=%s\n", entry.Message)
			}
			return nil
		}),
		// Fatal 级别日志 Hook
		zap.WithFatalHook(Hook{}),
	}
	logger := log.New(io.MultiWriter(os.Stdout, file), log.InfoLevel, opts...)
	defer logger.Sync()

	logger.Info("Info msg", log.String("val", "string"))
	logger.Warn("Warn msg", log.Int("val", 7))
	logger.Fatal("Fatal msg", log.Time("val", time.Now()))
}

type Hook struct{}

func (h Hook) OnWrite(ce *zapcore.CheckedEntry, field []zapcore.Field) {
	fmt.Printf("Fatal Hook: msg=%s, field=%+v\n", ce.Message, field)
}
```

控制台输出:

```log
{"level":"info","ts":"2023-03-19T22:02:25+08:00","caller":"examples/main.go:55","msg":"Info msg","val":"string"}
{"level":"warn","ts":"2023-03-19T22:02:25+08:00","caller":"examples/main.go:56","msg":"Warn msg","val":7}
Warn Hook: msg=Warn msg
{"level":"fatal","ts":"2023-03-19T22:02:25+08:00","caller":"examples/main.go:57","msg":"Fatal msg","val":"2023-03-19T22:02:25+08:00"}
Fatal Hook: msg=Fatal msg, field=[{Key:val Type:16 Integer:1679234545108924000 String: Interface:Local}]
```

`test.log` 输出:

```log
{"level":"info","ts":"2023-03-19T22:02:25+08:00","caller":"examples/main.go:55","msg":"Info msg","val":"string"}
{"level":"warn","ts":"2023-03-19T22:02:25+08:00","caller":"examples/main.go:56","msg":"Warn msg","val":7}
{"level":"fatal","ts":"2023-03-19T22:02:25+08:00","caller":"examples/main.go:57","msg":"Fatal msg","val":"2023-03-19T22:02:25+08:00"}
```

### 不同级别日志输出到不同位置

`Info` 级别日志输出到 `os.Stdout`，`Warn` 级别日志输出到 `test-warn.log`，其他级别日志不会输出。

```go
package main

import (
	"os"

	log "github.com/jianghushinian/gokit/log/zap"
)

func main() {
	file, _ := os.OpenFile("test-warn.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	tees := []log.TeeOption{
		{
			Out: os.Stdout,
			LevelEnablerFunc: func(level log.Level) bool {
				return level == log.InfoLevel
			},
		},
		{
			Out: file,
			LevelEnablerFunc: func(level log.Level) bool {
				return level == log.WarnLevel
			},
		},
	}
	logger := log.NewTee(tees)
	defer logger.Sync()

	logger.Info("Info tee msg")
	logger.Warn("Warn tee msg")
	logger.Error("Error tee msg") // 不会输出
}
```

控制台输出:

```log
{"level":"info","ts":"2023-03-19T22:06:25+08:00","msg":"Info tee msg"}
```

`test-warn.log` 输出:

```log
{"level":"warn","ts":"2023-03-19T22:06:25+08:00","msg":"Warn tee msg"}
```

### 日志轮转

`Warn` 以下级别日志按大小轮转，`Warn` 及以上级别日志按时间轮转。

```go
package main

import (
	log "github.com/jianghushinian/gokit/log/zap"
)

func main() {
	tees := []log.TeeOption{
		{
			Out: log.NewProductionRotateBySize("rotate-by-size.log"),
			LevelEnablerFunc: log.LevelEnablerFunc(func(level log.Level) bool {
				return level < log.WarnLevel
			}),
		},
		{
			Out: log.NewProductionRotateByTime("rotate-by-time.log"),
			LevelEnablerFunc: log.LevelEnablerFunc(func(level log.Level) bool {
				return level >= log.WarnLevel
			}),
		},
	}
	lts := log.NewTee(tees)
	defer lts.Sync()

	lts.Debug("Debug msg")
	lts.Info("Info msg")
	lts.Warn("Warn msg")
	lts.Error("Error msg")
}
```

`rotate-by-size.log` 输出:

```log
{"level":"debug","ts":"2023-03-19T22:50:54+08:00","msg":"Debug msg"}
{"level":"info","ts":"2023-03-19T22:50:54+08:00","msg":"Info msg"}
```

`rotate-by-time.log` 输出:

```log
{"level":"warn","ts":"2023-03-19T22:50:54+08:00","msg":"Warn msg"}
{"level":"error","ts":"2023-03-19T22:50:54+08:00","msg":"Error msg"}
```

更多使用详情请参考 [examples](./examples)。
