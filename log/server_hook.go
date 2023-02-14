package log

import "github.com/sirupsen/logrus"

type AppHook struct {
	Name string
}

func NewAppHook(name string) *AppHook {
	return &AppHook{Name: name}
}

func (s *AppHook) Fire(entry *logrus.Entry) error {
	entry.Data["srv_name"] = s.Name
	return nil
}

func (s *AppHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
