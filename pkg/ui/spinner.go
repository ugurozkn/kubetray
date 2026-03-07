package ui

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner represents a terminal spinner
type Spinner struct {
	message  string
	writer   io.Writer
	done     chan struct{}
	mu       sync.Mutex
	active   bool
	frameIdx int
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		writer:  os.Stdout,
		done:    make(chan struct{}),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	// Don't animate if not a terminal
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Fprintf(s.writer, "%s %s...\n", GrayText(spinnerFrames[0]), s.message)
		return
	}

	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.done:
				return
			case <-ticker.C:
				s.mu.Lock()
				if !s.active {
					s.mu.Unlock()
					return
				}
				frame := spinnerFrames[s.frameIdx]
				s.frameIdx = (s.frameIdx + 1) % len(spinnerFrames)
				fmt.Fprintf(s.writer, "\r%s %s", CyanText(frame), s.message)
				s.mu.Unlock()
			}
		}
	}()
}

// Stop stops the spinner without any status
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return
	}
	s.active = false
	close(s.done)
	s.clearLine()
}

// Success stops the spinner with a success message
func (s *Spinner) Success(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		fmt.Fprintf(s.writer, "%s %s\n", GreenText(SymbolSuccess), message)
		return
	}
	s.active = false
	close(s.done)
	s.clearLine()
	fmt.Fprintf(s.writer, "%s %s\n", GreenText(SymbolSuccess), message)
}

// Error stops the spinner with an error message
func (s *Spinner) Error(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		fmt.Fprintf(s.writer, "%s %s\n", RedText(SymbolError), message)
		return
	}
	s.active = false
	close(s.done)
	s.clearLine()
	fmt.Fprintf(s.writer, "%s %s\n", RedText(SymbolError), message)
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

func (s *Spinner) clearLine() {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Fprint(s.writer, "\r\033[K")
	}
}

// WithSpinner executes a function while showing a spinner
func WithSpinner(message string, fn func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()

	err := fn()

	if err != nil {
		spinner.Error(message + " - " + RedText("failed"))
		return err
	}

	spinner.Success(message)
	return nil
}
