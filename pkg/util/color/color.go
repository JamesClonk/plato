package color

import "github.com/fatih/color"

func Fail(format string, args ...any) string {
	return stringify(color.New(color.FgRed, color.BlinkSlow, color.Underline), format, args...)
}

func Red(format string, args ...any) string {
	return stringify(color.New(color.FgRed), format, args...)
}

func Green(format string, args ...any) string {
	return stringify(color.New(color.FgGreen), format, args...)
}

func Blue(format string, args ...any) string {
	return stringify(color.New(color.FgBlue), format, args...)
}

func Magenta(format string, args ...any) string {
	return stringify(color.New(color.FgMagenta), format, args...)
}

func Cyan(format string, args ...any) string {
	return stringify(color.New(color.FgCyan), format, args...)
}

func Yellow(format string, args ...any) string {
	return stringify(color.New(color.FgYellow), format, args...)
}

func White(format string, args ...any) string {
	return stringify(color.New(color.FgWhite), format, args...)
}

func Black(format string, args ...any) string {
	return stringify(color.New(color.FgBlack), format, args...)
}

func stringify(c *color.Color, format string, args ...any) string {
	sprint := c.SprintfFunc()
	if len(args) == 0 {
		return sprint(format)
	}
	return sprint(format, args...)
}
