package logs

import (
	"testing"
)

func TestNewLogOpts(t *testing.T) {
	opts := NewLogOpts()
	if opts == nil {
		t.Error("NewLogOpts() should return a valid pointer to Options")
	}
}

func TestDeb(t *testing.T) {
	opts := NewLogOpts()
	Deb(opts, "This is a debug message")
}

func TestInf(t *testing.T) {
	opts := NewLogOpts()
	Inf(opts, "This is an info message")
}

func TestWar(t *testing.T) {
	opts := NewLogOpts()
	War(opts, "This is a warning message")
}

func TestErr(t *testing.T) {
	opts := NewLogOpts()
	Err(opts, "This is an error message")
}

func TestFat(t *testing.T) {
	opts := NewLogOpts()
	Fat(opts, nil)
}

func TestGetCaller(t *testing.T) {
	c, err := getCaller()
	if err != nil {
		t.Error("getCaller() should not return an error")
	}
	if c == nil {
		t.Error("getCaller() should return a valid pointer to caller")
	}
}

func TestPrintDeb(t *testing.T) {
	opts := NewLogOpts()
	PrintDeb(opts, "This is a debug message")
}

func TestPrintInf(t *testing.T) {
	opts := NewLogOpts()
	PrintInf(opts, "This is an info message")
}

func TestPrintWar(t *testing.T) {
	opts := NewLogOpts()
	PrintWar(opts, "This is a warning message")
}

func TestPrintErr(t *testing.T) {
	opts := NewLogOpts()
	PrintErr(opts, "This is an error message")
}

func TestPrintFat(t *testing.T) {
	opts := NewLogOpts()
	PrintFat(opts, nil)
}

func TestCreateNewLog(t *testing.T) {
	opts := NewLogOpts()
	c, err := getCaller()
	if err != nil {
		t.Error("getCaller() should not return an error")
	}
	err = createNewLog(opts, DEBUG, c, "test")
	if err != nil {
		t.Error("createNewLog() should not return an error")
	}
}

func TestQueryLogs(t *testing.T) {
	_, err := queryLogs()
	if err != nil {
		t.Error("queryLogs() should not return an error")
	}
}

func TestPrintLogs(t *testing.T) {
	PrintLogs()
}

func TestPrintLogsFiltered(t *testing.T) {
	PrintLogs(FilterByStatus(DEBUG))
}

func TestPrintLogsFilteredMulty(t *testing.T) {
	PrintLogs(FilterByStatus(DEBUG), FilterByCallerFile("testing", CONTAINS))
}

func TestPrintLogsSort(t *testing.T) {
	PrintLogs(SortByStatus(DESC))
}

func TestPrintLogsSortFiltered(t *testing.T) {
	PrintLogs(SortByStatus(DESC), FilterByCallerFile("testing", CONTAINS))
}

func TestPrintLogsSortFilteredMulty(t *testing.T) {
	PrintLogs(SortByTimestamp(DESC), FilterByStatus(DEBUG), FilterByCallerFile("testing", CONTAINS))
}

func TestPrintLogsLimit(t *testing.T) {
	PrintLogs(AddLimit(3))
}

func TestPrintLogsLimitOffset(t *testing.T) {
	PrintLogs(AddLimit(3, 2))
}
