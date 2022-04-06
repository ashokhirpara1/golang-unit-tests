package prettylogs

import "github.com/sirupsen/logrus"

// Storage - hold logick for logging
type Storage interface {
	Info(fields logrus.Fields, msg string)
	Error(fields logrus.Fields, msg string)
	Fatal(fields logrus.Fields, msg string)
}
