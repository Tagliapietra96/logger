package logger

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

// Logger represents the logger configuration structure
// that holds the configuration options for the logger
// with the methods to interact with the logger
// and log messages
// The logger can be configured with the following options:
//   - Folder: (string) the folder path to store the logs data (by default it uses the binary folder)
//     to store the database file, otherwise it will use the current working directory
//   - Inline: (bool) if true the logs will be printed inline, otherwise they will be printed in a block
//   - Caller: (ShowCallerLevel) the level of caller information to show
//   - Timestamp: (ShowTimestampLevel) the level of timestamp information to show
//   - ShowTags: (bool) if true the logger will show the tags in the logs
//   - Tags: (string) the tags to add to the logs created with this logger
//   - SetFatal: (string, string) the title and message to show in the fatal error
//     alert when the Fatal method is called
//   - Copy: creates a copy of the logger with the same configurations
//
// The logger has the following methods to log messages:
//   - Debug: creates a debug log message in the database (it not will be printed)
//   - Info: creates an info log message in the database (it not will be printed)
//   - Warn: creates a warning log message in the database (it not will be printed)
//   - Error: creates an error log message in the database (it not will be printed)
//   - Fatal: creates a fatal log message in the database and exits the program (it not will be printed)
//     it will show an alert with the title and message set with SetFatal (only if the error passed is not nil)
//   - PrintDebug: prints a debug log message in the console (it not will be saved in the database)
//   - PrintInfo: prints an info log message in the console (it not will be saved in the database)
//   - PrintWarn: prints a warning log message in the console (it not will be saved in the database)
//   - PrintError: prints an error log message in the console (it not will be saved in the database)
//   - PrintFatal: prints a fatal log message in the console and exits the program (it not will be saved in the database)
//     if the error passed is not nil
//   - PrintLogs: prints the logs in the database based on the query configurations passed
type Logger struct {
	folderPath    string             // the folder path to store the logs data
	showTags      bool               // if true the logger will show the tags in the logs
	inline        bool               // if true the logs will be printed inline, otherwise they will be printed in a block
	showCaller    ShowCallerLevel    // the level of caller information to show
	showTimestamp ShowTimestampLevel // the level of timestamp information to show
	tags          []string           // the tags to add to the logs created with this logger
	fatalTitle    string             // the title to show in the fatal error alert
	fatalMessage  string             // the message to show in the fatal error alert
}

// New creates a new logger with the given tags
// the tags will be added to the logs created with this logger
// if no tags are passed it will create a logger without tags
// The new logger will have the following default configurations:
//   - folderPath: the bynary folder path (if it fails to get the path it will use an empty string)
//   - showTags: false
//   - inline: false
//   - showCaller: ShowCallerFile
//   - showTimestamp: ShowDateTime
//   - fatalTitle: "Fatal"
//   - fatalMessage: "An error occurred, please check the logs for more information"
//   - tags: the tags passed or an empty slice
//
// Check the Logger struct for more information about the logger configurations
// and the methods to interact with the logger and log messages
func New(tags ...string) *Logger {
	l := new(Logger)

	folder, err := os.Executable()
	if err != nil {
		folder = ""
	}

	if strings.Contains(folder, os.TempDir()) {
		folder, err = os.Getwd()
		if err != nil {
			folder = ""
		}
	}

	l.folderPath = folder
	l.showCaller = ShowCallerFile
	l.showTimestamp = ShowDateTime
	l.showTags = false
	l.fatalTitle = "Fatal"
	l.fatalMessage = "An error occurred, please check the logs for more information"
	l.tags = make([]string, 0)

	if len(tags) > 0 {
		l.tags = tags
	}

	return l
}

// Copy creates a copy of the logger with the same configurations
func (opts *Logger) Copy() *Logger {
	l := new(Logger)
	l.folderPath = opts.folderPath
	l.showTags = opts.showTags
	l.inline = opts.inline
	l.showCaller = opts.showCaller
	l.showTimestamp = opts.showTimestamp
	l.tags = append(make([]string, 0), opts.tags...)
	l.fatalTitle = opts.fatalTitle
	l.fatalMessage = opts.fatalMessage
	return l
}

// Folder sets the folder path to store the logs data
// Every log created with this logger will be stored in this folder
func (opts *Logger) Folder(path string) {
	opts.folderPath = path
}

// Inline sets the logger to print the logs inline
// if the inline parameter is true, otherwise it will print
// the logs in a block (like cards)
func (opts *Logger) Inline(inline bool) {
	opts.inline = inline
}

