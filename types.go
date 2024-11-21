package logs

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	fatalColor   = lipgloss.AdaptiveColor{Light: "201", Dark: "213"}
	errorColor   = lipgloss.AdaptiveColor{Light: "160", Dark: "196"}
	warningColor = lipgloss.AdaptiveColor{Light: "208", Dark: "214"}
	infoColor    = lipgloss.AdaptiveColor{Light: "33", Dark: "45"}
	debugColor   = lipgloss.AdaptiveColor{Light: "27", Dark: "33"}
	unknownColor = lipgloss.AdaptiveColor{Light: "244", Dark: "241"}
)

type QueryOperator interface {
	op() string
}

type NumericOperator int

const (
	EQUAL NumericOperator = iota
	NOT_EQUAL
	GREATER_THAN
	GREATER_THAN_OR_EQUAL
	LESS_THAN
	LESS_THAN_OR_EQUAL
)

func (no NumericOperator) op() string {
	var operator string

	switch no {
	case EQUAL:
		operator = "="
	case NOT_EQUAL:
		operator = "!="
	case GREATER_THAN:
		operator = ">"
	case GREATER_THAN_OR_EQUAL:
		operator = ">="
	case LESS_THAN:
		operator = "<"
	case LESS_THAN_OR_EQUAL:
		operator = "<="
	default:
		operator = ""
	}

	return operator
}

type StringOperator int

const (
	CONTAINS StringOperator = iota
	NOT_CONTAINS
	SAME
)

func (so StringOperator) op() string {
	var operator string

	switch so {
	case CONTAINS:
		operator = "LIKE"
	case NOT_CONTAINS:
		operator = "NOT LIKE"
	case SAME:
		operator = "="
	default:
		operator = ""
	}
	return operator
}

type SortOperator int

const (
	ASC SortOperator = iota
	DESC
)

func (so SortOperator) op() string {
	var operator string

	switch so {
	case ASC:
		operator = "ASC"
	case DESC:
		operator = "DESC"
	default:
		operator = ""
	}
	return operator
}

type LogField int

const (
	STATUS LogField = iota
	CONTEXT
	CALLER_FILE
	CALLER_LINE
	CALLER_FUNCTION
	MESSAGE
	TIMESTAMP
)

func (lf LogField) String() string {
	var label string

	switch lf {
	case STATUS:
		label = "status"
	case CONTEXT:
		label = "context"
	case CALLER_FILE:
		label = "caller_file"
	case CALLER_LINE:
		label = "caller_line"
	case CALLER_FUNCTION:
		label = "caller_function"
	case MESSAGE:
		label = "message"
	case TIMESTAMP:
		label = "time"
	default:
		label = ""
	}

	return label
}

// LogStatus represents the status of the log
//
//   - DEBUG: used for debugging purposes
//   - INFO: used for informational messages
//   - WARNING: used for warning messages
//   - ERROR: used for error messages
//   - FATAL: used for fatal messages
type LogStatus int

