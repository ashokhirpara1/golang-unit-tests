package prettylogs

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
	"unit-tests/prettylogs/logger"

	"github.com/sirupsen/logrus"
)

// Handler represented custom storage
type Handler struct {
	storage *logger.Handler
}

// Get create new Handler
func Get() *Handler {
	logger := logger.New()
	return &Handler{storage: logger}
}

// Info log handler
func (lg *Handler) Info(str string) {

	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a Function object this functions parent
	functionObject := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	method := extractFnName.ReplaceAllString(functionObject.Name(), "$1")

	fields := logrus.Fields{
		"method": method,
	}

	lg.storage.Info(fields, str)
}

// Error - log error
func (lg *Handler) Error(msg string, errr error) {
	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a Function object this functions parent
	functionObject := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	method := extractFnName.ReplaceAllString(functionObject.Name(), "$1")

	fields := logrus.Fields{
		"method": method,
	}

	msg = msg + ": " + errr.Error()
	lg.storage.Error(fields, msg)
}

// DBError - error in DB, send msg to slack
func (lg *Handler) DBError(msg string, errr error) {
	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a Function object this functions parent
	functionObject := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	method := extractFnName.ReplaceAllString(functionObject.Name(), "$1")

	fields := logrus.Fields{
		"method": method,
	}

	msg = msg + ": " + errr.Error()
	lg.storage.Error(fields, msg)

}

// Fatal - log fatal
func (lg *Handler) Fatal(method string, msg string, errr error) {
	fields := logrus.Fields{
		"method": method,
	}

	msg = msg + ": " + errr.Error()
	lg.storage.Fatal(fields, msg)
}

// Enter - start time of working of method
func (lg *Handler) Enter() (time.Time, string) {
	start := time.Now()

	// Skip this function, and fetch the PC and file for its parent
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a Function object this functions parent
	functionObject := runtime.FuncForPC(pc)

	// Regex to extract just the function name (and not the module path)
	extractFnName := regexp.MustCompile(`^.*\.(.*)$`)
	method := extractFnName.ReplaceAllString(functionObject.Name(), "$1")

	fields := logrus.Fields{
		"method": method,
	}

	msg := "Started " + method
	lg.storage.Info(fields, msg)

	return start, method
}

// Exit this is to track how much time function took
func (lg *Handler) Exit(start time.Time, name string) {
	elapsed := time.Since(start)

	fields := logrus.Fields{
		"method": name,
	}

	msg := fmt.Sprintf("%s took %s", name, elapsed)
	lg.storage.Info(fields, msg)
}
