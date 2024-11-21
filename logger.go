package logs

import (
	"fmt"
	"os"
	"time"
)

func NewLogOpts(configurations ...OptionConfiguration) *Options {
	opts := new(Options)
	opts.UseBinaryFolder(false)
	opts.Context("")
	opts.DefaultFatalMessage("An error occurred, please check the logs for more information")
	opts.DefaultFatalTitle("FATAL")

	for _, config := range configurations {
		config(opts)
	}

	return opts
}

func Deb(opts *Options, message string, args ...any) {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	err = createNewLog(opts, DEBUG, caller, formattedMessage)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
}

func Inf(opts *Options, message string, args ...any) {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	err = createNewLog(opts, INFO, caller, formattedMessage)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
}

func War(opts *Options, message string, args ...any) {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	err = createNewLog(opts, WARNING, caller, formattedMessage)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
}

func Err(opts *Options, message string, args ...any) {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	err = createNewLog(opts, ERROR, caller, formattedMessage)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}
}

func Fat(opts *Options, err error) {
	if err == nil {
		return
	}

	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	err = createNewLog(opts, FATAL, caller, err.Error())
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	}

	os.Exit(1)
}

func PrintDeb(opts *Options, message string, args ...any) {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	l := &Log{
		Status:         DEBUG,
		Context:        opts.GetContext(),
		CallerFile:     caller.file,
		CallerLine:     caller.line,
		CallerFunction: caller.funcion,
		Message:        formattedMessage,
		Timestamp:      time.Now().Format("2006-01-02 15:04:05"),
	}

	printLogs([]*Log{l})
}

func PrintInf(opts *Options, message string, args ...any) {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	l := &Log{
		Status:         INFO,
		Context:        opts.GetContext(),
		CallerFile:     caller.file,
		CallerLine:     caller.line,
		CallerFunction: caller.funcion,
		Message:        formattedMessage,
		Timestamp:      time.Now().Format("2006-01-02 15:04:05"),
	}

	printLogs([]*Log{l})
}

func PrintWar(opts *Options, message string, args ...any) {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	l := &Log{
		Status:         WARNING,
		Context:        opts.GetContext(),
		CallerFile:     caller.file,
		CallerLine:     caller.line,
		CallerFunction: caller.funcion,
		Message:        formattedMessage,
		Timestamp:      time.Now().Format("2006-01-02 15:04:05"),
	}

	printLogs([]*Log{l})
}

func PrintErr(opts *Options, message string, args ...any) {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	l := &Log{
		Status:         ERROR,
		Context:        opts.GetContext(),
		CallerFile:     caller.file,
		CallerLine:     caller.line,
		CallerFunction: caller.funcion,
		Message:        formattedMessage,
		Timestamp:      time.Now().Format("2006-01-02 15:04:05"),
	}

	printLogs([]*Log{l})
}

func PrintFat(opts *Options, err error) {
	if err == nil {
		return
	}

	caller, err := getCaller()
	if err != nil {
		fmt.Println("ERROR: [logger-pkg] failed to get the caller information")
	}

	l := &Log{
		Status:         FATAL,
		Context:        opts.GetContext(),
		CallerFile:     caller.file,
		CallerLine:     caller.line,
		CallerFunction: caller.funcion,
		Message:        err.Error(),
		Timestamp:      time.Now().Format("2006-01-02 15:04:05"),
	}

	printLogs([]*Log{l})
	os.Exit(1)
}

func PrintLogs(queryConfigs ...QueryConfiguration) {
	logs, err := queryLogs(queryConfigs...)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}

	printLogs(logs)
}
