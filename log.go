package logger

// log represents the log structure
type log struct {
	level          LogLevel
	tags           []string
	callerFile     string
	callerLine     int
	callerFunction string
	message        string
	timestamp      timestamp
}

// LogLevel represents the level of the log
//
//   - Debug: used for debugging purposes
//   - Info: used for informational messages
//   - Warning: used for warning messages
//   - Error: used for error messages
//   - Fatal: used for fatal messages
type LogLevel int

const (
	Debug   LogLevel = iota // debug level
	Info                    // info level
	Warning                 // warning level
	Error                   // error level
	Fatal                   // fatal level
)

// String returns the string representation of the LogLevel
// it returns the label of the level in uppercase
func (ls LogLevel) String() string {
	var s string
	switch ls {
	case Debug:
		s = "DEBUG"
	case Info:
		s = "INFO"
	case Warning:
		s = "WARNING"
	case Error:
		s = "ERROR"
	case Fatal:
		s = "FATAL"
	default:
		s = ""
	}

	return s
}

// LogField represents the fields of the log
// it is used in the query configuration to specify the fields to be returned
type LogField int

const (
	Level          LogField = iota // the level of the log
	Tags                           // the tags of the log
	CallerFile                     // the file of the caller
	Callerline                     // the line of the caller
	CallerFunction                 // the function of the caller
	Message                        // the message of the log
	Timestamp                      // the timestamp of the log
)

// String returns the string representation of the LogField
// it returns the label of the field in lowercase used in the database
func (lf LogField) String() string {
	var label string

	switch lf {
	case Level:
		label = "level"
	case Tags:
		label = "tags"
	case CallerFile:
		label = "caller_file"
	case Callerline:
		label = "caller_line"
	case CallerFunction:
		label = "caller_function"
	case Message:
		label = "message"
	case Timestamp:
		label = "time"
	default:
		label = ""
	}

	return label
}
