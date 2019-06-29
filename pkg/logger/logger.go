package logger

import "fmt"

var (
	logLevelMap = map[string]int {
		debugLevel: 1,
		traceLevel: 2,
		infoLevel:  3,
		warnLevel:  4,
		errLevel:   5,
	}
)

const (
	debugLevel = "DEBUG"
	traceLevel = "TRACE"
	infoLevel  = "INFO"
	warnLevel  = "WARN"
	errLevel   = "ERROR"

	LogLevel = infoLevel
)

func LogIfErr(err error, a ...interface{}) {
	if err != nil {
		Error(err, a)
	}
}

func Debug(format string, a ...interface{}) {
	log(debugLevel, format, a...)
}

func Trace(err error, a ...interface{}) {
	log(traceLevel, err.Error(), a...)
}

func Info(format string, a ...interface{}) {
	log(infoLevel, format, a...)
}

func Warn(err error, a ...interface{}) {
	log(warnLevel, err.Error(), a...)
}

func Error(err error, a ...interface{}) {
	log(errLevel, err.Error(), a...)
}

func log(logLevel string, format string, a ...interface{}) {
	if logLevelMap[logLevel] >= logLevelMap[LogLevel] {
		fmt.Println("[" + logLevel + "] " + fmt.Sprintf(format, a...))
	}
}
