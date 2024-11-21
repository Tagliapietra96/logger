package logs

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const table = `
CREATE TABLE IF NOT EXISTS logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	status INTEGER NOT NULL DEFAULT 0,
	context TEXT DEFAULT '',
	caller_file TEXT DEFAULT '',
	caller_line INTEGER DEFAULT 0,
	caller_function TEXT DEFAULT '',
	message TEXT DEFAULT '',
	time TEXT NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE INDEX IF NOT EXISTS id_index ON logs (id);
CREATE INDEX IF NOT EXISTS status_index ON logs (status);
CREATE INDEX IF NOT EXISTS context_index ON logs (context);
CREATE INDEX IF NOT EXISTS caller_file_index ON logs (caller_file);
CREATE INDEX IF NOT EXISTS caller_line_index ON logs (caller_line);
CREATE INDEX IF NOT EXISTS caller_function_index ON logs (caller_function);
CREATE INDEX IF NOT EXISTS message_index ON logs (message);
CREATE INDEX IF NOT EXISTS time_index ON logs (time);
`

const defaultQuery = "SELECT status, context, caller_file, caller_line, caller_function, message, time FROM logs"

func getDBConnection(useBinaryFolder bool) (*sql.DB, error) {
	var contextFolder, contextLabel string
	var db *sql.DB
	var err error

	if useBinaryFolder {
		contextFolder, err = os.Executable()
		contextLabel = "binary"
	} else {
		contextFolder, err = os.Getwd()
		contextLabel = "working"
	}

	if err != nil {
		return nil, errors.New("[logger-pkg] failed to get the current " + contextLabel + " directory: " + err.Error())
	}

	dbFilePath := filepath.Join(contextFolder, "logs_data.db")
	_, err = os.Stat(dbFilePath)

	if os.IsNotExist(err) {
		var dbFile *os.File
		dbFile, err = os.Create(dbFilePath)
		if err != nil {
			return nil, errors.New("[logger-pkg] failed to create the logs database file: " + err.Error())
		}
		dbFile.Close()
	} else if err != nil {
		return nil, errors.New("[logger-pkg] failed to check the logs database file: " + err.Error())
	}

	db, err = sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, errors.New("[logger-pkg] failed to open the logs database: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New("[logger-pkg] failed to get a connection to the logs database: " + err.Error())
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, errors.New("[logger-pkg] failed to generate the logs table: " + err.Error())
	}

	_, err = tx.Exec(table)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("[logger-pkg] failed to generate the logs table: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, errors.New("[logger-pkg] failed to generate the logs table: " + err.Error())
	}

	return db, nil
}

func createNewLog(opts *Options, status LogStatus, caller *caller, message string) error {
	db, err := getDBConnection(opts.GetUseBinaryFolder())
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
	}

	stmt, err := tx.Prepare("INSERT INTO logs (status, context, caller_file, caller_line, caller_function, message) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
	}

	_, err = stmt.Exec(int(status), opts.GetContext(), caller.file, caller.line, caller.funcion, message)
	if err != nil {
		tx.Rollback()
		return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
	}

	return nil
}

func queryLogs(configs ...QueryConfiguration) ([]*Log, error) {
	db, err := getDBConnection(false)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var query *strings.Builder
	query.WriteString(defaultQuery)
	for _, config := range configs {
		config(query)
	}
	query.WriteString(";")

	rows, err := db.Query(query.String())
	if err != nil {
		return nil, errors.New("[logger-pkg] failed to query the logs: " + err.Error())
	}
	defer rows.Close()

	var logs []*Log
	for rows.Next() {
		var status int
		var context, callerFile, callerFunction, message, time string
		var callerLine int

		err = rows.Scan(&status, &context, &callerFile, &callerLine, &callerFunction, &message, &time)
		if err != nil {
			return nil, errors.New("[logger-pkg] failed to scan the logs: " + err.Error())
		}

		logs = append(logs, &Log{
			Status:         LogStatus(status),
			Context:        context,
			CallerFile:     callerFile,
			CallerLine:     callerLine,
			CallerFunction: callerFunction,
			Message:        message,
			Timestamp:      time,
		})
	}

	return logs, nil
}
