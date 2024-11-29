package logger

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const table = `
CREATE TABLE IF NOT EXISTS logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	level INTEGER NOT NULL DEFAULT 0,
	caller_file TEXT DEFAULT '',
	caller_line INTEGER DEFAULT 0,
	caller_function TEXT DEFAULT '',
	message TEXT DEFAULT '',
	time TEXT NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE INDEX IF NOT EXISTS logs_id_index ON logs (id);
CREATE INDEX IF NOT EXISTS logs_level_index ON logs (level);
CREATE INDEX IF NOT EXISTS logs_caller_file_index ON logs (caller_file);
CREATE INDEX IF NOT EXISTS logs_caller_line_index ON logs (caller_line);
CREATE INDEX IF NOT EXISTS logs_caller_function_index ON logs (caller_function);
CREATE INDEX IF NOT EXISTS logs_message_index ON logs (message);
CREATE INDEX IF NOT EXISTS logs_time_index ON logs (time);

CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE INDEX IF NOT EXISTS tags_id_index ON tags (id);
CREATE INDEX IF NOT EXISTS tags_name_index ON tags (name);

CREATE TABLE IF NOT EXISTS log_tags (
    log_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (log_id, tag_id),
    FOREIGN KEY (log_id) REFERENCES logs(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS lt_log_id_index ON log_tags (log_id);
CREATE INDEX IF NOT EXISTS lt_tag_id_index ON log_tags (tag_id);
`

const defaultQuery = `
SELECT DISTINCT logs.id, logs.level, logs.caller_file, logs.caller_line, logs.caller_function, logs.message, logs.time
FROM logs
INNER JOIN log_tags ON logs.id = log_tags.log_id
INNER JOIN tags ON log_tags.tag_id = tags.id
`

type QueryOption func(*strings.Builder)

func getDBConnection(folderPath string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	dbFilePath := filepath.Join(folderPath, "logs_data.db")
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

func createNewLog(opts *Logger, log *log) error {
	db, err := getDBConnection(opts.folderPath)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
	}

	logstmt, err := tx.Prepare("INSERT INTO logs (level, caller_file, caller_line, caller_function, message) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
	}
	defer logstmt.Close()

	result, err := logstmt.Exec(int(log.level), log.callerFile, log.callerLine, log.callerFunction, log.message)
	if err != nil {
		tx.Rollback()
		return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
	}

	logId, err := result.LastInsertId()
	if err != nil || logId < 1 {
		tx.Rollback()
		return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
	}

	if len(log.tags) > 0 {
		for _, tag := range log.tags {
			tagstmt, err := tx.Prepare("INSERT OR IGNORE INTO tags (name) VALUES (?);")
			if err != nil {
				return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
			}
			defer tagstmt.Close()

			_, err = tagstmt.Exec(tag)
			if err != nil {
				tx.Rollback()
				return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
			}

			linkstmt, err := tx.Prepare("INSERT INTO log_tags (log_id, tag_id) VALUES (?, (SELECT id FROM tags WHERE name = ?));")
			if err != nil {
				return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
			}
			defer linkstmt.Close()

			_, err = linkstmt.Exec(logId, tag)
			if err != nil {
				tx.Rollback()
				return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return errors.New("[logger-pkg] failed to create a new log: " + err.Error())
	}

	return nil
}

func queryLogs(opts *Logger, configs ...QueryOption) ([]*log, error) {
	db, err := getDBConnection(opts.folderPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	query := new(strings.Builder)
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

	var logs []*log
	for rows.Next() {
		var id, level, callerLine int
		var callerFile, callerFunction, message, time string

		err = rows.Scan(&id, &level, &callerFile, &callerLine, &callerFunction, &message, &time)
		if err != nil {
			return nil, errors.New("[logger-pkg] failed to scan the logs: " + err.Error())
		}

		tags, err := getTagsForLog(db, id)
		if err != nil {
			return nil, errors.New("[logger-pkg] failed to get the tags for the logs: " + err.Error())
		}

		logs = append(logs, &log{
			level:          LogLevel(level),
			tags:           tags,
			callerFile:     callerFile,
			callerLine:     callerLine,
			callerFunction: callerFunction,
			message:        message,
			timestamp:      newTimestamp(time),
		})
	}

	return logs, nil
}

func getTagsForLog(db *sql.DB, logId int) ([]string, error) {
	tags := make([]string, 0)
	rows, err := db.Query("SELECT tags.name FROM tags INNER JOIN log_tags ON tags.id = log_tags.tag_id WHERE log_tags.log_id = ?", logId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
