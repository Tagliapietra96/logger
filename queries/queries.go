package queries

import (
	"fmt"
	"strings"
	"time"

	"github.com/Tagliapietra96/logger"
)

const defaultQuery = `
SELECT DISTINCT logs.id, logs.level, logs.caller_file, logs.caller_line, logs.caller_function, logs.message, logs.time
FROM logs
INNER JOIN log_tags ON logs.id = log_tags.log_id
INNER JOIN tags ON log_tags.tag_id = tags.id
`

func getOrder(order string) string {
	order = strings.ToUpper(order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}
	return order
}

// CustomQuery returns a QueryOption that appends the given query to the base query
// This is useful for custom queries that are not covered by the other QueryOptions
// Example:
//
//	queryOpt := queries.CustomQuery("WHERE level = 1 OR level = 3 ORDER BY time DESC")
//
// In this example te custom query is very simple, but it can be as complex as needed,
// as long as it is a valid SQL query.
// Note the base query is not needed, as it is already defined in the package.
// The custom query is appended to the base query.
// The resulting query will be:
//
//	SELECT id, level, caller_file, caller_line, caller_function, message, time FROM logs WHERE level = 1 OR level = 3 ORDER BY time DESC
//
// The main approach for this package is to use the other QueryOptions, as they are more specific and easier to use.
func CustomQuery(query string) logger.QueryOption {
	return func(sb *strings.Builder) {
		sb.WriteString(" ")
		sb.WriteString(query)
	}
}

func prepareFilter(config logger.QueryOption) logger.QueryOption {
	return func(sb *strings.Builder) {
		var filter, order, limit string
		s := sb.String()
		if s == "" {
			sb.WriteString(defaultQuery)
		}

		if strings.Contains(s, " WHERE ") {
			pieces := strings.Split(s, " WHERE ")
			filter = pieces[1]
			if strings.Contains(filter, " ORDER BY ") {
				pieces = strings.Split(filter, " ORDER BY ")
				filter = pieces[0]
				order = pieces[1]
				if strings.Contains(order, " LIMIT ") {
					pieces = strings.Split(order, " LIMIT ")
					order = pieces[0]
					limit = pieces[1]
				}
			} else if strings.Contains(filter, " LIMIT ") {
				pieces = strings.Split(filter, " LIMIT ")
				filter = pieces[0]
				limit = pieces[1]
			}
			sb.Reset()
			sb.WriteString(defaultQuery)
			sb.WriteString(" WHERE ")
			sb.WriteString(filter)
		} else if strings.Contains(s, " ORDER BY ") {
			pieces := strings.Split(s, " ORDER BY ")
			order = pieces[1]
			if strings.Contains(order, " LIMIT ") {
				pieces = strings.Split(order, " LIMIT ")
				order = pieces[0]
				limit = pieces[1]
			}
			sb.Reset()
			sb.WriteString(defaultQuery)
		} else if strings.Contains(s, " LIMIT ") {
			pieces := strings.Split(s, " LIMIT ")
			limit = pieces[1]
			sb.Reset()
			sb.WriteString(defaultQuery)
		}

		if !strings.Contains(s, " WHERE ") {
			sb.WriteString(" WHERE ")
		} else {
			sb.WriteString(" AND ")
		}

		config(sb)

		if order != "" {
			sb.WriteString(" ORDER BY ")
			sb.WriteString(order)
		}

		if limit != "" {
			sb.WriteString(" LIMIT ")
			sb.WriteString(limit)
		}
	}
}

