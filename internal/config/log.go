package config

import (
	"github.com/hs-zavet/comtools/logkit"
	"github.com/sirupsen/logrus"
)

func Logger(cfg Config) *logrus.Logger {
	logger := logkit.SetupLogger(cfg.Server.Log.Level, cfg.Server.Log.Format)
	return logger
}
