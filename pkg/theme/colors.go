package theme

import "github.com/fatih/color"

var (
	Info    = color.New(color.FgCyan).SprintFunc()
	Warn    = color.New(color.FgYellow).SprintFunc()
	Error   = color.New(color.FgRed).SprintFunc()
	Success = color.New(color.FgGreen).SprintFunc()
	Title   = color.New(color.Bold, color.FgHiWhite).SprintFunc()
	Unimportant = color.New(color.Faint, color.FgHiBlack).SprintFunc()
)
