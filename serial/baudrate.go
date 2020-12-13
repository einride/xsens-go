package serial

// BaudRate represents a serial communication baud rate.
type BaudRate int

//go:generate stringer -type BaudRate -trimprefix BaudRate

const (
	BaudRate4800    BaudRate = 4800
	BaudRate9600    BaudRate = 9600
	BaudRate19200   BaudRate = 19200
	BaudRate38400   BaudRate = 38400
	BaudRate57600   BaudRate = 57600
	BaudRate115200  BaudRate = 115200
	BaudRate230400  BaudRate = 230400
	BaudRate460800  BaudRate = 460800
	BaudRate921600  BaudRate = 921600
	BaudRate2000000 BaudRate = 2000000
)
