package main

import (
	"goconf.googlecode.com/hg"
	"fmt"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota;
	INFO;
	WARN;
	ERROR
)

type LogMsg struct {
	Level LogLevel
	Msg string
	Timestamp *time.Time
}

func NewLogMsg(level LogLevel, msg string) LogMsg {
	return LogMsg{ level, msg, time.UTC() }
}

func Logger(c* conf.ConfigFile, logchan chan LogMsg) {
	// TODO: read configuration about logging threshold, log target, etc.
	for {
		lm := <-logchan
		severity := LogLevelToString(lm.Level)
		timestamp := lm.Timestamp.String()
		fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", timestamp, severity, lm.Msg)
	}
}

func LogLevelToString(level LogLevel) string {
	var loglevel_names = map[LogLevel] string {
		DEBUG: "DEBUG",
		INFO: "INFO",
		WARN: "WARN",
		ERROR: "ERROR",
	}
	return loglevel_names[level]
}

