package logger

import (
	"log"
	"runtime"
	"strconv"
)

var LogLevel string

func LogFatal(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	pre := []string{"\033[31m", "FATAL: ", "\033[0m", file, strconv.Itoa(line)}
	preI := make([]interface{}, len(pre))
	for i, s := range pre {
		preI[i] = s
	}
	v = append(preI, v...)
	log.Fatalln(v...)
}

func LogInfo(v ...interface{}) {
	pre := []string{"\033[32m", "INFO: ", "\033[0m"}
	preI := make([]interface{}, len(pre))
	for i, s := range pre {
		preI[i] = s
	}
	v = append(preI, v...)
	log.Println(v...)
}

func LogError(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	pre := []string{"\033[31m", "ERROR: ", "\033[0m", file, strconv.Itoa(line)}
	preI := make([]interface{}, len(pre))
	for i, s := range pre {
		preI[i] = s
	}
	v = append(preI, v...)
	log.Println(v...)
}

func LogWarn(v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	pre := []string{"\033[33m", "WARN: ", "\033[0m", file, strconv.Itoa(line)}
	preI := make([]interface{}, len(pre))
	for i, s := range pre {
		preI[i] = s
	}
	v = append(preI, v...)
	log.Println(v...)
}

func LogDebug(v ...interface{}) {
	if LogLevel == "debug" {
		pre := []string{"\033[34m", "DEBUG: ", "\033[0m"}
		preI := make([]interface{}, len(pre))
		for i, s := range pre {
			preI[i] = s
		}
		v = append(preI, v...)
		log.Println(v...)
	}
}
