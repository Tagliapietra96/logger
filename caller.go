package logger

import (
	"errors"
	"path/filepath"
	"runtime"
)

// ShowCallerLevel is an enum to define the level of caller information to be shown
type ShowCallerLevel int

const (
	HideCaller         ShowCallerLevel = iota // hide the caller information
	ShowCallerFile                            // show the caller file only main.go
	ShowCallerLine                            // show the caller file and line main.go:10
	ShowCallerFunction                        // show the caller file, line and function main.go:10 - main.main
)

// getCaller appends the caller information to a log, such as the file, line and function
func getCaller(l *log) error {
	// get the caller information by runtime
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return errors.New("[logger-pkg] failed to get the caller information")
	}

	l.callerFile = filepath.Base(file)
	l.callerLine = line

	f := runtime.FuncForPC(pc)
	l.callerFunction = f.Name()
	return nil
}