// Caller sets the level of caller information to show
// in the logs based on the level parameter
// the level can be one of the following:
//   - ShowCallerFile: shows the file of the caller
//   - ShowCallerLine: shows the caller file and line
//   - ShowCallerFunction: shows the caller file, line and function
//   - HideCaller: hides the caller information
func (opts *Logger) Caller(level ShowCallerLevel) {
	opts.showCaller = level
}

// Timestamp sets the level of timestamp information to show
// in the logs based on the level parameter
// the level can be one of the following:
//   - ShowFullTimestamp: shows the full timestamp with date and time
//   - ShowDateTime: shows the timestamp with date and time
//   - ShowTime: shows the timestamp with time only
//   - HideTimestamp: hides the timestamp
func (opts *Logger) Timestamp(level ShowTimestampLevel) {
	opts.showTimestamp = level
}

// ShowTags sets the logger to show the tags in the logs
// if the show parameter is true, otherwise it will hide the tags
func (opts *Logger) ShowTags(show bool) {
	opts.showTags = show
}

// Tags adds the tags to the logger
// the tags will be added to the logs created with this logger
func (opts *Logger) Tags(tags ...string) {
	opts.tags = append(opts.tags, tags...)
}

// SetTags sets the tags to the logger
// this method replaces the current tags with the new ones
func (opts *Logger) SetTags(tags ...string) {
	opts.tags = append(make([]string, 0), tags...)
}

// SetFatal sets the title and message to show in the fatal error
// alert when the Fatal method is called
func (opts *Logger) SetFatal(title, message string) {
	opts.fatalTitle = title
	opts.fatalMessage = message
}

// Debug creates a debug log message in the database
// with the message and arguments passed
// it formats the message with the arguments using fmt.Sprintf
// The new log is created in the database, but it is not printed
// if it fails to create the log it will return an error
func (opts *Logger) Debug(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	log, err := newLog(Debug, opts.tags, formattedMessage)
	if err != nil {
		return err
	}
	return createNewLog(opts, log)
}

// Info creates an info log message in the database
// with the message and arguments passed
// it formats the message with the arguments using fmt.Sprintf
// The new log is created in the database, but it is not printed
// if it fails to create the log it will return an error
func (opts *Logger) Info(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	log, err := newLog(Info, opts.tags, formattedMessage)
	if err != nil {
		return err
	}
	return createNewLog(opts, log)
}

// Warn creates a warning log message in the database
// with the message and arguments passed
// it formats the message with the arguments using fmt.Sprintf
// The new log is created in the database, but it is not printed
// if it fails to create the log it will return an error
func (opts *Logger) Warn(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	log, err := newLog(Warning, opts.tags, formattedMessage)
	if err != nil {
		return err
	}
	return createNewLog(opts, log)
}

// Error creates an error log message in the database
// with the message and arguments passed
// it formats the message with the arguments using fmt.Sprintf
// The new log is created in the database, but it is not printed
// if it fails to create the log it will return an error
func (opts *Logger) Error(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	log, err := newLog(Error, opts.tags, formattedMessage)
	if err != nil {
		return err
	}
	return createNewLog(opts, log)
}

// Fatal creates a fatal log message in the database only if the error passed is not nil
// it uses the error message as the message of the log
// The new log is created in the database, but it is not printed
// it will show an alert with the title and message set with SetFatal
// this method will exit the program with code 1
// if it fails to create the log it will return an error
func (opts *Logger) Fatal(e error) error {
	if e == nil {
		return nil
	}

	log, err := newLog(Fatal, opts.tags, e.Error())
	if err != nil {
		return err
	}

	err = createNewLog(opts, log)
	if err != nil {
		return err
	}

	beeep.Alert(opts.fatalTitle, opts.fatalMessage, "")
	os.Exit(1)
	return nil
}

// PrintDebug prints a debug log message in the console
// with the message and arguments passed
// it formats the message with the arguments using fmt.Sprintf
// The new log is not created in the database
// if it fails to print the log it will return an error
func (opts *Logger) PrintDebug(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	l, err := newLog(Debug, opts.tags, formattedMessage)
	if err != nil {
		return err
	}
	printLogs(opts, []*log{l})
	return nil
}

// PrintInfo prints an info log message in the console
// with the message and arguments passed
// it formats the message with the arguments using fmt.Sprintf
// The new log is not created in the database
// if it fails to print the log it will return an error
func (opts *Logger) PrintInfo(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	l, err := newLog(Info, opts.tags, formattedMessage)
	if err != nil {
		return err
	}
	printLogs(opts, []*log{l})
	return nil
}

