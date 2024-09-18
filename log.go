package cache

import "github.com/anchore/go-logger"

func traceLog(log logger.Logger, message string, fields ...any) {
	if log == nil {
		return
	}
	log.WithFields(fields...).Trace(message)
}

func warnLog(log logger.Logger, message string, fields ...any) {
	if log == nil {
		return
	}
	log.WithFields(fields...).Warn(message)
}