func prepareSort(config logger.QueryOption) logger.QueryOption {
	return func(sb *strings.Builder) {
		var base, order, limit string
		s := sb.String()
		if s == "" {
			sb.WriteString(defaultQuery)
		}

		if strings.Contains(s, " ORDER BY ") {
			pieces := strings.Split(s, " ORDER BY ")
			base = pieces[0]
			order = pieces[1]
			if strings.Contains(order, " LIMIT ") {
				pieces = strings.Split(order, " LIMIT ")
				order = pieces[0]
				limit = pieces[1]
			}
			sb.Reset()
			sb.WriteString(base)
			sb.WriteString(" ORDER BY ")
			sb.WriteString(order)
		} else if strings.Contains(s, " LIMIT ") {
			pieces := strings.Split(s, " LIMIT ")
			base = pieces[0]
			limit = pieces[1]
			sb.Reset()
			sb.WriteString(base)
		}

		if !strings.Contains(s, " ORDER BY ") {
			sb.WriteString(" ORDER BY ")
		} else {
			sb.WriteString(", ")
		}

		config(sb)

		if limit != "" {
			sb.WriteString(" LIMIT ")
			sb.WriteString(limit)
		}
	}
}

// AddFilters returns a QueryOption that appends the given filters to the base query
// This is useful to add multiple filters to the query
// Example:
//
//	queryOpt := queries.AddFilters(queries.LevelEqual(logger.Info), queries.CallerFileLike("main"))
//
// Note: this metod is useful if you want to use custom QueryOprions to filter the logs
// If you use the QueryOptions provided by this package, you can avoid to use this method
// because every QueryOption already has the logic to add the filter to the query without
// conflicting with other QueryOptions.
func AddFilters(configs ...logger.QueryOption) logger.QueryOption {
	return func(sb *strings.Builder) {
		for _, config := range configs {
			prepareFilter(config)
		}
	}
}

// AddSorts returns a QueryOption that appends the given sorts to the base query
// This is useful to add multiple sorts to the query
// Example:
//
//	queryOpt := queries.AddSorts(logger.SortLevel("DESC"), queries.SortTimestamp("ASC"))
//
// Note: this metod is useful if you want to use custom QueryOprions to sort the logs
// If you use the QueryOptions provided by this package, you can avoid to use this method
// because every QueryOption already has the logic to add the sort to the query without
// conflicting with other QueryOptions.
func AddSorts(configs ...logger.QueryOption) logger.QueryOption {
	return func(sb *strings.Builder) {
		for _, config := range configs {
			prepareSort(config)
		}
	}
}

// AddLimit returns a QueryOption that appends the given limit and offset to the base query
// This is useful to add the limit and offset to the query
//
//   - If only one argument is provided, it will be used as the limit
//   - If two arguments are provided, the first will be used as the limit and the second as the offset
//   - If no arguments are provided, the method does nothing
//
// Example:
//
//	queryOpt := queries.AddLimit(10)
//	queryOpt := queries.AddLimit(10, 5)
//
// In the first example, the query will have a limit of 10 logs and no offset
// In the second example, the query will have a limit of 10 logs and an offset of 5
func AddLimit(limitAndOffset ...int) logger.QueryOption {
	return func(sb *strings.Builder) {
		if len(limitAndOffset) == 0 {
			return
		}

		var base string
		s := sb.String()
		if s == "" {
			sb.WriteString(defaultQuery)
		}

		if strings.Contains(s, " LIMIT ") {
			pieces := strings.Split(s, " LIMIT ")
			base = pieces[0]
			sb.Reset()
			sb.WriteString(base)
		}

		sb.WriteString(" LIMIT ")
		sb.WriteString(fmt.Sprintf("%d", limitAndOffset[0]))

		if len(limitAndOffset) > 1 {
			sb.WriteString(" OFFSET ")
			sb.WriteString(fmt.Sprintf("%d", limitAndOffset[1]))
		}
	}
}

// HasTags returns a QueryOption that filters the logs by the given tags
// the logs must have at least one of the given tags
// Example:
//
//	queryOpt := queries.HasTags("tag1", "tag2")
//
// In this example, the query will return all the logs with the tags set to tag1 or tag2
// or any other tag with the string "tag1" or "tag2" in its name
// The query will return the logs with at least one of the given tags
func HasTags(tag string, tags ...string) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		for i, tag := range append(tags, tag) {
			sb.WriteString(fmt.Sprintf("tags.name LIKE '%%%s%%'", tag))
			if i < len(tags)-1 {
				sb.WriteString(" OR ")
			}
		}
	})
}

