package utils

import (
	"github.com/zjxpcyc/tinylogger"
)

var logger tinylogger.LogService

func SetLogger(l tinylogger.LogService) {
	logger = l
}

func GetLogger() tinylogger.LogService {
	return logger
}
