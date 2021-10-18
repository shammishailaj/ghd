package utils

import (
	log "github.com/sirupsen/logrus"
)

type Utils struct {
	Log *log.Logger
}

func NewUtils(l *log.Logger) *Utils {
	return &Utils{
		Log: l,
	}
}