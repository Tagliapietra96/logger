package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Tagliapietra96/tui"
	"github.com/Tagliapietra96/tui/opts"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

// getTerminalSize function returns the width and height of the terminal.
// It returns the width and height of the terminal as integers.
// If the terminal size cannot be determined, it returns 0, 0.
func getTerminalSize() (int, int) {
	w, h, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return 0, 0
	}

	return w, h
}

func printLogs(logs []*log) {
	w := 100
	tw, _ := getTerminalSize()
	if tw > 0 && tw < w {
		w = tw - 4
	}
	page := tui.NewStyle(opts.Margin(1, 2, 1, 1), opts.Width(w))
	for _, log := range logs {
		l := tui.NewStyle(opts.Padding(0, 1))
		l = l.Border(lipgloss.RoundedBorder(), true)
		tui.Config(&l, opts.FitWidth(w))
		color := tui.ColorMuted
		switch log.level {
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
		}

		tui.Config(&l, opts.Color(nil, nil, color))

		logTitle := tui.NewStyle(opts.Color(nil, nil, tui.ColorMuted), opts.Width(w-4)).Border(lipgloss.NormalBorder(), false, false, true, false)
		level := tui.Render(log.level.String(), opts.Color(color))
		timestamp := tui.Render(log.timestamp.String(), opts.Color(tui.ColorMuted), opts.Right)
		caller := tui.Render(fmt.Sprintf("at %s:%d - %s", filepath.Base(log.callerFile), log.callerLine, log.callerFunction), opts.Color(tui.ColorMuted), opts.Left)
		tags := tui.NewStyle(opts.Color(tui.ColorLightMuted), opts.Right)
		for _, tag := range log.tags {
			tui.ConcatWith(&tags, " ･ ", fmt.Sprintf("🔖 %s", tag))
		}
		tui.ConcatLn(
			&logTitle,
			lipgloss.JoinHorizontal(lipgloss.Top,
				level,
				lipgloss.PlaceHorizontal(w-4-lipgloss.Width(level)-lipgloss.Width(timestamp), lipgloss.Center, ""),
				timestamp,
			),
			lipgloss.JoinHorizontal(lipgloss.Top,
				caller,
				lipgloss.PlaceHorizontal(w-4-lipgloss.Width(caller)-lipgloss.Width(tags.String()), lipgloss.Center, ""),
				tags.String(),
			),
		)

		message := tui.Render(log.message, opts.Left, opts.Padding(1, 0), opts.Width(w-4))
		tui.Concat(&l, logTitle.String(), message)
		tui.Concat(&page, l.String())
	}

	fmt.Print(page.String())
	println("")
}
