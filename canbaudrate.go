package xsens

//go:generate stringer -type CANBaudRate -trimprefix CANBaudRate

type CANBaudRate int8

const (
	CANBaudRate1M   CANBaudRate = 0x0C
	CANBaudRate800k CANBaudRate = 0x0B
	CANBaudRate500k CANBaudRate = 0x0A
	CANBaudRate250k CANBaudRate = 0x00
	CANBaudRate125k CANBaudRate = 0x01
	CANBaudRate100k CANBaudRate = 0x02
	CANBaudRate83k3 CANBaudRate = 0x03
	CANBaudRate62k5 CANBaudRate = 0x04
	CANBaudRate50k  CANBaudRate = 0x05
	CANBaudRate33k3 CANBaudRate = 0x06
	CANBaudRate20k  CANBaudRate = 0x07
	CANBaudRate10k  CANBaudRate = 0x08
	CANBaudRate5k   CANBaudRate = 0x09
)
