package mlog

const (
	Debug = 1 << iota
	Info
	Notice
	Warning
	Error
	Critical
	Alert
	Emergency
)

const LevelAll = Debug | Info | Notice | Warning | Error | Critical | Alert | Emergency

const (
	_ int8 = iota
	TFile
	TSmart
	TMultiHandler
)

var (
	levelPrefix = map[int]string{
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

	levelString = map[string]int{
		"debug":     Debug,
		"info":      Info,
		"notice":    Notice,
		"warning":   Warning,
		"error":     Error,
		"critical":  Critical,
		"alert":     Alert,
		"emergency": Emergency,
	}
)

func ConvertLogLevel(level string) (l int) {
	var e bool

	if l, e = levelString[level]; !e {
		return
	}

	return l
}

func Prefix(l int) string {
	return levelPrefix[l]
}
