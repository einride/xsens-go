package xsens

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

// MeasurementData contains measurement data read from the Xsens device.
type MeasurementData struct {

	// DataIdentifiers contains a list of all data types present in the measurement data instance.
	DataIdentifiers []DataIdentifier

	// PacketCounter contains the packet counter.
	PacketCounter PacketCounter

	// StatusByte contains the status byte.
	StatusByte StatusByte

	// StatusWord contains the status word.
	StatusWord StatusWord

	// SampleTimeFine contains the fine-grained sample time.
	//
	// See also the documentation for SampleTimeFine.
	SampleTimeFine SampleTimeFine

	// SampleTimeCoarse contains the coarse-grained sample time.
	//
	// See also the documentation for SampleTimeCoarse.
	SampleTimeCoarse SampleTimeCoarse

	// BaroPressure contains the current barometric pressure.
	BaroPressure BaroPressure

	// UTCTime contains the current UTC time.
	UTCTime UTCTime

	// DeltaV contains the delta velocity value of the SDI output in m/s.
	DeltaV VectorXYZ

	// Acceleration contains the calibrated acceleration vector in x, y, and z axes in m/s 2 .
	Acceleration VectorXYZ

	// FreeAcceleration contains the free acceleration vector in x, y, and z axes in m/s 2 .
	FreeAcceleration VectorXYZ

	// AccelerationHR contains the high-resolution calibrated acceleration vector in x, y, and z axes in m/s 2 .
	//
	// For the MTi 1-series, with the exception of the MTi-7, the output data rate is 1000 Hz
	// based on the internal clock of the IMU which is not aligned with other data; data has
	// not been processed in the SDI algorithm. It has been calibrated with the Xsens
	// calibration parameters (except for g-sensitivity).
	//
	// For the MTi-7, the output data is 800 Hz based on the internal clock of the IMUs
	// which are not aligned with other data; data has not been processed in the SDI
	// algorithm. It has been calibrated with the Xsens calibration parameters (except for g-sensitivity).
	//
	// For the MTi 100-series and MTi-G-710, the output data is 1000 Hz, synchronized with
	// the internal clock of the MTi 100-series (10 ppm; 1 ppm with GNSS ClockSync). The
	// data has been processed in the SDI algorithm. Note that AccelerationHR is not
	// grouped with messages coming out at the same time.
	AccelerationHR VectorXYZ

	// DeltaQ contains the delta quaternion value of the SDI output.
	DeltaQ Quaternion

	// RateOfTurn contains the calibrated rate of turn vector in x, y, and z axes in rad/s.
	RateOfTurn VectorXYZ

	// RateOfTurnHR contains the high-resolution calibrated rate of turn vector in x, y, and z axes in rad/s.
	//
	// For the MTi 1-series, with the exception of the MTi-7, the output data rate is 1000 Hz
	// based on the internal clock of the IMU which is not aligned with other data; data has
	// not been processed in the SDI algorithm. It has been calibrated with the Xsens
	// calibration parameters (except for g-sensitivity).
	//
	// For the MTi-7, the output data is 800 Hz based on the internal clock of the IMUs
	// which are not aligned with other data; data has not been processed in the SDI
	// algorithm. It has been calibrated with the Xsens calibration parameters (except for g-
	// sensitivity).
	//
	// For the MTi 100-series and MTi-G-710, the output data is 1000 Hz, synchronized with
	// the internal clock of the MTi 100-series (10 ppm; 1 ppm with GNSS ClockSync). The
	// data has been processed in the SDI algorithm. Note that RateOfTurnHR is not
	// grouped with messages coming out at the same time.
	RateOfTurnHR VectorXYZ

	// Quaternion contains the orientation output expressed as a quaternion
	Quaternion Quaternion

	// EulerAngles contains the three Euler angles in degrees that represent the orientation of the device.
	EulerAngles VectorXYZ

	// RotationMatrix contains the rotation matrix (DCM) that represents the orientation of the MT.
	RotationMatrix RotationMatrix

	// Temperature contains the internal temperature of the sensor in degrees Celsius
	Temperature Scalar

	// AltitudeEllipsoid contains the altitude of the MTi-G in meters above the WGS-84 Ellipsoid.
	AltitudeEllipsoid Scalar

	// PositionECEF contains the position of the MTi-G in the Earth-Centered, Earth-Fixed (ECEF) coordinate
	// system in meters.
	//
	// Note that position in ECEF cannot be represented in Fixed Point values because of the limited range of fixed point
	// representations. Use double or float representation instead.
	PositionECEF VectorXYZ

	// LatLon contains the latitude and longitude in degrees of the MTi-G position
	LatLon LatLon

	// VelocityXYZ contains the X, Y and Z components of the MTi-G velocity in m/s.
	VelocityXYZ VectorXYZ

	// MagneticField contains the magnetic field value in x, y, and z axes in arbitrary units.
	//
	// Magnetic field is normalized to 1.0 during calibration.
	MagneticField VectorXYZ

	// GNSSPVTData contains the current position, velocity and time data.
	GNSSPVTData GNSSPVTData

	// GNSSSatInfo contains info on the currently used GNSS satellites.
	GNSSSatInfo GNSSSatInfo
}

