package mlog

import "strings"

type Level uint8

func NewLevel(level string) (l Level) {
	level = strings.ToLower(level)

	return levelString[level]
}

func (l Level) String() string {
	return levelPrefix[l]
}

const (
	Debug Level = 1 << iota
	Info
	Notice
	Warning
	Error
	Critical
	Alert
	Emergency
)

const LevelAll = ^Level(0)

const (
	TFile uint8 = iota
	TBare
	TMultiHandler
)

var (
	levelPrefix = map[Level]string{
		Debug:     "[DEBUG]",
		Info:      "[INFO]",
		Notice:    "[NOTICE]",
		Warning:   "[WARNING]",
		Error:     "[ERROR]",
		Critical:  "[CRITICAL]",
		Alert:     "[ALERT]",
		Emergency: "[EMERGENCY]",
		LevelAll:  "[LA]",
	}

	levelString = map[string]Level{
		"debug":     Debug,
		"debug+":    LevelAll,
		"info":      Info,
		"info+":     ^Debug,
		"notice":    Notice,
		"notice+":   ^Debug & ^Info,
		"warning":   Warning,
		"warning+":  ^Debug & ^Info & ^Notice,
		"error":     Error,
		"error+":    Error | Critical | Alert | Emergency,
		"critical":  Critical,
		"alert":     Alert,
		"emergency": Emergency,
	}
)
