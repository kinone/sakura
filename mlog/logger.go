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

	AddHandler(Handler)
}

type Option struct {
	Type  int8
	Level string
	File  string
}

type Logger struct {
	handlers []Handler
}

func NewLogger(opt *Option) (l *Logger) {
	if nil == opt {
		opt = &Option{}
	}

	var handler Handler

	switch opt.Type {
	case TFile:
		handler = NewLevelhandler(opt.File, opt.Level)
	case TNull:
		handler = &NullHandler{}
	default:
		handler = NewSmartHandler(NewLevelhandler(opt.File, opt.Level))
	}

	l = &Logger{
		handlers: []Handler{handler},
	}

	return
}

func (l *Logger) AddHandler(h Handler) {
	l.handlers = append(l.handlers, h)
}

func (l *Logger) Print(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{NoLevel, "", v})
	}
}

func (l *Logger) Println(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{NoLevel, "", v})
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{NoLevel, format, v})
	}
}

func (l *Logger) Debug(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Debug, "", v})
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Debug, format, v})
	}
}

func (l *Logger) Info(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Info, "", v})
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Info, format, v})
	}
}

func (l *Logger) Notice(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Notice, "", v})
	}
}

func (l *Logger) Noticef(format string, v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Notice, format, v})
	}
}

func (l *Logger) Warning(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Warning, "", v})
	}
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Warning, format, v})
	}
}

func (l *Logger) Error(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Error, "", v})
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Error, format, v})
	}
}

func (l *Logger) Alert(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Alert, "", v})
	}
}

func (l *Logger) Alertf(format string, v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Alert, format, v})
	}
}

func (l *Logger) Critical(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Critical, "", v})
	}
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Critical, format, v})
	}
}

func (l *Logger) Emergency(v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Emergency, "", v})
	}
}

func (l *Logger) Emergencyf(format string, v ...interface{}) {
	for _, h := range l.handlers {
		h.Log(&Record{Emergency, format, v})
	}
}

func (l *Logger) Reload() (err error) {
	for _, h := range l.handlers {
		if err = h.Reload(); nil != err {
			return
		}
	}

	return
}

func (l *Logger) Close() {
	for _, h := range l.handlers {
		h.Close()
	}
}
