package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Table represents a simple table for terminal output
type Table struct {
	headers []string
	rows    [][]string
	writer  io.Writer
}

// NewTable creates a new table with the given headers
func NewTable(headers ...string) *Table {
	return &Table{
		headers: headers,
		rows:    make([][]string, 0),
		writer:  os.Stdout,
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(cols ...string) {
	t.rows = append(t.rows, cols)
}

// Render prints the table
func (t *Table) Render() {
	if len(t.headers) == 0 {
		return
	}

	// Calculate column widths
	widths := make([]int, len(t.headers))
	for i, h := range t.headers {
		widths[i] = len(h)
	}
	for _, row := range t.rows {
		for i, col := range row {
			if i < len(widths) && len(col) > widths[i] {
				widths[i] = len(col)
			}
		}
	}

	// Print headers
	headerParts := make([]string, len(t.headers))
	for i, h := range t.headers {
		headerParts[i] = BoldText(padRight(strings.ToUpper(h), widths[i]))
	}
	fmt.Fprintln(t.writer, strings.Join(headerParts, "  "))

	// Print rows
	for _, row := range t.rows {
		rowParts := make([]string, len(t.headers))
		for i := range t.headers {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			rowParts[i] = padRight(val, widths[i])
		}
		fmt.Fprintln(t.writer, strings.Join(rowParts, "  "))
	}
}

func padRight(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}

// Box prints content in a simple box
func Box(title string, lines ...string) {
	maxLen := len(title)
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	border := strings.Repeat("─", maxLen+2)

	fmt.Printf("┌%s┐\n", border)
	if title != "" {
		fmt.Printf("│ %s │\n", BoldText(padRight(title, maxLen)))
		fmt.Printf("├%s┤\n", border)
	}
	for _, line := range lines {
		fmt.Printf("│ %s │\n", padRight(line, maxLen))
	}
	fmt.Printf("└%s┘\n", border)
}

// Header prints a section header
func Header(text string) {
	fmt.Printf("\n%s\n", BoldText(text))
	fmt.Println(strings.Repeat("─", len(text)))
}

// Blank prints an empty line
func Blank() {
	fmt.Println()
}
