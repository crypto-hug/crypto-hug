package log

import (
	"github.com/sirupsen/logrus"
	//	"gopkg.in/natefinch/lumberjack.v2"
	//	"os"
)

var consoleLogger = newConsoleLogger()

// func ChangeSetting() {
// 	// log.SetFormatter(&log.TextFormatter{
// 	// 	FullTimestamp: true})

// 	// log.SetOutput(&lumberjack.Logger{
// 	// 	Filename:   "/Users/vbraun/code/repos-public/crypto-hug/foo.log",
// 	// 	MaxSize:    500, // megabytes
// 	// 	MaxBackups: 3,
// 	// 	MaxAge:     28,   //days
// 	// 	Compress:   true, // disabled by default
// 	// })

// 	// //log.SetOutput(os.Stdout)
// 	//consoleLogger.Formatter = &logrus.TextFormatter{}
// }

func newConsoleLogger() *logrus.Logger {
	result := logrus.New()
	result.Formatter = defaultFormatter()
	result.Level = logrus.DebugLevel
	return result
}

func defaultFormatter() *logrus.TextFormatter {
	return &logrus.TextFormatter{}
}

type Logger struct {
	innerLogger []*logrus.Entry
}

type More logrus.Fields

func new(fields logrus.Fields) *Logger {
	consoleLogger.Formatter = defaultFormatter()
	consoleEntry := consoleLogger.WithFields(fields)

	result := Logger{innerLogger: []*logrus.Entry{consoleEntry}}
	return &result
}

func (self *Logger) log(lvl logrus.Level, msg string, more More) {
	for _, innerLogger := range self.innerLogger {
		if more != nil && len(more) > 0 {
			innerLogger = innerLogger.WithFields(logrus.Fields(more))
		}

		if lvl == logrus.DebugLevel {
			innerLogger.Debug(msg)
		} else if lvl == logrus.InfoLevel {
			innerLogger.Info(msg)
		} else if lvl == logrus.WarnLevel {
			innerLogger.Warn(msg)
		} else if lvl == logrus.ErrorLevel {
			innerLogger.Error(msg)
		} else if lvl == logrus.PanicLevel {
			innerLogger.Panic(msg)
		} else if lvl == logrus.FatalLevel {
			innerLogger.Fatal(msg)
		}
	}
}

func Global() *Logger {
	return NewLog("global")
}

func NewLog(sys string) *Logger {
	result := new(logrus.Fields{"sys": sys})

	return result
}

func NewLogWithFields(sys string, fields logrus.Fields) *Logger {
	if fields == nil {
		return NewLog(sys)
	}

	fields["sys"] = sys
	return new(fields)
}

func (self *Logger) Debug(msg string, more More) {
	self.log(logrus.DebugLevel, msg, more)
}

func (self *Logger) Info(msg string, more More) {
	self.log(logrus.InfoLevel, msg, more)
}

func (self *Logger) Warn(msg string, more More) {
	self.log(logrus.WarnLevel, msg, more)
}

func (self *Logger) Error(msg string, more More) {
	self.log(logrus.ErrorLevel, msg, more)
}

func (self *Logger) Panic(msg string, more More) {
	self.log(logrus.PanicLevel, msg, more)
}

func (self *Logger) Fatal(msg string, more More) {
	self.log(logrus.FatalLevel, msg, more)
}