// UnmarshalMTData2 sets *m to the data contained in the provided MTData2 message.
func (m *MeasurementData) UnmarshalMTData2(msg Message) error {
	if msg.Identifier() != MessageIdentifierMTData2 {
		return xerrors.Errorf("message is not %v", MessageIdentifierMTData2)
	}
	*m = MeasurementData{DataIdentifiers: m.DataIdentifiers[:0]}
	i := 0
	data := MTData2(msg.Data())
	for i < len(data) {
		packet, err := data.PacketAt(i)
		if err != nil {
			return err
		}
		i += len(packet)
		m.DataIdentifiers = append(m.DataIdentifiers, packet.Identifier())
		if measurementDataType, ok := m.getMeasurementDataType(packet.Identifier().DataType); ok {
			if err := measurementDataType.unmarshalMTData2Packet(packet); err != nil {
				return xerrors.Errorf("packet: %v: %w", packet, err)
			}
		}
	}
	return nil
}

func (m *MeasurementData) MarshalText() (string, error) {
	var buf bytes.Buffer
	for _, dataIdentifier := range m.DataIdentifiers {
		if _, err := fmt.Fprintf(&buf, "%v:\n", dataIdentifier); err != nil {
			return "", err
		}
		if data, ok := m.getMeasurementDataType(dataIdentifier.DataType); ok {
			// pretty-print
			dataText := fmt.Sprintf("%+v", data)
			// strip braces
			dataText = strings.TrimSuffix(strings.TrimPrefix(dataText, "&{"), "}")
			// add line breaks after each field
			dataText = strings.Replace(dataText, " ", "\n  ", -1)
			if _, err := fmt.Fprintf(&buf, "  %s\n", dataText); err != nil {
				return "", err
			}
		}
	}
	return strings.TrimSpace(buf.String()), nil
}

type MeasurementDataType interface {
	unmarshalMTData2Packet(MTData2Packet) error
}

func (m *MeasurementData) getMeasurementDataType(dataType DataType) (MeasurementDataType, bool) {
	switch dataType {
	case DataTypeDeltaV:
		return &m.DeltaV, true
	case DataTypeAcceleration:
		return &m.Acceleration, true
	case DataTypeFreeAcceleration:
		return &m.FreeAcceleration, true
	case DataTypeAccelerationHR:
		return &m.AccelerationHR, true
	case DataTypeDeltaQ:
		return &m.DeltaQ, true
	case DataTypeRateOfTurn:
		return &m.RateOfTurn, true
	case DataTypeRateOfTurnHR:
		return &m.RateOfTurnHR, true
	case DataTypeQuaternion:
		return &m.Quaternion, true
	case DataTypeEulerAngles:
		return &m.EulerAngles, true
	case DataTypeRotationMatrix:
		return &m.RotationMatrix, true
	case DataTypeTemperature:
		return &m.Temperature, true
	case DataTypeAltitudeEllipsoid:
		return &m.AltitudeEllipsoid, true
	case DataTypePositionECEF:
		return &m.PositionECEF, true
	case DataTypeLatLon:
		return &m.LatLon, true
	case DataTypeVelocityXYZ:
		return &m.VelocityXYZ, true
	case DataTypeStatusByte:
		return &m.StatusByte, true
	case DataTypeStatusWord:
		return &m.StatusWord, true
	case DataTypeUTCTime:
		return &m.UTCTime, true
	case DataTypePacketCounter:
		return &m.PacketCounter, true
	case DataTypeSampleTimeFine:
		return &m.SampleTimeFine, true
	case DataTypeSampleTimeCoarse:
		return &m.SampleTimeCoarse, true
	case DataTypeBaroPressure:
		return &m.BaroPressure, true
	case DataTypeMagneticField:
		return &m.MagneticField, true
	case DataTypeGNSSPVTData:
		return &m.GNSSPVTData, true
	case DataTypeGNSSSatInfo:
		return &m.GNSSSatInfo, true
	}
	return nil, false
}
