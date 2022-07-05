package lib

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type Logger struct{}

func exists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func createLogger(logType string, caller string) *log.Logger {
	now := time.Now().Format("2006-01-02")
	path := fmt.Sprintf("logs/%s", now)

	if !exists(path) {
		err := os.MkdirAll(path, os.ModePerm)

		if err != nil {
			log.Fatalf("Error create \"%s\" directory", path)
		}
	}

	path = fmt.Sprintf("%s/%s.log", path, logType)

	if exists(path) {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)

		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}

		return _createLogger(file, caller)
	} else {
		file, err := os.Create(path)

		if err != nil {
			log.Fatalf("Error create \"%s\" directory", path)
		}

		return _createLogger(file, caller)
	}
}

func _createLogger(file *os.File, caller string) *log.Logger {
	now := time.Now().Format("2006-01-02 15:04:05")

	return log.New(file, fmt.Sprintf("%s - %s - ", now, caller), 0)
}

func (Logger) Info(message string) {
	caller := getCaller()

	_logger := createLogger("info", caller)
	_logger.Println(message)
}

func (Logger) Error(message string, exit bool) {
	caller := getCaller()

	_logger := createLogger("error", caller)
	_logger.Println(message)

	if exit {
		os.Exit(1)
	}
}

func (Logger) InfoLine(message string) {
	caller := getCaller()

	_logger := createLogger("infoLine", caller)
	_logger.Println(message)
}

func (Logger) InfoFreshdesk(message string) {
	caller := getCaller()

	_logger := createLogger("infoFreshdesk", caller)
	_logger.Println(message)
}

func (Logger) InfoMiddleware(message string) {
	caller := getCaller()

	_logger := createLogger("infoMiddleware", caller)
	_logger.Println(message)
}

func getCaller() string {
	pc, _, _, ok := runtime.Caller(2)

	if ok {
		function := runtime.FuncForPC(pc)
		_, line := function.FileLine(pc)

		return fmt.Sprintf("%s:%d", function.Name(), line)
	} else {
		return "Unknown"
	}
}
