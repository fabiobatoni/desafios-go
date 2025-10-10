package main

import (
	"fmt"
	"log/slog"
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

func main() {
	fmt.Println(convertLevel(DEBUG))
}

func convertLevel(l int) string {
	switch {
	case l == 0:
		CustomizeLog(l)
		return "INFO"
	case l == 1:
		CustomizeLog(l)
		return "INFO"
	case l == 2:
		CustomizeLog(l)
		return "WARNING"
	case l == 3:
		CustomizeLog(l)
		return "CRITICAL"
	default:
		return "ERROR"
	}
}

func CustomizeLog(l int) {
	logger := slog.Default()

	if l >= 2 {
		logger.Warn("Ixi azedou !")
	}
}
