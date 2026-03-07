package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"

	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Gray    = "\033[90m"

	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
)

// Symbols for terminal output
const (
	SymbolSuccess = "✓"
	SymbolError   = "✗"
	SymbolWarning = "⚠"
	SymbolInfo    = "ℹ"
	SymbolArrow   = "→"
	SymbolDot     = "•"
)

var colorsEnabled = true

func init() {
	// Disable colors if not a terminal or NO_COLOR is set
	if !term.IsTerminal(int(os.Stdout.Fd())) || os.Getenv("NO_COLOR") != "" {
		colorsEnabled = false
	}
}

// DisableColors disables colored output
func DisableColors() {
	colorsEnabled = false
}

// EnableColors enables colored output
func EnableColors() {
	colorsEnabled = true
}

// colorize wraps text with color codes if colors are enabled
func colorize(color, text string) string {
	if !colorsEnabled {
		return text
	}
	return color + text + Reset
}

// Color functions

func RedText(text string) string {
	return colorize(Red, text)
}

func GreenText(text string) string {
	return colorize(Green, text)
}

func YellowText(text string) string {
	return colorize(Yellow, text)
}

func BlueText(text string) string {
	return colorize(Blue, text)
}

func CyanText(text string) string {
	return colorize(Cyan, text)
}

func GrayText(text string) string {
	return colorize(Gray, text)
}

func BoldText(text string) string {
	return colorize(Bold, text)
}

func DimText(text string) string {
	return colorize(Dim, text)
}

// Status output functions

func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", GreenText(SymbolSuccess), msg)
}

func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", RedText(SymbolError), msg)
}

func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", YellowText(SymbolWarning), msg)
}

func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", BlueText(SymbolInfo), msg)
}

func Step(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("  %s %s\n", GrayText(SymbolArrow), msg)
}

func SubStep(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("    %s %s\n", GrayText(SymbolDot), msg)
}