// PrintWarn prints a warning log message in the console
// with the message and arguments passed
// it formats the message with the arguments using fmt.Sprintf
// The new log is not created in the database
// if it fails to print the log it will return an error
func (opts *Logger) PrintWarn(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	l, err := newLog(Warning, opts.tags, formattedMessage)
	if err != nil {
		return err
	}
	printLogs(opts, []*log{l})
	return nil
}

// PrintError prints an error log message in the console
// with the message and arguments passed
// it formats the message with the arguments using fmt.Sprintf
// The new log is not created in the database
// if it fails to print the log it will return an error
func (opts *Logger) PrintError(message string, args ...any) error {
	formattedMessage := fmt.Sprintf(message, args...)
	l, err := newLog(Error, opts.tags, formattedMessage)
	if err != nil {
		return err
	}
	printLogs(opts, []*log{l})
	return nil
}

// PrintFatal prints a fatal log message in the console and exits the program
// with the message and arguments passed only if the error passed is not nil
// it formats the message with the arguments using fmt.Sprintf
// The new log is not created in the database
// if it fails to print the log it will return an error
func (opts *Logger) PrintFatal(e error) error {
	if e == nil {
		return nil
	}

	l, err := newLog(Fatal, opts.tags, e.Error())
	if err != nil {
		return err
	}

	printLogs(opts, []*log{l})
	os.Exit(1)
	return nil
}

// PrintLogs prints the logs in the database based on the query options passed
// if it fails to query the logs it will return an error
func (opts *Logger) PrintLogs(queryOptions ...QueryOption) error {
	logs, err := queryLogs(opts, queryOptions...)
	if err != nil {
		return err
	}

	printLogs(opts, logs)
	return nil
}

// Export exports the logs in the database based on the query options passed
// to the export type passed
// the export type defines the format of the exported logs
// the export type can be one of the following:
//   - LOG: exports the logs in a .log file
//   - JSON: exports the logs in a .json file
//   - CSV: exports the logs in a .csv file
//
// the target folder for the exported file will be the folder path set in the logger
//
// this method returns the path of the exported file and an error if it fails to export the logs
func (opts *Logger) Export(exportType ExportType, queryOptions ...QueryOption) (string, error) {
	logs, err := queryLogs(opts, queryOptions...)
	if err != nil {
		return "", err
	}

	switch exportType {
	case JSON:
		return exportJson(logs, opts.folderPath)
	case CSV:
		return exportCSV(logs, opts.folderPath)
	default: // LOG
		return exportLogFile(logs, opts.folderPath)
	}
}

func createExportFile(filePath string) (*os.File, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		err := os.Remove(filePath)
		if err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func exportJson(logs []*log, folder string) (string, error) {
	filePath := filepath.Join(folder, fmt.Sprintf("%s_logs.json", time.Now().Format("20060102150405")))
	file, err := createExportFile(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	if len(logs) == 0 {
		_, err = file.WriteString("[]")
		if err != nil {
			return "", err
		}
		return filePath, nil
	}

	_, err = file.WriteString("[\n")
	if err != nil {
		return "", err
	}

	for i, log := range logs {
		if i > 0 {
			_, err = file.WriteString(",\n")
			if err != nil {
				return "", err
			}
		}

		_, err = file.WriteString(log.toJSON())
		if err != nil {
			return "", err
		}
	}

	_, err = file.WriteString("\n]")
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func exportCSV(logs []*log, folder string) (string, error) {
	filePath := filepath.Join(folder, fmt.Sprintf("%s_logs.csv", time.Now().Format("20060102150405")))
	file, err := createExportFile(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"level", "tags", "timestamp", "caller_file", "caller_line", "caller_function", "message"})
	if err != nil {
		return "", err
	}

	for _, log := range logs {
		err = writer.Write([]string{
			log.level.String(),
			strings.Join(log.tags, "|"),
			log.timestamp.String(),
			log.callerFile,
			fmt.Sprintf("%d", log.callerLine),
			log.callerFunction,
			log.message,
		})
		if err != nil {
			return "", err
		}
	}
	return filePath, nil
}

func exportLogFile(logs []*log, folder string) (string, error) {
	filePath := filepath.Join(folder, fmt.Sprintf("%s_logs.log", time.Now().Format("20060102150405")))
	file, err := createExportFile(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	for i, log := range logs {
		if i > 0 {
			_, err = file.WriteString("\n")
			if err != nil {
				return "", err
			}
		}

		_, err := file.WriteString(log.String())
		if err != nil {
			return "", err
		}
	}
	return filePath, nil
}
