package logger

import (
	"github.com/sirupsen/logrus"

	"bytes"
	"os"
)

//Handler - hold logrus
type Handler struct {
	*logrus.Logger
}

type logOutputSplitter struct{}

// Splits log output, error and fatal to stderr and the rest to stdout
func (splitter *logOutputSplitter) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte("level=debug")) || bytes.Contains(p, []byte("level=info")) ||
		bytes.Contains(p, []byte("level=trace")) || bytes.Contains(p, []byte("level=warn")) {
		return os.Stdout.Write(p)
	}
	return os.Stderr.Write(p)
}

// New - return Handler
func New() *Handler {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// logger the debug severity or above.
	logger.SetLevel(logrus.DebugLevel)

	return &Handler{logger}
}

// Info handler
func (l *Handler) Info(fields logrus.Fields, msg string) {
	l.Logger.WithFields(fields).Info(msg)
}

// Error handler
func (l *Handler) Error(fields logrus.Fields, msg string) {
	l.Logger.WithFields(fields).Error(msg)
}

// Fatal handler
func (l *Handler) Fatal(fields logrus.Fields, msg string) {
	l.Logger.WithFields(fields).Fatal(msg)
}
