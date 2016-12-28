package log

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	return &Logger{
		Logger: logrus.New(),
	}
}

