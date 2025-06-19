package presenter

import "github.com/sirupsen/logrus"

type Presenters struct {
	log *logrus.Entry
}

func NewPresenters(log *logrus.Entry) Presenters {
	return Presenters{
		log: log,
	}
}
