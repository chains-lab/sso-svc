package config

import (
	"strings"

	"github.com/chains-lab/sso-svc/internal/logger"
	"github.com/sirupsen/logrus"
)

func (c *Config) GetLogger() logger.Logger {
	base := logrus.New()

	lvl, err := logrus.ParseLevel(strings.ToLower(c.Logger.Level))
	if err != nil {
		base.Warnf("invalid log level '%s', defaulting to 'info'", c.Logger.Level)
		lvl = logrus.InfoLevel
	}
	base.SetLevel(lvl)

	switch strings.ToLower(c.Logger.Format) {
	case "json":
		base.SetFormatter(&logrus.JSONFormatter{})
	default:
		base.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	return logger.NewWithBase(base)
}
