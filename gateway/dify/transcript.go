package dify

import (
	"strings"
	"sync"
)

// Line is one relayed message of a transcript.
type Line struct {
	Username string
	Text     string
}

// Transcript keeps the most recent relayed messages of a gateway in memory.
// It is safe for concurrent use.
type Transcript struct {
	mu    sync.Mutex
	lines []Line
	next  int
	full  bool
}

// NewTranscript returns a transcript keeping the last size messages.
func NewTranscript(size int) *Transcript {
	return &Transcript{lines: make([]Line, size)}
}

// Add records a message, dropping the oldest one when the transcript is full.
func (t *Transcript) Add(l Line) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.lines[t.next] = l
	t.next = (t.next + 1) % len(t.lines)

	if t.next == 0 {
		t.full = true
	}
}

// Lines returns the recorded messages, oldest first.
func (t *Transcript) Lines() []Line {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.full {
		return append([]Line(nil), t.lines[:t.next]...)
	}

	return append(append([]Line(nil), t.lines[t.next:]...), t.lines[:t.next]...)
}

// Format renders lines as "username: text", one message per line.
func Format(lines []Line) string {
	var b strings.Builder

	for _, l := range lines {
		b.WriteString(l.Username)
		b.WriteString(": ")
		b.WriteString(l.Text)
		b.WriteByte('\n')
	}

	return b.String()
}