// LevelEqual returns a QueryOption that filters the logs by the given level
// Example:
//
//	queryOpt := queries.LevelEqual(logger.Info)
//
// In this example, the query will return all the logs with the level set to Info
func LevelEqual(level logger.LogLevel) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.level = %d", level))
	})
}

// LevelNotEqual returns a QueryOption that filters the logs by the levels different from the given level
// Example:
//
//	queryOpt := queries.LevelNotEqual(logger.Info)
//
// In this example, the query will return all the logs with the level different from Info
func LevelNotEqual(level logger.LogLevel) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.level != %d", level))
	})
}

// LevelGreaterThan returns a QueryOption that filters the logs by the levels greater than the given level
// Example:
//
//	queryOpt := queries.LevelGreaterThan(logger.Info) // warning, error, fatal
//
// In this example, the query will return all the logs with the level greater than Info
func LevelGreaterThan(level logger.LogLevel) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.level > %d", level))
	})
}

// LevelLessThan returns a QueryOption that filters the logs by the levels less than the given level
// Example:
//
//	queryOpt := queries.LevelLessThan(logger.Info) // debug
//
// In this example, the query will return all the logs with the level less than Info
func LevelLessThan(level logger.LogLevel) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.level < %d", level))
	})
}

// LevelBetween returns a QueryOption that filters the logs by the levels between the given start and end levels
// Example:
//
//	queryOpt := queries.LevelBetween(logger.Info, logger.Warning) // info, warning
//
// In this example, the query will return all the logs with the level between Info and Warning
func LevelBetween(start, end logger.LogLevel) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.level BETWEEN %d AND %d", start, end))
	})
}

// CallerFileLike returns a QueryOption that filters the logs by the given file
// Example:
//
//	queryOpt := queries.CallerFileLike("main.go")
//
// In this example, the query will return all the logs with the caller file set to main.go
// or any other file with the string "main.go" in its name
func CallerFileLike(file string) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_file LIKE '%%%s%%'", file))
	})
}

// CallerFileNotLike returns a QueryOption that filters the logs by the files different from the given file
// Example:
//
//	queryOpt := queries.CallerFileNotLike("main.go")
//
// In this example, the query will return all the logs with the caller file different from main.go
// or any other file without the string "main.go" in its name
func CallerFileNotLike(file string) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_file NOT LIKE '%%%s%%'", file))
	})
}

// CallerLineEqual returns a QueryOption that filters the logs by the given line
// Example:
//
//	queryOpt := queries.CallerLineEqual(10)
//
// In this example, the query will return all the logs with the caller line set to 10
func CallerLineEqual(line int) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_line = %d", line))
	})
}

// CallerLineNotEqual returns a QueryOption that filters the logs by the lines different from the given line
// Example:
//
//	queryOpt := queries.CallerLineNotEqual(10)
//
// In this example, the query will return all the logs with the caller line different from 10
func CallerLineNotEqual(line int) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_line != %d", line))
	})
}

// CallerLineGreaterThan returns a QueryOption that filters the logs by the lines greater than the given line
// Example:
//
//	queryOpt := queries.CallerLineGreaterThan(10)
//
// In this example, the query will return all the logs with the caller line greater than 10
func CallerLineGreaterThan(line int) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_line > %d", line))
	})
}

// CallerLineLessThan returns a QueryOption that filters the logs by the lines less than the given line
// Example:
//
//	queryOpt := queries.CallerLineLessThan(10)
//
// In this example, the query will return all the logs with the caller line less than 10
func CallerLineLessThan(line int) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_line < %d", line))
	})
}

// CallerLineBetween returns a QueryOption that filters the logs by the lines between the given start and end lines
// Example:
//
//	queryOpt := queries.CallerLineBetween(10, 20)
//
// In this example, the query will return all the logs with the caller line between 10 and 20
func CallerLineBetween(start, end int) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_line BETWEEN %d AND %d", start, end))
	})
}

