package util

import "log"

type StreamLogger struct {
	logger *log.Logger
}

func (sl *StreamLogger) Write(p []byte) (n int, err error) {
	sl.logger.Print(string(p))
	return len(p), nil
}
