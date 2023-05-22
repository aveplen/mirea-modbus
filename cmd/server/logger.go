package main

import "fmt"

type Logger struct {
	staticPrefix  string
	dynamicPrefix func() string
	buffer        []string

	size   int
	window int
}

func NewLogger(
	staticPrefix string,
	dynamicPrefix func() string,
	size int,
	window int,
) *Logger {
	return &Logger{
		staticPrefix:  staticPrefix,
		dynamicPrefix: dynamicPrefix,
		buffer:        make([]string, 0, size),
		window:        window,
		size:          size,
	}
}

func (l *Logger) GetBuffer() []string {
	return append(make([]string, 0, len(l.buffer)), l.buffer...)
}

func (l *Logger) logMessage(level string, message string) {
	if len(l.buffer) == cap(l.buffer) {
		l.buffer = append(make([]string, 0, l.size), l.buffer[:l.window]...)
	}

	l.buffer = append(
		l.buffer,
		fmt.Sprintf(
			"%s [%s] %s %s",
			l.dynamicPrefix(),
			level,
			l.staticPrefix,
			message,
		),
	)
}

func (l *Logger) Error(err error) {
	l.logMessage("EROR", err.Error())
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Errorf(format, args...))
}

func (l *Logger) Info(message string) {
	l.logMessage("INFO", message)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(message string) {
	l.logMessage("DEBG", message)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(message string) {
	l.logMessage("WARN", message)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}
