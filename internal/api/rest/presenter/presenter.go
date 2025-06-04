package presenter

import (
	"github.com/sirupsen/logrus"
)

type Presenter struct {
	log *logrus.Entry
}

func NewPresenter(log *logrus.Entry) Presenter {
	return Presenter{
		log: log,
	}
}
