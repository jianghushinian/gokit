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
	// 开箱即用
	{
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

	// 选项
	{
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

	// 不同级别日志输出到不同位置
	{
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

	// 日志轮转
	{
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
}

type Hook struct{}

func (h Hook) OnWrite(ce *zapcore.CheckedEntry, field []zapcore.Field) {
	fmt.Printf("Fatal Hook: msg=%s, field=%+v\n", ce.Message, field)
}
