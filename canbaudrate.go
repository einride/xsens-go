package xsens

//go:generate stringer -type CANBaudRateID -trimprefix CANBaudRateID

type (
	CANBaudRateID int8
	CANBaudRate   int
)

const (
	CANBaudRate1M   CANBaudRateID = 0x0C
	CANBaudRate800k CANBaudRateID = 0x0B
	CANBaudRate500k CANBaudRateID = 0x0A
	CANBaudRate250k CANBaudRateID = 0x00
	CANBaudRate125k CANBaudRateID = 0x01
	CANBaudRate100k CANBaudRateID = 0x02
	CANBaudRate83k3 CANBaudRateID = 0x03
	CANBaudRate62k5 CANBaudRateID = 0x04
	CANBaudRate50k  CANBaudRateID = 0x05
	CANBaudRate33k3 CANBaudRateID = 0x06
	CANBaudRate20k  CANBaudRateID = 0x07
	CANBaudRate10k  CANBaudRateID = 0x08
	CANBaudRate5k   CANBaudRateID = 0x09
)

func (c CANBaudRate) ID() (CANBaudRateID, error) {
	switch c {
	case 5000:
		return CANBaudRate5k, nil
	case 10000:
		return CANBaudRate10k, nil
	case 20000:
		return CANBaudRate20k, nil
	case 33300:
		return CANBaudRate33k3, nil
	case 50000:
		return CANBaudRate50k, nil
	case 62500:
		return CANBaudRate62k5, nil
	case 83300:
		return CANBaudRate83k3, nil
	case 100000:
		return CANBaudRate100k, nil
	case 125000:
		return CANBaudRate125k, nil
	case 250000:
		return CANBaudRate250k, nil
	case 500000:
		return CANBaudRate500k, nil
	case 800000:
		return CANBaudRate800k, nil
	case 1000000:
		return CANBaudRate1M, nil
	default:
		return -1, nil
	}
}
