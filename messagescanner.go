package xsens

import (
	"bufio"
	"io"
)

// MessageScanner scans for Xsens messages.
type MessageScanner struct {
	sc *bufio.Scanner
}

// NewMessageScanner creates a message scanner that scans for Xsens messages in the provided io.Reader.
func NewMessageScanner(r io.Reader) *MessageScanner {
	m := &MessageScanner{sc: bufio.NewScanner(r)}
	m.sc.Split(ScanMessages)
	m.sc.Buffer(make([]byte, 0, maxLengthOfMessage), maxLengthOfMessage)
	return m
}

// Scan to the next message.
func (m *MessageScanner) Scan() bool {
	return m.sc.Scan()
}

// Message returns the latest scanned message.
func (m *MessageScanner) Message() Message {
	return Message(m.sc.Bytes())
}

// Err returns the error encountered while scanning, if any.
func (m *MessageScanner) Err() error {
	return m.sc.Err()
}
