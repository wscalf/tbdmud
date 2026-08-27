package game

import (
	"fmt"
	"log/slog"
)

type Log struct{}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Debug(format string, args ...any) {
	message := fmt.Sprintf(format, args...) //TODO: condition formatting on whether or not this log level is enabled. Also, log message queue?
	slog.Debug(message)
}
