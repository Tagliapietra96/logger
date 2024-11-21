package logs

import (
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
)

func NewOpts(configurations ...OptionConfiguration) *Options {
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

func Deb(opts *Options, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	return createNewLog(opts, DEBUG, caller, formattedMessage)
}

func Inf(opts *Options, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	return createNewLog(opts, INFO, caller, formattedMessage)
}

func War(opts *Options, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	return createNewLog(opts, WARNING, caller, formattedMessage)
}

func Err(opts *Options, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	return createNewLog(opts, ERROR, caller, formattedMessage)
}

func Fat(opts *Options, e error) error {
	if e == nil {
		return nil
	}

	caller, err := getCaller()
	if err != nil {
		return err
	}

	err = createNewLog(opts, FATAL, caller, e.Error())
	if err != nil {
		return err
	}

	beeep.Alert(opts.GetDefaultFatalTitle(), opts.GetDefaultFatalMessage(), "")
	os.Exit(1)
	return nil
}

func PrintDeb(opts *Options, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
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
	return nil
}

func PrintInf(opts *Options, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
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
	return nil
}

func PrintWar(opts *Options, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
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
	return nil
}

func PrintErr(opts *Options, message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
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
	return nil
}

func PrintFat(opts *Options, e error) error {
	if e == nil {
		return nil
	}

	caller, err := getCaller()
	if err != nil {
		return err
	}

	l := &Log{
		Status:         FATAL,
		Context:        opts.GetContext(),
		CallerFile:     caller.file,
		CallerLine:     caller.line,
		CallerFunction: caller.funcion,
		Message:        e.Error(),
		Timestamp:      time.Now().Format("2006-01-02 15:04:05"),
	}

	printLogs([]*Log{l})
	os.Exit(1)
	return nil
}

func PrintLogs(queryConfigs ...QueryConfiguration) error {
	logs, err := queryLogs(queryConfigs...)
	if err != nil {
		return err
	}

	printLogs(logs)
	return nil
}
