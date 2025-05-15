package controllers

import (
	"github.com/sirupsen/logrus"
)

type Controller struct {
	log *logrus.Entry
}

func NewController(log *logrus.Entry) Controller {
	return Controller{
		log: log,
	}
}
