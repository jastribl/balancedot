package textscanner

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Scanner is the scanner
type Scanner struct {
	reader    *strings.Reader
	scanner   *bufio.Scanner
	bytesRead int64
}

// NewScanner creates a new Scanner given some text as a byte array
func NewScanner(s []byte) *Scanner {
	reader := strings.NewReader(string(s))
	return &Scanner{
		reader:    reader,
		scanner:   bufio.NewScanner(reader),
		bytesRead: 0,
	}
}

// EatToLine eats lines until and including the desired line
func (m *Scanner) EatToLine(s string) bool {
	bytesToRollbackTo := m.bytesRead

	for m.scanner.Scan() {
		line := m.scanner.Text()
		m.bytesRead += int64(len(line) + 2)
		if line == s {
			return true
		}
	}

	// Roll-Back
	m.reader.Seek(bytesToRollbackTo, io.SeekStart)
	m.scanner = bufio.NewScanner(m.reader)
	m.bytesRead = bytesToRollbackTo

	return false
}

// EatToLineContainsWithCallback eats lines until the line contains the given string and passes lines to fn
func (m *Scanner) EatToLineContainsWithCallback(sub string, fn func(string) error) (bool, error) {
	for m.scanner.Scan() {
		line := m.scanner.Text()
		m.bytesRead += int64(len(line) + 2)
		if strings.Contains(line, sub) {
			return true, fn(line)
		}

	}
	return false, nil
}

// ProcessToAndEatLine eats lines passing them to fn1 until the line contains the given string and passes that line to fn2
func (m *Scanner) ProcessToAndEatLine(sub string, fn1 func(string) error, fn2 func(string) error) error {
	for m.scanner.Scan() {
		line := m.scanner.Text()
		m.bytesRead += int64(len(line) + 2)
		if strings.Contains(line, sub) {
			return fn2(line)
		}

		if err := fn1(line); err != nil {
			return err
		}
	}

	return fmt.Errorf("Unable to find substring string '%s'", sub)
}
