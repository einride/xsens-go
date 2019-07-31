package xsens

// DataType represents an Xsens data type.
type DataType uint16

// Data group: Temperature.
const (
	DataTypeTemperature DataType = 0x0810
)

// Data group: Timestamp.
const (
	DataTypeUTCTime          DataType = 0x1010
	DataTypePacketCounter    DataType = 0x1020
	DataTypeITOW             DataType = 0x1030
	DataTypeGPSAge           DataType = 0x1040
	DataTypePressureAge      DataType = 0x1050
	DataTypeSampleTimeFine   DataType = 0x1060
	DataTypeSampleTimeCoarse DataType = 0x1070
)

// Data group: Orientation.
const (
	DataTypeQuaternion     DataType = 0x2010
	DataTypeRotationMatrix DataType = 0x2020
	DataTypeEulerAngles    DataType = 0x2030
)

// Data group: Pressure.
const (
	DataTypeBaroPressure DataType = 0x3010
)

// Data group: Acceleration.
const (
	DataTypeDeltaV           DataType = 0x4010
	DataTypeAcceleration     DataType = 0x4020
	DataTypeFreeAcceleration DataType = 0x4030
	DataTypeAccelerationHR   DataType = 0x4040
)

// Data group: Position.
const (
	DataTypeAltitudeEllipsoid DataType = 0x5020
	DataTypePositionECEF      DataType = 0x5030
	DataTypeLatLon            DataType = 0x5040
)

// Data group: GNSS.
const (
	DataTypeGNSSPVTData DataType = 0x7010
	DataTypeGNSSSatInfo DataType = 0x7020
)

// Data group: Angular velocity.
const (
	DataTypeRateOfTurn   DataType = 0x8020
	DataTypeDeltaQ       DataType = 0x8030
	DataTypeRateOfTurnHR DataType = 0x8040
)

// Data group: GPS.
const (
	DataTypeGPSDOP     DataType = 0x8830
	DataTypeGPSSOL     DataType = 0x8840
	DataTypeGPSTimeUTC DataType = 0x8880
	DataTypeGPSSVInfo  DataType = 0x88a0
)

// Data group: Magnetic.
const (
	DataTypeMagneticField DataType = 0xc020
)

// Data group: Velocity.
const (
	DataTypeVelocityXYZ DataType = 0xd010
)

// Data group: Status.
const (
	DataTypeStatusByte DataType = 0xe010
	DataTypeStatusWord DataType = 0xe020
)

// HasPrecision returns true for data types which support configurable output precision.
func (d DataType) HasPrecision() bool {
	switch d {
	case
		// temperature
		DataTypeTemperature,
		// orientation
		DataTypeQuaternion, DataTypeRotationMatrix, DataTypeEulerAngles,
		// acceleration
		DataTypeDeltaV, DataTypeAcceleration, DataTypeFreeAcceleration, DataTypeAccelerationHR,
		// position
		DataTypeAltitudeEllipsoid, DataTypePositionECEF, DataTypeLatLon,
		// angular velocity
		DataTypeDeltaQ, DataTypeRateOfTurn, DataTypeRateOfTurnHR,
		// velocity
		DataTypeVelocityXYZ,
		// magnetic
		DataTypeMagneticField:
		return true
	default:
		return false
	}
}

// HasCoordinateSystem returns true for data types which support configurable coordinate system.
func (d DataType) HasCoordinateSystem() bool {
	switch d {
	case
		// orientation
		DataTypeQuaternion, DataTypeRotationMatrix, DataTypeEulerAngles,
		// velocity
		DataTypeVelocityXYZ:
		return true
	default:
		return false
	}
}