// CallerFunctionLike returns a QueryOption that filters the logs by the given function
// Example:
//
//	queryOpt := queries.CallerFunctionLike("main.main")
//
// In this example, the query will return all the logs with the caller function set to main.main
// or any other function with the string "main.main" in its name
func CallerFunctionLike(function string) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_function LIKE '%%%s%%'", function))
	})
}

// CallerFunctionNotLike returns a QueryOption that filters the logs by the functions different from the given function
// Example:
//
//	queryOpt := queries.CallerFunctionNotLike("main.main")
//
// In this example, the query will return all the logs with the caller function different from main.main
// or any other function without the string "main.main" in its name
func CallerFunctionNotLike(function string) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_function NOT LIKE '%%%s%%'", function))
	})
}

// MessageLike returns a QueryOption that filters the logs by the given message
// Example:
//
//	queryOpt := queries.MessageLike("error")
//
// In this example, the query will return all the logs with the message set to error
// or any other message with the string "error" in its content
func MessageLike(message string) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.message LIKE '%%%s%%'", message))
	})
}

// MessageNotLike returns a QueryOption that filters the logs by the messages different from the given message
// Example:
//
//	queryOpt := queries.MessageNotLike("error")
//
// In this example, the query will return all the logs with the message different from error
// or any other message without the string "error" in its content
func MessageNotLike(message string) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.message NOT LIKE '%%%s%%'", message))
	})
}

// TimestampEqual returns a QueryOption that filters the logs by the given timestamp
// Example:
//
//	queryOpt := queries.TimestampEqual(time.Now())
//
// In this example, the query will return all the logs with the timestamp set to the current time
// it consider both date and time
func TimestampEqual(timestamp time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.time = '%s'", timestamp.Format("2006-01-02 15:04:05")))
	})
}

// TimestampNotEqual returns a QueryOption that filters the logs by the timestamps different from the given timestamp
// Example:
//
//	queryOpt := queries.TimestampNotEqual(time.Now())
//
// In this example, the query will return all the logs with the timestamp different from the current time
// it consider both date and time
func TimestampNotEqual(timestamp time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.time != '%s'", timestamp.Format("2006-01-02 15:04:05")))
	})
}

// TimestampGreaterThan returns a QueryOption that filters the logs by the timestamps greater than the given timestamp
// Example:
//
//	queryOpt := queries.TimestampGreaterThan(time.Now())
//
// In this example, the query will return all the logs with the timestamp greater than the current time
// it consider both date and time
func TimestampGreaterThan(timestamp time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.time > '%s'", timestamp.Format("2006-01-02 15:04:05")))
	})
}

// TimestampLessThan returns a QueryOption that filters the logs by the timestamps less than the given timestamp
// Example:
//
//	queryOpt := queries.TimestampLessThan(time.Now())
//
// In this example, the query will return all the logs with the timestamp less than the current time
// it consider both date and time
func TimestampLessThan(timestamp time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.time < '%s'", timestamp.Format("2006-01-02 15:04:05")))
	})
}

// TimestampBetween returns a QueryOption that filters the logs by the timestamps between the given start and end timestamps
// Example:
//
//	queryOpt := queries.TimestampBetween(time.Now().Add(-time.Hour), time.Now())
//
// In this example, the query will return all the logs with the timestamp between one hour ago and the current time
// it consider both date and time
func TimestampBetween(start, end time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.time BETWEEN '%s' AND '%s'", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05")))
	})
}

// DateEqual returns a QueryOption that filters the logs by the given date
// Example:
//
//	queryOpt := queries.DateEqual(time.Now())
//
// In this example, the query will return all the logs with the date set to the current date
// it consider only the date, not the time
func DateEqual(date time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("DATE(logs.time) = '%s'", date.Format("2006-01-02")))
	})
}

// DateNotEqual returns a QueryOption that filters the logs by the dates different from the given date
// Example:
//
//	queryOpt := queries.DateNotEqual(time.Now())
//
// In this example, the query will return all the logs with the date different from the current date
// it consider only the date, not the time
func DateNotEqual(date time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("DATE(logs.time) != '%s'", date.Format("2006-01-02")))
	})
}

