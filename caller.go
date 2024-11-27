package logger

import (
	"errors"
	"runtime"
)

// caller is a utility struct to store the caller information
type caller struct {
	file    string
	line    int
	funcion string
}

// getCaller returns the caller information, such as the file, line and function
func getCaller() (*caller, error) {
	c := new(caller) // setup the caller struct
	c.file = "unknown"
	c.line = 0
	c.funcion = "unknown"

	// get the caller information by runtime
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return c, errors.New("[logger-pkg] failed to get the caller information")
	}

	// update the caller struct with the information
	c.file = file
	c.line = line

	f := runtime.FuncForPC(pc) // get the function information
	c.funcion = f.Name()       // update the caller struct with the function name
	return c, nil
}

// ShowCallerLevel is an enum to define the level of caller information to be shown
type ShowCallerLevel int

const (
	HideCaller         ShowCallerLevel = iota // hide the caller information
	ShowCallerFile                            // show the caller file only main.go
	ShowCallerLine                            // show the caller file and line main.go:10
	ShowCallerFunction                        // show the caller file, line and function main.go:10 - main.main
)
