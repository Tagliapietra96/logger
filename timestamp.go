package logger

import (
	"time"

	"github.com/Tagliapietra96/tui"
	"github.com/Tagliapietra96/tui/opts"
)

func newTimestamp(s string) timestamp {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return timestamp(t)
}

type timestamp time.Time

func (t timestamp) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}

func (t timestamp) toString(level ShowTimestampLevel) string {
	var layout string
	switch level {
	case ShowDate:
		layout = "2006-01-02"
	case ShowDateTime:
		layout = "2006-01-02 15:04:05"
	case ShowFullTimestamp:
		layout = "Monday 2006-01-02 15:04:05"
	default:
		return ""
	}
	return tui.Render(time.Time(t).Format(layout), opts.Muted)
}

// ShowTimestampLevel is an enum to define the level of timestamp information to be shown
// the level can be:
//   - HideTimestamp: hide the timestamp information
//   - ShowDate: show the date 2006-01-02
//   - ShowDateTime: show the date and time 2006-01-02 15:04:05
//   - ShowFullTimestamp: show the full timestamp Monday 2006-01-02 15:04:05
type ShowTimestampLevel int

const (
	HideTimestamp     ShowTimestampLevel = iota // hide the timestamp information
	ShowDate                                    // show the date 2006-01-02
	ShowDateTime                                // show the date and time 2006-01-02 15:04:05
	ShowFullTimestamp                           // show the full timestamp Monday 2006-01-02 15:04:05
)
