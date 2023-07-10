package util

import "log"

type StreamLogger struct {
	Logger *log.Logger
}

func (sl *StreamLogger) Write(p []byte) (n int, err error) {
	sl.Logger.Print(string(p))
	return len(p), nil
}
