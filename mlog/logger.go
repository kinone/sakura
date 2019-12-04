package mlog

type StdLogger interface {
	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type LevelLogger interface {
	StdLogger

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Notice(v ...interface{})
	Noticef(format string, v ...interface{})

	Warning(v ...interface{})
	Warningf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Critical(v ...interface{})
	Criticalf(format string, v ...interface{})

	Alert(v ...interface{})
	Alertf(format string, v ...interface{})

	Emergency(v ...interface{})
	Emergencyf(format string, v ...interface{})

	Reload() error
	Close()
}

type Option struct {
	Type  int8
	Level string
	File  string
}

type Logger struct {
	h Handler
}

func NewLogger(opt *Option) (l *Logger) {
	if nil == opt {
		opt = &Option{}
	}

	f := func() (h *FileHandler) {
		h = NewFileHandler(opt.File)
		l := ConvertLogLevel(opt.Level)
		h.filter = LevelFilter(l)

		return h
	}

	var handler Handler

	switch opt.Type {
	case TFile:
		handler = f()
	case TNull:
		handler = &NullHandler{}
	default:
		handler = NewSmartHandler(f())
	}

	l = &Logger{
		h: handler,
	}

	return
}

func (l *Logger) Print(v ...interface{}) {
	l.h.Log(NoLevel, v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.h.Log(NoLevel, v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.h.Logf(NoLevel, format, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.h.Log(Debug, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.h.Logf(Debug, format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.h.Log(Info, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.h.Logf(Info, format, v...)
}

func (l *Logger) Notice(v ...interface{}) {
	l.h.Log(Notice, v...)
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	l.h.Logf(Notice, format, v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.h.Log(Warning, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.h.Logf(Warning, format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.h.Log(Error, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.h.Logf(Error, format, v...)
}

func (l *Logger) Alert(v ...interface{}) {
	l.h.Log(Alert, v...)
}

func (l *Logger) Alertf(format string, v ...interface{}) {
	l.h.Logf(Alert, format, v...)
}

func (l *Logger) Critical(v ...interface{}) {
	l.h.Log(Critical, v...)
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
	l.h.Logf(Critical, format, v...)
}

func (l *Logger) Emergency(v ...interface{}) {
	l.h.Log(Emergency, v...)
}

func (l *Logger) Emergencyf(format string, v ...interface{}) {
	l.h.Logf(Emergency, format, v...)
}

func (l *Logger) Reload() error {
	return l.h.Reload()
}

func (l *Logger) Close() {
	l.h.Close()
}
