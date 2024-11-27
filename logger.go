package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
)

func New(useBinaryFolder bool, tags ...string) *Logger {
	l := new(Logger)
	l.useBinaryFolder = useBinaryFolder
	l.showCaller = ShowCallerFile
	l.showTimestamp = ShowDateTime
	l.showTags = false
	l.fatalTitle = "Fatal"
	l.fatalMessage = "An error occurred, please check the logs for more information"
	l.tags = make([]string, 0)

	if len(tags) > 0 {
		l.showTags = true
		l.tags = tags
	}

	return l
}

type Logger struct {
	useBinaryFolder bool
	showTags        bool
	showCaller      ShowCallerLevel
	showTimestamp   ShowTimestampLevel
	tags            []string
	fatalTitle      string
	fatalMessage    string
}

func (opts *Logger) BinaryFolder(use bool) {
	opts.useBinaryFolder = use
}

func (opts *Logger) Caller(level ShowCallerLevel) {
	opts.showCaller = level
}

func (opts *Logger) Timestamp(level ShowTimestampLevel) {
	opts.showTimestamp = level
}

func (opts *Logger) ShowTags(show bool) {
	opts.showTags = show
}

func (opts *Logger) Tags(tags ...string) {
	opts.tags = append(opts.tags, tags...)
}

func (opts *Logger) SetFatal(title, message string) {
	opts.fatalTitle = title
	opts.fatalMessage = message
}

func (opts *Logger) Debug(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	return createNewLog(opts, Debug, caller, formattedMessage)
}

func (opts *Logger) Info(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	return createNewLog(opts, Info, caller, formattedMessage)
}

func (opts *Logger) Warn(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	return createNewLog(opts, Warning, caller, formattedMessage)
}

func (opts *Logger) Error(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	return createNewLog(opts, Error, caller, formattedMessage)
}

func (opts *Logger) Fatal(e error) error {
	if e == nil {
		return nil
	}

	caller, err := getCaller()
	if err != nil {
		return err
	}

	err = createNewLog(opts, Fatal, caller, e.Error())
	if err != nil {
		return err
	}

	beeep.Alert(opts.fatalTitle, opts.fatalMessage, "")
	os.Exit(1)
	return nil
}

func (opts *Logger) PrintDebug(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	l := &log{
		level:          Debug,
		tags:           opts.tags,
		callerFile:     caller.file,
		callerLine:     caller.line,
		callerFunction: caller.funcion,
		message:        formattedMessage,
		timestamp:      timestamp(time.Now()),
	}

	printLogs([]*log{l})
	return nil
}

func (opts *Logger) PrintInfo(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	l := &log{
		level:          Info,
		tags:           opts.tags,
		callerFile:     caller.file,
		callerLine:     caller.line,
		callerFunction: caller.funcion,
		message:        formattedMessage,
		timestamp:      timestamp(time.Now()),
	}

	printLogs([]*log{l})
	return nil
}

func (opts *Logger) PrintWarn(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	l := &log{
		level:          Warning,
		tags:           opts.tags,
		callerFile:     caller.file,
		callerLine:     caller.line,
		callerFunction: caller.funcion,
		message:        formattedMessage,
		timestamp:      timestamp(time.Now()),
	}

	printLogs([]*log{l})
	return nil
}

func (opts *Logger) PrintError(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	caller, err := getCaller()
	if err != nil {
		return err
	}

	l := &log{
		level:          Error,
		tags:           opts.tags,
		callerFile:     caller.file,
		callerLine:     caller.line,
		callerFunction: caller.funcion,
		message:        formattedMessage,
		timestamp:      timestamp(time.Now()),
	}

	printLogs([]*log{l})
	return nil
}

func (opts *Logger) PrintFatal(e error) error {
	if e == nil {
		return nil
	}

	caller, err := getCaller()
	if err != nil {
		return err
	}

	l := &log{
		level:          Fatal,
		tags:           opts.tags,
		callerFile:     caller.file,
		callerLine:     caller.line,
		callerFunction: caller.funcion,
		message:        e.Error(),
		timestamp:      timestamp(time.Now()),
	}

	printLogs([]*log{l})
	os.Exit(1)
	return nil
}

func (opts *Logger) PrintLogs(queryConfigs ...QueryConfiguration) error {
	logs, err := queryLogs(queryConfigs...)
	if err != nil {
		return err
	}

	printLogs(logs)
	return nil
}
