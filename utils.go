package logs

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

func getCaller() (*caller, error) {
	c := new(caller)
	c.file = "unknown"
	c.line = 0
	c.funcion = "unknown"

	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return c, errors.New("[logger-pkg] failed to get the caller information")
	}

	c.file = filepath.Base(file)
	c.line = line

	f := runtime.FuncForPC(pc)
	c.funcion = f.Name()
	return c, nil
}

func addRigthPadding(s string, width int) string {
	l := runewidth.StringWidth(s)
	if l >= width {
		return s
	}
	return s + strings.Repeat(" ", width-l)
}

func printLogs(logs []*Log) {
	statusWidth := 0
	contextWidth := 0
	callerWidth := 0
	timeWidth := 0
	weekdayWidth := 0

	for _, log := range logs {
		l := runewidth.StringWidth(log.Status.String())
		if l > statusWidth {
			statusWidth = l
		}

		l = runewidth.StringWidth(fmt.Sprintf("[%s]", log.Context))
		if l > contextWidth {
			contextWidth = l
		}

		l = runewidth.StringWidth(fmt.Sprintf("[%s:%d - function: %s]", log.CallerFile, log.CallerLine, log.CallerFunction))
		if l > callerWidth {
			callerWidth = l
		}

		l = runewidth.StringWidth(fmt.Sprintf("[%s]", log.Timestamp))
		if l > timeWidth {
			timeWidth = l
		}

		t, _ := time.Parse("2006-01-02 15:04:05", log.Timestamp)
		wd := t.Weekday().String()
		l = runewidth.StringWidth(fmt.Sprintf("[%s]", wd))
		if l > weekdayWidth {
			weekdayWidth = l
		}
	}

	for _, log := range logs {
		var sb strings.Builder
		if timeWidth > 2 {
			sb.WriteString(addRigthPadding(fmt.Sprintf("[%s]", log.Timestamp), timeWidth))
		}

		if weekdayWidth > 2 {
			t, _ := time.Parse("2006-01-02 15:04:05", log.Timestamp)
			wd := t.Weekday().String()
			sb.WriteString(addRigthPadding(fmt.Sprintf("[%s]", wd), weekdayWidth))
			sb.WriteString(" ")
		}

		if contextWidth > 2 {
			sb.WriteString(addRigthPadding(fmt.Sprintf("[%s]", log.Context), contextWidth))
			sb.WriteString(" ")
		}

		if statusWidth > 0 {
			sb.WriteString(addRigthPadding(log.Status.String(), statusWidth))
			sb.WriteString(" ")
		}

		if callerWidth > 9 {
			sb.WriteString(addRigthPadding(fmt.Sprintf("[%s:%d - function: %s]", log.CallerFile, log.CallerLine, log.CallerFunction), callerWidth))
			sb.WriteString(" ")
		}

		sb.WriteString(log.Message)
		fmt.Println(sb.String())
	}
}
