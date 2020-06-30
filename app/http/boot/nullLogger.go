package h_boot

import (
	"log"
	"os"
)

func NullLogger() *log.Logger {
	devNull, _ := os.Open(os.DevNull)
	return log.New(devNull, "", log.LstdFlags)
}