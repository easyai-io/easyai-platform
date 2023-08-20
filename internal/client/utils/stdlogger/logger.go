package stdlogger

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/easyai-io/easyai-platform/pkg/contextx"
)

// SetDebugLevel used by root, Debug Level
func SetDebugLevel() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&customFormatter{})
}

// SetInfoLevel used by root, InfoLevel
func SetInfoLevel() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&customFormatter{})
}

// DebugSensitive only for maintainer/developer of easyai platform
func DebugSensitive(ctx context.Context, v string, args ...interface{}) {
	if !_isAdmin(ctx) {
		return
	}
	log.Debugf(v, args...)
}

var (
	_adminUserSet = map[string]bool{
		"shuaiyy": true,
		"admin":   true,
	}

	_isAdmin = func(ctx context.Context) bool {
		return _adminUserSet[contextx.FromUserUID(ctx)]
	}
)

// Red 函数将输入字符串渲染为红色
func Red(s ...interface{}) string {
	return color.New(color.FgRed).Sprint(s...)
}

// Blue 函数将输入字符串渲染为蓝色
func Blue(s ...interface{}) string {
	return color.New(color.FgBlue).Sprint(s...)
}

// Green 函数将输入字符串渲染为绿色
func Green(s ...interface{}) string {
	return color.New(color.FgGreen).Sprint(s...)
}

// Yellow 函数将输入字符串渲染为黄色
func Yellow(s ...interface{}) string {
	return color.New(color.FgYellow).Sprint(s...)
}

type customFormatter struct {
}

// Format 实现 logrus.Formatter 接口
func (f *customFormatter) Format(entry *log.Entry) ([]byte, error) {
	// 设置时间格式
	timeFormat := "2006/01/02 15:04:05"
	var result string

	// 拼接日志字符串
	levelString := fmt.Sprintf("[%s]", entry.Level)

	switch entry.Level {
	case log.DebugLevel:
		result = fmt.Sprintf("%s %s %s\n", entry.Time.Format(timeFormat), Yellow(levelString), entry.Message)

	case log.WarnLevel:
		result = fmt.Sprintf("%s %s %s\n", entry.Time.Format(timeFormat), Yellow(levelString), Yellow(entry.Message))

	case log.ErrorLevel:
		result = fmt.Sprintf("%s %s %s\n", entry.Time.Format(timeFormat), Red(levelString), Red(entry.Message))

	case log.PanicLevel:
		result = fmt.Sprintf("%s %s %s\n", entry.Time.Format(timeFormat), Red(levelString), entry.Message)

	case log.InfoLevel:
		result = fmt.Sprintf("%s %s %s\n", entry.Time.Format(timeFormat), Blue(levelString), entry.Message)

	default:
		result = fmt.Sprintf("%s %s %s\n", entry.Time.Format(timeFormat), Blue(levelString), entry.Message)
	}

	return []byte(result), nil
}

// 全都调用logrus中的实现

// Info level ,nomal level
func Info(args ...interface{}) {
	log.Info(args...)
}

// Debug Level
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Warn logs a message with severity WARN.
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error logs a message with severity ERROR.
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal logs a message with severity ERROR followed by a call to os.Exit().
func Fatal(args ...interface{}) {
	log.Panic(args...)
}

// Infof level ,nomal level
func Infof(v string, args ...interface{}) {
	log.Infof(v, args...)
}

// Debugf Level
func Debugf(v string, args ...interface{}) {
	log.Debugf(v, args...)
}

// Warnf logs a message with severity WARN.
func Warnf(v string, args ...interface{}) {
	log.Warnf(v, args...)
}

// Errorf logs a message with severity ERROR.
func Errorf(v string, args ...interface{}) {
	log.Errorf(v, args...)
}

// Fatalf logs a message with severity ERROR followed by a call to os.Exit().
func Fatalf(v string, args ...interface{}) {
	log.Panicf(v, args...)
}
