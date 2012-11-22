package main

import (
	"fmt"
	"github.com/akrennmair/goconf"
	"log"
	"log/syslog"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type LogMsg struct {
	Level     LogLevel
	Msg       string
	Timestamp time.Time
}

type SysLogger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

func NewSysLogger() *SysLogger {
	sl := &SysLogger{}
	sl.debug, _ = syslog.NewLogger(syslog.LOG_DEBUG, log.LstdFlags)
	sl.info, _ = syslog.NewLogger(syslog.LOG_INFO, log.LstdFlags)
	sl.warn, _ = syslog.NewLogger(syslog.LOG_WARNING, log.LstdFlags)
	sl.error, _ = syslog.NewLogger(syslog.LOG_ERR, log.LstdFlags)
	return sl
}

func (sl *SysLogger) Log(level LogLevel, msg string) {
	switch level {
	case DEBUG:
		sl.debug.Print(msg)
	case INFO:
		sl.info.Print(msg)
	case WARN:
		sl.warn.Print(msg)
	case ERROR:
		sl.error.Print(msg)
	}
}

func NewLogMsg(level LogLevel, msg string) LogMsg {
	return LogMsg{level, msg, time.Now().UTC()}
}

func Logger(c *conf.ConfigFile, logchan chan LogMsg) {
	use_syslog := false
	logmethod, err := c.GetString("", "log_method")
	if err == nil && logmethod == "syslog" {
		use_syslog = true
	}
	sl := NewSysLogger()
	for {
		lm := <-logchan
		timestamp := lm.Timestamp.String()
		if use_syslog {
			sl.Log(lm.Level, lm.Msg)
		} else {
			severity := LogLevelToString(lm.Level)
			fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", timestamp, severity, lm.Msg)
		}
	}
}

func LogLevelToString(level LogLevel) string {
	var loglevel_names = map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
	}
	return loglevel_names[level]
}
