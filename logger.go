// Package logger
// Date: 2022/12/15 23:17:48
// Author: Amu
// Description:
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	rotator "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once

type Logger struct {
	*zap.Logger
	name    string
	lock    sync.Mutex
	fields  []zap.Field
	loggers map[string]*Logger
}

func (l *Logger) CreateLogger(options ...Option) {
	l.lock.Lock()
	defer l.lock.Unlock()
	config := &Config{
		name:                "std",
		logFile:             "default.log",
		logLevel:            InfoLevel,
		logFormat:           "text",
		logFileRotationTime: time.Hour * 24,
		logFileMaxAge:       time.Hour * 24 * 7,
		logOutput:           "stdout",
		logFileSuffix:       ".%Y%m%d",
	}
	for _, option := range options {
		option(config)
	}

	if _, ok := l.loggers[config.name]; ok {
		return
	}
	encoder := getEncoder(config)
	writer := getWriter(config)
	level := config.logLevel

	newLogger := &Logger{
		Logger: zap.New(
			zapcore.NewCore(encoder, writer, level),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		),
		name:    config.name,
		loggers: make(map[string]*Logger),
	}
	l.loggers[config.name] = newLogger
}

func (l *Logger) NewField(key string, value any) *Logger {
	switch value := value.(type) {
	case int:
		l.fields = append(l.fields, zap.Int(key, value))
	case int32:
		l.fields = append(l.fields, zap.Int32(key, value))
	case int64:
		l.fields = append(l.fields, zap.Int64(key, value))
	case float32:
		l.fields = append(l.fields, zap.Float32(key, value))
	case float64:
		l.fields = append(l.fields, zap.Float64(key, value))
	case string:
		l.fields = append(l.fields, zap.String(key, value))
	case time.Time:
		l.fields = append(l.fields, zap.Time(key, value))
	case time.Duration:
		l.fields = append(l.fields, zap.Duration(key, value))
	case bool:
		l.fields = append(l.fields, zap.Bool(key, value))
	}
	return l
}

func (l *Logger) Debug(args ...any) {
	l.Logger.Debug(fmt.Sprint(args...), l.fields...)
}

func (l *Logger) Debugf(args ...any) {
	l.Logger.Debug(fmt.Sprintf(args[0].(string), args[1:]...), l.fields...)
}

func (l *Logger) Info(args ...any) {
	l.Logger.Info(fmt.Sprint(args...), l.fields...)
}

func (l *Logger) Infof(args ...any) {
	l.Logger.Info(fmt.Sprintf(args[0].(string), args[1:]...), l.fields...)
}

func (l *Logger) Warn(args ...any) {
	l.Logger.Warn(fmt.Sprint(args...), l.fields...)
}

func (l *Logger) Warnf(args ...any) {
	l.Logger.Warn(fmt.Sprintf(args[0].(string), args[1:]...), l.fields...)
}

func (l *Logger) Error(args ...any) {
	l.Logger.Error(fmt.Sprint(args...), l.fields...)
}

func (l *Logger) Errorf(args ...any) {
	l.Logger.Error(fmt.Sprintf(args[0].(string), args[1:]...), l.fields...)
}

func (l *Logger) Fatal(args ...any) {
	l.Logger.Fatal(fmt.Sprint(args...), l.fields...)
}

func (l *Logger) Fatalf(args ...any) {
	l.Logger.Fatal(fmt.Sprintf(args[0].(string), args[1:]...), l.fields...)
}

func (l *Logger) Panic(args ...any) {
	l.Logger.Panic(fmt.Sprint(args...), l.fields...)
}

func (l *Logger) Panicf(args ...any) {
	l.Logger.Panic(fmt.Sprintf(args[0].(string), args[1:]...), l.fields...)
}

func getEncoder(config *Config) zapcore.Encoder {
	if config.logFormat == "text" {
		return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
			TimeKey:          "timestamp",
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
		})
	}

	var baseConfig = zapcore.EncoderConfig{
		// 下面以 Key 结尾的参数表示，Json格式日志中的 key
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel, // 日志级别的以大写还是小写输出
		EncodeTime:     encodeTime,  // timestamp 时间字段的时间字符串格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   encodeCaller, // caller 路径
	}
	return zapcore.NewJSONEncoder(baseConfig)
}

func encodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(level.CapitalString())
}

func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(TimeFormat))
}

func encodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(caller.TrimmedPath())
}

func getWriter(config *Config) zapcore.WriteSyncer {
	if config.logOutput == "stdout" {
		return zapcore.AddSync(os.Stdout)
	}
	logFilePath := config.logFile
	if !filepath.IsAbs(config.logFile) {
		abspath, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), config.logFile))
		logFilePath = abspath
	}

	_log, _ := rotator.New(
		filepath.Join(logFilePath+config.logFileSuffix),
		// 生成软连接，指向最新的日志文件
		rotator.WithLinkName(logFilePath),
		// 保留文件期限
		rotator.WithMaxAge(config.logFileMaxAge),
		// 日志文件的切割间隔
		rotator.WithRotationTime(config.logFileRotationTime),
	)
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(_log), zapcore.AddSync(os.Stdout))
}
