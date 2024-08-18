package common

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type termColors string

const (
	// Reset
	reset termColors = "\033[0m"
	// Regular Colors
	black  termColors = "\033[0;30m"
	red    termColors = "\033[0;31m"
	yellow termColors = "\033[0;33m"
	blue   termColors = "\033[0;34m"
	green  termColors = "\033[0;32m"
)

type Logger struct {
	LogMutex sync.Mutex
	Target   io.Writer
	Module   atomic.Value
	C        *gin.Context
}

//GENERAL SCHEMA:
// {color}[TIME][MODULE] MESSAGE{reset}

func NewLogger(target io.Writer, module string, c *gin.Context) *Logger {
	logger := &Logger{
		Target: target,
		C:      c,
	}
	logger.Module.Store(module)
	return logger
}

func (l *Logger) SetPrefix(prefix string) {
	l.Module.Store(prefix)
}

func getLocalTime() string {
	return time.Now().Format("2006-01-02 - 15:04:05")
}

func getIP(c *gin.Context) string {
	var clientIP string
	if c != nil {
		clientIP = c.ClientIP()
	} else {
		clientIP = "System"
	}
	return clientIP
}

func (l *Logger) Println(message interface{}) {
	l.LogMutex.Lock()
	defer l.LogMutex.Unlock()

	module_name := l.Module.Load().(string)
	_, _ = l.Target.Write([]byte(string(green) + "[" + getLocalTime() + "][" + getIP(l.C) + "][" + module_name + "] " + message.(string) + string(reset) + "\n"))
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.LogMutex.Lock()
	defer l.LogMutex.Unlock()

	module_name := l.Module.Load().(string)
	formattedMessage := fmt.Sprintf(format, args...)
	_, _ = l.Target.Write([]byte(string(green) + "[" + getLocalTime() + "][" + getIP(l.C) + "][" + module_name + "] " + formattedMessage + string(reset) + "\n"))
}

func (l *Logger) PrintfError(format string, args ...interface{}) {
	l.LogMutex.Lock()
	defer l.LogMutex.Unlock()

	module_name := l.Module.Load().(string)
	formattedMessage := fmt.Sprintf(format, args...)
	_, _ = l.Target.Write([]byte(string(red) + "[" + getLocalTime() + "][" + getIP(l.C) + "][" + module_name + "] " + formattedMessage + string(reset) + "\n"))
}

func (l *Logger) PrintfWarning(format string, args ...interface{}) {
	l.LogMutex.Lock()
	defer l.LogMutex.Unlock()

	module_name := l.Module.Load().(string)
	formattedMessage := fmt.Sprintf(format, args...)
	_, _ = l.Target.Write([]byte(string(yellow) + "[" + getLocalTime() + "][" + getIP(l.C) + "][" + module_name + "] " + formattedMessage + string(reset) + "\n"))
}

func (l *Logger) PrintfInfo(format string, args ...interface{}) {
	l.LogMutex.Lock()
	defer l.LogMutex.Unlock()

	module_name := l.Module.Load().(string)
	formattedMessage := fmt.Sprintf(format, args...)
	_, _ = l.Target.Write([]byte(string(blue) + "[" + getLocalTime() + "][" + getIP(l.C) + "][" + module_name + "] " + formattedMessage + string(reset) + "\n"))
}
