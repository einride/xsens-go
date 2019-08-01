package serial

// BaudRate represents a serial communication baud rate.
type BaudRate int

const (
	BaudRate4800 BaudRate = iota
	BaudRate9600
	BaudRate19200
	BaudRate38400
	BaudRate57600
	BaudRate115200
	BaudRate230400
	BaudRate460800
	BaudRate921600
	BaudRate2000000
)