const (
	DEBUG LogStatus = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// String returns the string representation of the LogStatus
// it returns the label of the status in uppercase, bold and colored
func (ls LogStatus) String() string {
	var label string
	var color lipgloss.TerminalColor

	switch ls {
	case DEBUG:
		label = "DEBUG"
		color = debugColor
	case INFO:
		label = "INFO"
		color = infoColor
	case WARNING:
		label = "WARNING"
		color = warningColor
	case ERROR:
		label = "ERROR"
		color = errorColor
	case FATAL:
		label = "FATAL"
		color = fatalColor
	default:
		label = "UNKNOWN"
		color = unknownColor
	}

	return lipgloss.NewStyle().Foreground(color).Inline(true).Bold(true).Render(label)
}

type caller struct {
	file    string
	line    int
	funcion string
}

type Log struct {
	Status         LogStatus
	Context        string
	CallerFile     string
	CallerLine     int
	CallerFunction string
	Message        string
	Timestamp      string
}

// Options represents the options for the logger
// it is used to configure the logger functions
// it has the following fields:
//
//   - useBinaryFolder: a boolean that indicates if the logs database should be stored in the binary folder, if true it will be stored in the binary folder, if false it will be stored in the working directory
//   - context: a string that represents the context of the log, it is used to identify the logs (e.g. the name of the application, the name of the service, etc.), by default it is empty. It accepts multiple contexts separated by commas
//   - defaultFatalTitle: a string that represents the default title for the fatal notifications, by default it is 'FATAL'
//   - defaultFatalMessage: a string that represents the default message for the fatal notifications, by default it is 'An error occurred, please check the logs for more information'
//
// Please note that the fatal notifications are not the messages printed to the terminal, they are the messages visible via real-time os notifications
type Options struct {
	useBinaryFolder     bool
	context             string
	defaultFatalTitle   string
	defaultFatalMessage string
}

// UseBinaryFolder sets the useBinaryFolder field of the Options
func (o *Options) UseBinaryFolder(use bool) {
	o.useBinaryFolder = use
}

// Context sets the context field of the Options
func (o *Options) Context(context string) {
	o.context = context
}

// AddContext adds a context to the context field of the Options
// it appends the context to the existing context separated by a comma
func (o *Options) AddContext(context string) {
	if o.context == "" {
		o.context = context
	} else {
		o.context += ", " + context
	}
}

// DefaultFatalTitle sets the defaultFatalTitle field of the Options
func (o *Options) DefaultFatalTitle(title string) {
	o.defaultFatalTitle = title
}

// DefaultFatalMessage sets the defaultFatalMessage field of the Options
func (o *Options) DefaultFatalMessage(message string) {
	o.defaultFatalMessage = message
}

// GetUseBinaryFolder returns the useBinaryFolder field of the Options
func (o *Options) GetUseBinaryFolder() bool {
	return o.useBinaryFolder
}

// GetContext returns the context field of the Options
func (o *Options) GetContext() string {
	return o.context
}

// GetDefaultFatalTitle returns the defaultFatalTitle field of the Options
func (o *Options) GetDefaultFatalTitle() string {
	return o.defaultFatalTitle
}

// GetDefaultFatalMessage returns the defaultFatalMessage field of the Options
func (o *Options) GetDefaultFatalMessage() string {
	return o.defaultFatalMessage
}

type OptionConfiguration func(*Options)

func WithBinFolder(use bool) OptionConfiguration {
	return func(o *Options) {
		o.UseBinaryFolder(use)
	}
}

func WithContext(context string) OptionConfiguration {
	return func(o *Options) {
		o.AddContext(context)
	}
}

func WithFatalTitle(title string) OptionConfiguration {
	return func(o *Options) {
		o.DefaultFatalTitle(title)
	}
}

func WithFatalMessage(message string) OptionConfiguration {
	return func(o *Options) {
		o.DefaultFatalMessage(message)
	}
}

type QueryConfiguration func(*strings.Builder)

func CustomQuery(query string) QueryConfiguration {
	return func(sb *strings.Builder) {
		sb.WriteString(" ")
		sb.WriteString(query)
	}
}

func prepareFilter(config QueryConfiguration) QueryConfiguration {
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

func prepareSort(config QueryConfiguration) QueryConfiguration {
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

		config(sb)

		if limit != "" {
			sb.WriteString(" LIMIT ")
			sb.WriteString(limit)
		}
	}
}

func AddFilters(configs ...QueryConfiguration) QueryConfiguration {
	return func(sb *strings.Builder) {
		for _, config := range configs {
			prepareFilter(config)
		}
	}
}

func AddSorts(configs ...QueryConfiguration) QueryConfiguration {
	return func(sb *strings.Builder) {
		for _, config := range configs {
			prepareSort(config)
		}
	}
}

func AddLimit(limitAndOffset ...int) QueryConfiguration {
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

func FilterByStatus(status LogStatus, operator NumericOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s %d", STATUS.String(), operator.op(), status))
	})
}

func FilterByContext(context string, operator StringOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			context = fmt.Sprintf("%%%s%%", context)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", CONTEXT.String(), operator.op(), context))
	})
}

func FilterByCallerFile(file string, operator StringOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			file = fmt.Sprintf("%%%s%%", file)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", CALLER_FILE.String(), operator.op(), file))
	})
}

func FilterByCallerLine(line int, operator NumericOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s %d", CALLER_LINE.String(), operator.op(), line))
	})
}

func FilterByCallerFunction(function string, operator StringOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			function = fmt.Sprintf("%%%s%%", function)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", CALLER_FUNCTION.String(), operator.op(), function))
	})
}

func FilterByMessage(message string, operator StringOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			message = fmt.Sprintf("%%%s%%", message)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", MESSAGE.String(), operator.op(), message))
	})
}

func FilterByTimestamp(timestamp string, operator QueryOperator) QueryConfiguration {
	return prepareFilter(func(sb *strings.Builder) {
		if operator == CONTAINS || operator == NOT_CONTAINS {
			timestamp = fmt.Sprintf("%%%s%%", timestamp)
		}
		sb.WriteString(fmt.Sprintf("%s %s '%s'", TIMESTAMP.String(), operator.op(), timestamp))
	})
}

func SortByStatus(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", STATUS.String(), order.op()))
	})
}

func SortByContext(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", CONTEXT.String(), order.op()))
	})
}

func SortByCallerFile(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", CALLER_FILE.String(), order.op()))
	})
}

func SortByCallerLine(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", CALLER_LINE.String(), order.op()))
	})
}

func SortByCallerFunction(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", CALLER_FUNCTION.String(), order.op()))
	})
}

func SortByMessage(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", MESSAGE.String(), order.op()))
	})
}

func SortByTimestamp(order SortOperator) QueryConfiguration {
	return prepareSort(func(sb *strings.Builder) {
		sb.WriteString(fmt.Sprintf("%s %s", TIMESTAMP.String(), order.op()))
	})
}
