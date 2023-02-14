package log

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Option func(*Logger)

func WithLogPath(path string) Option {
	return func(l *Logger) {
		l.path = path
	}
}

func WithLevel(level uint32) Option {
	return func(l *Logger) {
		l.SetLevel(logrus.Level(level))
	}
}

type Logger struct {
	path string
	*logrus.Logger
}

func New(opts ...Option) *Logger {
	l := logrus.New()
	_log := &Logger{
		Logger: l,
	}

	for _, opt := range opts {
		opt(_log)
	}

	if _log.path != "" {
		logger := &lumberjack.Logger{
			Filename:   _log.path,
			MaxSize:    1, // MB
			MaxBackups: 3,
			MaxAge:     28, // Days
		}
		_log.SetOutput(logger)
		_log.SetFormatter(&logrus.JSONFormatter{})
	}

	return _log
}
