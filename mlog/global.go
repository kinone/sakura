package mlog

const (
	Debug     = 100
	Info      = 200
	Notice    = 250
	Warning   = 300
	Error     = 400
	Critical  = 500
	Alert     = 550
	Emergency = 600
	NoLevel   = 900
)

const (
	TSmart = iota
	TFile
	TNull
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
		NoLevel:   "[NL]",
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

func ConvertLogLevel(level string) int {
	l, e := levelString[level]
	if !e {
		return Debug
	}

	return l
}

func Prefix(l int) string {
	return levelPrefix[l]
}
