package utils

import (
	log "github.com/sirupsen/logrus"
)

func ErrorHandler(msg error) {
	log.Fatalln(msg.Error())
}
