package logx

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogOpt struct {
	LogPath string // 日志路径
	Env     string // 环境：dev开发环境 prod生产环境
}

type Log struct {
	*logrus.Logger
}

func New(opt LogOpt) *Log {
	l := logrus.New()
	if opt.Env == "prod" {
		logger := &lumberjack.Logger{
			Filename:   opt.LogPath,
			MaxSize:    1, // MB
			MaxBackups: 3,
			MaxAge:     28, // Days
		}
		l.SetOutput(logger)
		l.SetFormatter(&logrus.JSONFormatter{})
	} else if opt.Env == "dev" {
		l.SetLevel(logrus.DebugLevel)
	}
	return &Log{l}
}

func (l *Log) WrapError(err error) error {
	if err != nil {
		l.Error(err.Error())
	}
	return err
}
