package config

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Config) GetLogger() *logrus.Logger {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(strings.ToLower(c.Logger.Level))
	if err != nil {
		logger.Warnf("invalid log level '%s', defaulting to 'info'", c.Logger.Level)
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	switch strings.ToLower(c.Logger.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		fallthrough
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return logger
}