// DateGreaterThan returns a QueryOption that filters the logs by the dates greater than the given date
// Example:
//
//	queryOpt := queries.DateGreaterThan(time.Now())
//
// In this example, the query will return all the logs with the date greater than the current date
// it consider only the date, not the time
func DateGreaterThan(date time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("DATE(logs.time) > '%s'", date.Format("2006-01-02")))
	})
}

// DateLessThan returns a QueryOption that filters the logs by the dates less than the given date
// Example:
//
//	queryOpt := queries.DateLessThan(time.Now())
//
// In this example, the query will return all the logs with the date less than the current date
// it consider only the date, not the time
func DateLessThan(date time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("DATE(logs.time) < '%s'", date.Format("2006-01-02")))
	})
}

// DateBetween returns a QueryOption that filters the logs by the dates between the given start and end dates
// Example:
//
//	queryOpt := queries.DateBetween(time.Now().Add(-24*time.Hour), time.Now())
//
// In this example, the query will return all the logs with the date between 24 hours ago and the current date
// it consider only the date, not the time
func DateBetween(start, end time.Time) logger.QueryOption {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("DATE(logs.time) BETWEEN '%s' AND '%s'", start.Format("2006-01-02"), end.Format("2006-01-02")))
	})
}

// SortLevel returns a QueryOption that sorts the logs by the level
// Example:
//
//	queryOpt := queries.SortLevel("DESC")
//
// In this example, the query will return the logs sorted by the level in descending order
// it accept only "ASC"/"asc" or "DESC"/"desc" as order. If the order is not valid, it will default to "ASC"
func SortLevel(order string) logger.QueryOption {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.level %s", getOrder(order)))
	})
}

// SortTags returns a QueryOption that sorts the logs by the tags
// Example:
//
//	queryOpt := queries.SortTags("ASC")
//
// In this example, the query will return the logs sorted by the tags in ascending order
// it accept only "ASC"/"asc" or "DESC"/"desc" as order. If the order is not valid, it will default to "ASC"
func SortCallerFile(order string) logger.QueryOption {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_file %s", getOrder(order)))
	})
}

// SortCallerLine returns a QueryOption that sorts the logs by the line of the caller
// Example:
//
//	queryOpt := queries.SortCallerLine("ASC")
//
// In this example, the query will return the logs sorted by the line of the caller in ascending order
// it accept only "ASC"/"asc" or "DESC"/"desc" as order. If the order is not valid, it will default to "ASC"
func SortCallerLine(order string) logger.QueryOption {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_line %s", getOrder(order)))
	})
}

// SortCallerFunction returns a QueryOption that sorts the logs by the function of the caller
// Example:
//
//	queryOpt := queries.SortCallerFunction("ASC")
//
// In this example, the query will return the logs sorted by the function of the caller in ascending order
// it accept only "ASC"/"asc" or "DESC"/"desc" as order. If the order is not valid, it will default to "ASC"
func SortCallerFunction(order string) logger.QueryOption {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.caller_function %s", getOrder(order)))
	})
}

// SortMessage returns a QueryOption that sorts the logs by the message
// Example:
//
//	queryOpt := queries.SortMessage("ASC")
//
// In this example, the query will return the logs sorted by the message in ascending order
// it accept only "ASC"/"asc" or "DESC"/"desc" as order. If the order is not valid, it will default to "ASC"
func SortMessage(order string) logger.QueryOption {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.message %s", getOrder(order)))
	})
}

// SortTimestamp returns a QueryOption that sorts the logs by the timestamp
// Example:
//
//	queryOpt := queries.SortTimestamp("ASC")
//
// In this example, the query will return the logs sorted by the timestamp in ascending order
// it accept only "ASC"/"asc" or "DESC"/"desc" as order. If the order is not valid, it will default to "ASC"
func SortTimestamp(order string) logger.QueryOption {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("logs.time %s", getOrder(order)))
	})
}
