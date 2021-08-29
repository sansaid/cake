package utils

import (
	"log"
)

func Check(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
