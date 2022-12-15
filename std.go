// Package logger
// Date: 2022/12/15 23:18:51
// Author: Amu
// Description:
package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var std *Logger

func init() {
	once.Do(func() {
		std = &Logger{
			Logger: zap.New(
				zapcore.NewCore(
					zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
						TimeKey:          "time",
						LevelKey:         "level",
						NameKey:          "logger",
						CallerKey:        "caller",
						MessageKey:       "message",
						StacktraceKey:    "stacktrace",
						LineEnding:       zapcore.DefaultLineEnding,
						EncodeLevel:      encodeLevel,
						EncodeTime:       encodeTime,
						EncodeDuration:   zapcore.SecondsDurationEncoder,
						EncodeCaller:     encodeCaller,
						ConsoleSeparator: " || ",
					}),
					zapcore.AddSync(os.Stdout),
					InfoLevel,
				),
				zap.AddCaller(),
				zap.AddCallerSkip(2),
			),
			name:    "std",
			loggers: make(map[string]*Logger),
		}
	})
}

func InitLogger(options ...Option) {
	config := &Config{
		name:                "std",
		logFile:             "default.log",
		logLevel:            InfoLevel,
		logFormat:           "text",
		logFileRotationTime: time.Hour * 24,
		logFileMaxAge:       time.Hour * 24 * 7,
		logOutput:           "stdout",
		logFileSuffix:       ".%Y%m%d",
		logSeparator:        " || ",
	}
	for _, option := range options {
		option(config)
	}
	fmt.Printf("config: %#v\n", config)

	encoder := getEncoder(config)
	writer := getWriter(config)
	level := config.logLevel

	std = &Logger{
		Logger: zap.New(
			zapcore.NewCore(encoder, writer, level),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		),
		name:    config.name,
		loggers: make(map[string]*Logger),
	}
}

func CreateLogger(options ...Option) {
	std.CreateLogger(options...)
}

func GetLoggerByName(name string) *Logger {
	if _, ok := std.loggers[name]; ok {
		return std.loggers[name]
	}
	return nil
}

func NewField(key string, value any) *Logger {
	return std.NewField(key, value)
}

func Debug(args ...any) {
	std.Debug(args...)
}

func Debugf(args ...any) {
	std.Debug(args...)
}

func Info(args ...any) {
	std.Info(args...)
}

func Infof(args ...any) {
	std.Infof(args...)
}

func Warn(args ...any) {
	std.Warn(args...)
}

func Warnf(args ...any) {
	std.Warnf(args...)
}

func Error(args ...any) {
	std.Error(args...)
}

func Errorf(args ...any) {
	std.Errorf(args...)
}

func Fatal(args ...any) {
	std.Fatal(args...)
}

func Fatalf(args ...any) {
	std.Fatalf(args...)
}

func Panic(args ...any) {
	std.Panic(args...)
}

func Panicf(args ...any) {
	std.Panicf(args...)
}
