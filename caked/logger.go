package main

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
}
