package zap

import (
	"io"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RotateConfig struct {
	// 共用配置
	Filename string // 完整文件名
	MaxAge   int    // 保留旧日志文件的最大天数

	// 按时间轮转配置
	RotationTime time.Duration // 日志文件轮转时间

	// 按大小轮转配置
	MaxSize    int  // 日志文件最大大小（MB）
	MaxBackups int  // 保留日志文件的最大数量
	Compress   bool // 是否对日志文件进行压缩归档
	LocalTime  bool // 是否使用本地时间，默认 UTC 时间
}

// NewProductionRotateByTime 创建按时间轮转的 io.Writer
func NewProductionRotateByTime(filename string) io.Writer {
	return NewRotateByTime(NewProductionRotateConfig(filename))
}

// NewProductionRotateBySize 创建按大小轮转的 io.Writer
func NewProductionRotateBySize(filename string) io.Writer {
	return NewRotateBySize(NewProductionRotateConfig(filename))
}

func NewProductionRotateConfig(filename string) *RotateConfig {
	return &RotateConfig{
		Filename: filename,
		MaxAge:   30, // 日志保留 30 天

		RotationTime: time.Hour * 24, // 24 小时轮转一次

		MaxSize:    100, // 100M
		MaxBackups: 100,
		Compress:   true,
		LocalTime:  false,
	}
}

func NewRotateByTime(cfg *RotateConfig) io.Writer {
	opts := []rotatelogs.Option{
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge) * time.Hour * 24),
		rotatelogs.WithRotationTime(cfg.RotationTime),
		rotatelogs.WithLinkName(cfg.Filename),
	}
	if !cfg.LocalTime {
		rotatelogs.WithClock(rotatelogs.UTC)
	}
	filename := strings.SplitN(cfg.Filename, ".", 2)
	l, _ := rotatelogs.New(
		filename[0]+".%Y-%m-%d-%H-%M-%S."+filename[1],
		opts...,
	)
	return l
}

func NewRotateBySize(cfg *RotateConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  cfg.LocalTime,
		Compress:   cfg.Compress,
	}
}
