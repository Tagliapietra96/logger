package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tagliapietra96/tui"
	"github.com/Tagliapietra96/tui/opts"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

func printLogs(lopts *Logger, logs []*log) {
	var strLogs []string
	w := 100

	if lopts.inline {
		w = 130
	}

	tw, _, err := term.GetSize(os.Stdout.Fd())
	if tw > 0 && tw < w && err == nil {
		w = tw - 4
	}

	page := tui.NewStyle(opts.Margin(1, 2, 1, 1), opts.Width(w))
	if lopts.inline {
		strLogs = getInlineLogs(w, lopts, logs)
	} else {
		strLogs = getBlockLogs(w, lopts, logs)
	}

	tui.Concat(&page, strLogs...)
	fmt.Print(page.String())
	println("")
}

func getInlineLogs(w int, lopts *Logger, logs []*log) []string {
	var lw, tw, cw, tgw, mw int

	if w <= 75 && lopts.showTimestamp == ShowFullTimestamp {
		lopts.Timestamp(ShowDateTime)
	}

	levels := make([]string, 0, len(logs))
	timestamps := make([]string, 0, len(logs))
	callers := make([]string, 0, len(logs))
	tags := make([]string, 0, len(logs))
	messages := make([]string, 0, len(logs))

	for _, log := range logs {
		level := log.level.toString()
		timestamp := log.timestamp.toString(lopts.showTimestamp)
		caller := log.getCaller(lopts.inline, lopts.showCaller)
		tag := ""
		if lopts.showTags && len(log.tags) > 0 {
			tag = strings.Join(log.getTags(), ", ")
			if tgw < lipgloss.Width(tag)+2 {
				tgw = lipgloss.Width(tag) + 2
			}
		}

		if lw < lipgloss.Width(level)+2 {
			lw = lipgloss.Width(level) + 2
		}

		if lopts.showTimestamp != HideTimestamp {
			if tw < lipgloss.Width(timestamp)+2 {
				tw = lipgloss.Width(timestamp) + 2
			}
		}

		if lopts.showCaller != HideCaller {
			if cw < lipgloss.Width(caller)+2 {
				cw = lipgloss.Width(caller) + 2
			}
		}

		if mw < lipgloss.Width(log.message)+1 {
			mw = lipgloss.Width(log.message) + 1
		}

		levels = append(levels, level)
		timestamps = append(timestamps, timestamp)
		callers = append(callers, caller)
		tags = append(tags, tag)
		messages = append(messages, log.message)
	}

	if w <= 75 {
		lopts.showTags = false
		mw += tgw
		tgw = 0
	}

	if w <= 60 {
		lopts.Caller(HideCaller)
		mw += cw
		cw = 0
	}

	if lw+tw+cw+tgw+mw > w {
		for lw+tw+cw+tgw+mw > w {
			if tw > 12 {
				tw--
			}

			if lopts.showCaller > ShowCallerLine {
				if cw > 1 {
					cw--
				}
			}

			if mw > 1 {
				mw--
			}
		}
	}

	rows := make([]string, 0, len(logs))

	for i := range len(logs) {
		var ts, lvl, cl, tg, msg string
		row := tui.NewStyle(opts.Color(nil, nil, tui.ColorMuted))
		if i != 0 {
			row = row.Border(lipgloss.NormalBorder(), true, false, false, false)
		}

		if lopts.showTimestamp != HideTimestamp {
			ts = tui.Render(timestamps[i], opts.Width(tw), opts.Muted)
		}

		if lopts.showCaller != HideCaller {
			cl = tui.Render(callers[i], opts.Width(cw), opts.Muted)
		}

		if lopts.showTags {
			tg = tui.Render(tags[i], opts.Width(tgw), opts.LightMuted)
		}

		lvl = tui.Render(levels[i], opts.Width(lw), opts.Color(logs[i].level.color()))
		msg = tui.Render(messages[i], opts.Width(mw))
		rows = append(rows, row.Render(lipgloss.JoinHorizontal(lipgloss.Top, ts, tg, lvl, cl, msg)))
	}

	return rows
}

func getBlockLogs(w int, lopts *Logger, logs []*log) []string {
	result := make([]string, 0, len(logs))
	for _, log := range logs {
		var timestamp, caller, tags string
		l := tui.NewStyle(opts.Padding(0, 1))
		l = l.Border(lipgloss.RoundedBorder(), true)
		tui.Config(&l, opts.FitWidth(w))
		color := log.level.color()

		tui.Config(&l, opts.Color(nil, nil, color))

		logTitle := tui.NewStyle(opts.Color(nil, nil, tui.ColorMuted), opts.Width(w-4)).Border(lipgloss.NormalBorder(), false, false, true, false)
		level := log.level.toString()

		if lopts.showTimestamp != HideTimestamp {
			timestamp = tui.Render(log.timestamp.toString(lopts.showTimestamp), opts.Right)
		}

		if lopts.showCaller != HideCaller {
			caller = log.getCaller(lopts.inline, lopts.showCaller)
		}

		if lopts.showTags && len(log.tags) > 0 {
			tags = tui.Render(strings.Join(log.getTags(), " ･ "))
		}

		var titlefirtsRow, titleSecondRow string
		if w-4-lipgloss.Width(level)-lipgloss.Width(timestamp) > 0 {
			titlefirtsRow = lipgloss.JoinHorizontal(lipgloss.Top, level, lipgloss.PlaceHorizontal(w-4-lipgloss.Width(level)-lipgloss.Width(timestamp), lipgloss.Center, ""), timestamp)
		} else {
			titlefirtsRow = level
			if timestamp != "" {
				titlefirtsRow += "\n" + timestamp
			}
		}

		if w-4-lipgloss.Width(caller)-lipgloss.Width(tags) > 0 {
			titleSecondRow = lipgloss.JoinHorizontal(lipgloss.Top, caller, lipgloss.PlaceHorizontal(w-4-lipgloss.Width(caller)-lipgloss.Width(tags), lipgloss.Center, ""), tags)
		} else {
			titleSecondRow = caller + "\n" + tags
		}

		tui.ConcatLn(&logTitle, titlefirtsRow, titleSecondRow)

		message := tui.Render(log.message, opts.Left, opts.Padding(1, 0), opts.Width(w-4))
		tui.Concat(&l, logTitle.String(), message)
		result = append(result, l.String())
	}

	return result
}
