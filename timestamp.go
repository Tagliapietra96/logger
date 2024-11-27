package logger

import "time"

func newTimestamp(s string) timestamp {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return timestamp(t)
}

type timestamp time.Time

func (t timestamp) String() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}

// ShowTimestampLevel is an enum to define the level of timestamp information to be shown
type ShowTimestampLevel int

const (
	HideTimestamp     ShowTimestampLevel = iota // hide the timestamp information
	ShowDate                                    // show the date 2006-01-02
	ShowDateTime                                // show the date and time 2006-01-02 15:04:05
	ShowFullTimestamp                           // show the full timestamp Monday 2006-01-02 15:04:05
)
