package common

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type termColor string

const (
	termRed    termColor = "\033[31m"
	termGreen  termColor = "\033[32m"
	termGray   termColor = "\033[90m"
	termYellow termColor = "\033[33m"
	termReset  termColor = "\033[0m"
)

type Logger struct {
	module string
	level  LogLevel
}

func NewLogger(module string, level LogLevel) *Logger {
	return &Logger{
		module: module,
		level:  level,
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	now := time.Now().UTC().Format(time.RFC3339)
	if l.level <= DebugLevel {
		fmt.Printf("%s[%s][%s][DEBUG] %s%s\n", termGreen, now, l.module, fmt.Sprintf(msg, args...), termReset)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	now := time.Now().UTC().Format(time.RFC3339)
	if l.level <= InfoLevel {
		fmt.Printf("[%s][%s][INFO] %s\n", now, l.module, fmt.Sprintf(msg, args...))
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	now := time.Now().UTC().Format(time.RFC3339)
	if l.level <= WarnLevel {
		fmt.Printf("[%s][%s][WARN] %s\n", now, l.module, fmt.Sprintf(msg, args...))
	}
}
