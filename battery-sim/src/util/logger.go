package util

import (
	"fmt"
	"log"
)

type Logger struct {
	correlationID string
}

func NewLogger(correlationID string) Logger {
	return Logger{correlationID: correlationID}
}

func (l *Logger) logf(format string, params ...interface{}) {
	log.Printf("[%s] "+format, append([]interface{}{l.correlationID}, params...)...)
}

func (l *Logger) Info(format string, params ...interface{}) {
	log.Printf("[%s] "+format, append([]interface{}{l.correlationID}, params...)...)
}

func (l *Logger) Fatal(err error) {
	log.Fatalf("[%s] %v", l.correlationID, err)
}

func (l *Logger) Fatalf(err error, format string, params ...interface{}) {
	msg := fmt.Sprintf(format, params...)
	log.Fatalf("[%s] %s %v", l.correlationID, msg, err)
}

func (l *Logger) ErrorWithMsg(msg string, err error, params ...interface{}) {
	l.logf("%s: %v", append([]interface{}{msg, err}, params...)...)
}
