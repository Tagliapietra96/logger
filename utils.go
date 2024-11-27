package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
)

func addRigthPadding(s string, width int) string {
	l := runewidth.StringWidth(s)
	if l >= width {
		return s
	}
	return s + strings.Repeat(" ", width-l)
}

func printLogs(logs []*log) {
	statusWidth := 0
	contextWidth := 0
	callerWidth := 0
	timeWidth := 0
	weekdayWidth := 0

	for _, log := range logs {
		l := runewidth.StringWidth(log.level.String())
		if l > statusWidth {
			statusWidth = l
		}

		l = runewidth.StringWidth(fmt.Sprintf("[%s]", log.tags))
		if l > contextWidth {
			contextWidth = l
		}

		l = runewidth.StringWidth(fmt.Sprintf("[%s:%d - function: %s]", log.callerFile, log.callerLine, log.callerFunction))
		if l > callerWidth {
			callerWidth = l
		}

		l = runewidth.StringWidth(fmt.Sprintf("[%s]", log.timestamp))
		if l > timeWidth {
			timeWidth = l
		}

		t, _ := time.Parse("2006-01-02 15:04:05", log.timestamp.String())
		wd := t.Weekday().String()
		l = runewidth.StringWidth(fmt.Sprintf("[%s]", wd))
		if l > weekdayWidth {
			weekdayWidth = l
		}
	}

	for _, log := range logs {
		var sb strings.Builder
		if timeWidth > 2 {
			sb.WriteString(addRigthPadding(fmt.Sprintf("[%s]", log.timestamp), timeWidth))
		}

		if weekdayWidth > 2 {
			t, _ := time.Parse("2006-01-02 15:04:05", log.timestamp.String())
			wd := t.Weekday().String()
			sb.WriteString(addRigthPadding(fmt.Sprintf("[%s]", wd), weekdayWidth))
			sb.WriteString(" ")
		}

		if contextWidth > 2 {
			sb.WriteString(addRigthPadding(fmt.Sprintf("[%s]", log.tags), contextWidth))
			sb.WriteString(" ")
		}

		if statusWidth > 0 {
			sb.WriteString(addRigthPadding(log.level.String(), statusWidth))
			sb.WriteString(" ")
		}

		if callerWidth > 9 {
			sb.WriteString(addRigthPadding(fmt.Sprintf("[%s:%d - function: %s]", log.callerFile, log.callerLine, log.callerFunction), callerWidth))
			sb.WriteString(" ")
		}

		sb.WriteString(log.message)
		fmt.Println(sb.String())
	}
}
