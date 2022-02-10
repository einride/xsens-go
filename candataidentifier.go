package xsens

import "fmt"

//go:generate stringer -type CANDataIdentifier -trimprefix CANDataIdentifier

type CANDataIdentifier uint8

const (
	CANDataIdentifierInvalid            CANDataIdentifier = 0x00
	CANDataIdentifierError              CANDataIdentifier = 0x01
	CANDataIdentifierWarning            CANDataIdentifier = 0x02
	CANDataIdentifierSampleTime         CANDataIdentifier = 0x05
	CANDataIdentifierGroupCounter       CANDataIdentifier = 0x06
	CANDataIdentifierUtcTime            CANDataIdentifier = 0x07
	CANDataIdentifierStatusWord         CANDataIdentifier = 0x11
	CANDataIdentifierQuaternion         CANDataIdentifier = 0x21
	CANDataIdentifierEulerAngles        CANDataIdentifier = 0x22
	CANDataIdentifierDeltaV             CANDataIdentifier = 0x31
	CANDataIdentifierRateOfTurn         CANDataIdentifier = 0x32
	CANDataIdentifierDeltaQ             CANDataIdentifier = 0x33
	CANDataIdentifierAcceleration       CANDataIdentifier = 0x34
	CANDataIdentifierFreeAcceleration   CANDataIdentifier = 0x35
	CANDataIdentifierMagneticField      CANDataIdentifier = 0x41
	CANDataIdentifierTemperature        CANDataIdentifier = 0x51
	CANDataIdentifierBaroPressure       CANDataIdentifier = 0x52
	CANDataIdentifierRateOfTurnHR       CANDataIdentifier = 0x61
	CANDataIdentifierAccelerationHR     CANDataIdentifier = 0x62
	CANDataIdentifierLatLong            CANDataIdentifier = 0x71
	CANDataIdentifierAltitudeEllipsoid  CANDataIdentifier = 0x72
	CANDataIdentifierPositionEcefX      CANDataIdentifier = 0x73
	CANDataIdentifierPositionEcefY      CANDataIdentifier = 0x74
	CANDataIdentifierPositionEcefZ      CANDataIdentifier = 0x75
	CANDataIdentifierVelocityXYZ        CANDataIdentifier = 0x76
	CANDataIdentifierGnssReceiverStatus CANDataIdentifier = 0x79
	CANDataIdentifierGnssReceiverDop    CANDataIdentifier = 0x7A
)

func (i *CANDataIdentifier) UnmarshalText(text []byte) error {
	knownIDs := []CANDataIdentifier{
		CANDataIdentifierInvalid,
		CANDataIdentifierFreeAcceleration,
		CANDataIdentifierError,
		CANDataIdentifierMagneticField,
		CANDataIdentifierWarning,
		CANDataIdentifierTemperature,
		CANDataIdentifierSampleTime,
		CANDataIdentifierBaroPressure,
		CANDataIdentifierGroupCounter,
		CANDataIdentifierRateOfTurnHR,
		CANDataIdentifierUtcTime,
		CANDataIdentifierAccelerationHR,
		CANDataIdentifierStatusWord,
		CANDataIdentifierLatLong,
		CANDataIdentifierQuaternion,
		CANDataIdentifierAltitudeEllipsoid,
		CANDataIdentifierEulerAngles,
		CANDataIdentifierPositionEcefX,
		CANDataIdentifierPositionEcefY,
		CANDataIdentifierDeltaV,
		CANDataIdentifierPositionEcefZ,
		CANDataIdentifierRateOfTurn,
		CANDataIdentifierVelocityXYZ,
		CANDataIdentifierDeltaQ,
		CANDataIdentifierGnssReceiverStatus,
		CANDataIdentifierAcceleration,
		CANDataIdentifierGnssReceiverDop,
	}
	for _, d := range knownIDs {
		if d.String() == string(text) {
			*i = d
			return nil
		}
	}
	return fmt.Errorf("unknown CANDataIdentifier %s", string(text))
}
