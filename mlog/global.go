package mlog

type Level uint8

func NewLevel(level string) (l Level) {
	var e bool

	if l, e = levelString[level]; !e {
		return
	}

	return
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
	_ int8 = iota
	TFile
	TSmart
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
		"info":      Info,
		"info+":     LevelAll & ^Debug,
		"notice":    Notice,
		"notice+":   LevelAll & ^Debug & ^Info,
		"warning":   Warning,
		"warning+":  LevelAll & ^Debug & ^Info & ^Notice,
		"error":     Error,
		"error+":    Error | Critical | Alert | Emergency,
		"critical":  Critical,
		"alert":     Alert,
		"emergency": Emergency,
	}
)
