package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

//Logger Defines a Logger message definition
type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

//NewLogger Returns a customer Logger Formatter using the logurs logger.
func NewLogger(wr io.Writer, level string, format string) Logger {
	if wr == nil {
		wr = os.Stderr
	}

	lr := logrus.New()
	lr.SetOutput(wr)
	lr.SetFormatter(&logrus.TextFormatter{})
	if format == "json" {
		lr.SetFormatter(&logrus.JSONFormatter{})
	}

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.WarnLevel
		lr.Warnf("failed to parse log-level '%s', defaulting to 'warning'", level)
	}
	lr.SetLevel(lvl)

	return &logrusLogger{
		Entry: logrus.NewEntry(lr),
	}
}

type logrusLogger struct {
	*logrus.Entry
}