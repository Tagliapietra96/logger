package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/Tagliapietra96/tui"
	"github.com/Tagliapietra96/tui/opts"
	"github.com/charmbracelet/lipgloss"
)

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

func newLog(level LogLevel, tags []string, message string) (*log, error) {
	l := &log{
		level:     level,
		tags:      tags,
		message:   message,
		timestamp: timestamp(time.Now()),
	}

	err := getCaller(l)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (l *log) getTags() []string {
	result := make([]string, 0, len(l.tags))
	for _, tag := range l.tags {
		result = append(result, "🔖"+tag)
	}

	return result
}

func (l *log) getCaller(inline bool, level ShowCallerLevel) string {
	if level == HideCaller {
		return ""
	}

	c := tui.NewStyle(opts.Muted)
	if !inline {
		tui.Concat(&c, "at ")
	} else {
		tui.Concat(&c, "<")
	}

	if level >= ShowCallerFile {
		tui.Concat(&c, l.callerFile)
	}

	if level >= ShowCallerLine {
		tui.Concat(&c, fmt.Sprintf(":%d", l.callerLine))
	}

	if level >= ShowCallerFunction {
		tui.Concat(&c, " - ", l.callerFunction)
	}

	if inline {
		tui.Concat(&c, ">")
	}

	return c.String()
}

func (l *log) toJSON() string {
	var b strings.Builder
	b.WriteString("{\n")
	b.WriteString(fmt.Sprintf("\t\"level\": \"%s\",\n", l.level.String()))
	b.WriteString("\t\"tags\": [")
	for i, tag := range l.tags {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("\"%s\"", tag))
	}
	b.WriteString("],\n")
	b.WriteString(fmt.Sprintf("\t\"caller_file\": \"%s\",\n", l.callerFile))
	b.WriteString(fmt.Sprintf("\t\"caller_line\": %d,\n", l.callerLine))
	b.WriteString(fmt.Sprintf("\t\"caller_function\": \"%s\",\n", l.callerFunction))
	b.WriteString(fmt.Sprintf("\t\"message\": \"%s\",\n", l.message))
	b.WriteString(fmt.Sprintf("\t\"time\": \"%s\"\n", l.timestamp.String()))
	b.WriteString("}")
	return b.String()
}

func (l *log) String() string {
	return fmt.Sprintf(
		"%s [%s] <%s:%d - %s> %s: %s",
		l.timestamp.String(),
		strings.Join(l.tags, ", "),
		l.callerFile,
		l.callerLine,
		l.callerFunction,
		l.level.String(),
		l.message,
	)
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

func (ls LogLevel) color() lipgloss.TerminalColor {
	var color lipgloss.TerminalColor
	switch ls {
	case Debug:
		color = tui.ColorLink
	case Info:
		color = tui.ColorInfo
	case Warning:
		color = tui.ColorWarning
	case Error:
		color = tui.ColorError
	case Fatal:
		color = tui.ColorAccent
	default:
		color = tui.ColorMuted
	}

	return color
}

func (ls LogLevel) toString() string {
	s := ls.String()
	color := ls.color()
	return tui.Render(s, opts.Color(color))
}

// ExportType represents the type of the export
// it is used to specify the type of the export to be done
// the type can be:
//   - JSON: export the logs in JSON format
//   - CSV: export the logs in CSV format
//   - LOG: export the logs in LOG format
type ExportType int

const (
	JSON ExportType = iota // export the logs in JSON
	CSV                    // export the logs in CSV
	LOG                    // export the logs in LOG
)
