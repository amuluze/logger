// Package logger
// Date: 2022/12/15 23:19:05
// Author: Amu
// Description:
package logger

import (
	"errors"
	"testing"
)

func TestLogger(t *testing.T) {
	std.Info("hello status: ", 200)

	err := errors.New("request failed")
	std.Errorf("this is a error message: %s", err)

	std.NewField("hello", false)
	std.Info("hello bool")
	std.Errorf("hello, this is test with new field")
}

func TestInitLogger(t *testing.T) {
	InitLogger(
		SetLogLevel("error"),
		SetLogFormat("text"),
	)

	std.Info("hello", "world", errors.New("this is a error message"))
	err := errors.New("test error log")
	std.Errorf("hello %v", err)

	std.Info("good")
}

func TestCreateLogger(t *testing.T) {
	CreateLogger(
		SetName("logger"),
		SetLogLevel("info"),
		SetLogFormat("text"),
	)

	log := GetLoggerByName("logger")
	log.Info("hello logger")

	log.Error("test create")
}
