package xsens

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"
)

// MeasurementData is a generic interface for any measurement data produced by an Xsens device.
type MeasurementData interface {
	UnmarshalMTData2Packet(MTData2Packet) error
	MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error)
}

// Scalar contains a single scalar value.
type Scalar float64

// String returns a string representation of the scalar.
func (s *Scalar) String() string {
	return strconv.FormatFloat(float64(*s), 'f', -1, 64)
}

func (s *Scalar) UnmarshalMTData2Packet(packet MTData2Packet) error {
	var err error
	switch packet.Identifier().Precision {
	case PrecisionFloat32:
		var value float32
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &value); err == nil {
			*s = Scalar(value)
		}
	case PrecisionFP1220:
		var value FP1220
		if err = value.fromBinary(packet.Data()); err == nil {
			*s = Scalar(value.Float64())
		}
	case PrecisionFP1632:
		var value FP1632
		if err = value.fromBinary(packet.Data()); err == nil {
			*s = Scalar(value.Float64())
		}
	case PrecisionFloat64:
		err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, s)
	default:
		err = fmt.Errorf("invalid precision: %v", packet.Identifier().Precision)
	}
	if err != nil {
		return fmt.Errorf("precision %v: %w", packet.Identifier().Precision, err)
	}
	return nil
}

func (s *Scalar) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(id.Precision.Size(), id)
	switch id.Precision {
	case PrecisionFloat32:
		binary.BigEndian.PutUint32(packet.Data(), math.Float32bits(float32(*s)))
	case PrecisionFP1220:
		var f FP1220
		f.FromFloat64(float64(*s))
		f.toBinary(packet.Data())
	case PrecisionFP1632:
		var f FP1632
		f.FromFloat64(float64(*s))
		f.toBinary(packet.Data())
	case PrecisionFloat64:
		binary.BigEndian.PutUint64(packet.Data(), math.Float64bits(float64(*s)))
	}
	return packet, nil
}

// VectorXYZ contains a vector with x, y and z-components.
type VectorXYZ struct {
	X, Y, Z float64
}

func (t *VectorXYZ) UnmarshalMTData2Packet(packet MTData2Packet) error {
	var err error
	switch packet.Identifier().Precision {
	case PrecisionFloat32:
		fields := struct {
			X, Y, Z float32
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.X = float64(fields.X)
			t.Y = float64(fields.Y)
			t.Z = float64(fields.Z)
		}
	case PrecisionFP1220:
		fields := struct {
			X, Y, Z FP1220
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.X = fields.X.Float64()
			t.Y = fields.Y.Float64()
			t.Z = fields.Z.Float64()
		}
	case PrecisionFP1632:
		fields := struct {
			X, Y, Z FP1632
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.X = fields.X.Float64()
			t.Y = fields.Y.Float64()
			t.Z = fields.Z.Float64()
		}
	case PrecisionFloat64:
		err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, t)
	default:
		err = fmt.Errorf("invalid precision: %v", packet.Identifier().Precision)
	}
	if err != nil {
		return fmt.Errorf("precision %v: %w", packet.Identifier().Precision, err)
	}
	return nil
}

func (t *VectorXYZ) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(id.Precision.Size()*3, id)
	switch id.Precision {
	case PrecisionFloat32:
		binary.BigEndian.PutUint32(packet.Data(), math.Float32bits(float32(t.X)))
		binary.BigEndian.PutUint32(packet.Data()[id.Precision.Size()*1:], math.Float32bits(float32(t.Y)))
		binary.BigEndian.PutUint32(packet.Data()[id.Precision.Size()*2:], math.Float32bits(float32(t.Z)))
	case PrecisionFP1220:
		var x, y, z FP1220
		x.FromFloat64(t.X)
		y.FromFloat64(t.Y)
		z.FromFloat64(t.Z)
		x.toBinary(packet.Data())
		y.toBinary(packet.Data()[id.Precision.Size():])
		z.toBinary(packet.Data()[id.Precision.Size()*2:])
	case PrecisionFP1632:
		var x, y, z FP1632
		x.FromFloat64(t.X)
		y.FromFloat64(t.Y)
		z.FromFloat64(t.Z)
		x.toBinary(packet.Data())
		y.toBinary(packet.Data()[id.Precision.Size():])
		z.toBinary(packet.Data()[id.Precision.Size()*2:])
	case PrecisionFloat64:
		binary.BigEndian.PutUint64(packet.Data(), math.Float64bits((t.X)))
		binary.BigEndian.PutUint64(packet.Data()[id.Precision.Size()*1:], math.Float64bits((t.Y)))
		binary.BigEndian.PutUint64(packet.Data()[id.Precision.Size()*2:], math.Float64bits((t.Z)))
	}
	return packet, nil
}

// Quaternion contains a quaternion with q0, q1, q2 and q3-components.
type Quaternion struct {
	Q0, Q1, Q2, Q3 float64
}

func (t *Quaternion) UnmarshalMTData2Packet(packet MTData2Packet) error {
	var err error
	switch packet.Identifier().Precision {
	case PrecisionFloat32:
		fields := struct {
			Q0, Q1, Q2, Q3 float32
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.Q0 = float64(fields.Q0)
			t.Q1 = float64(fields.Q1)
			t.Q2 = float64(fields.Q2)
			t.Q3 = float64(fields.Q3)
		}
	case PrecisionFP1220:
		fields := struct {
			Q0, Q1, Q2, Q3 FP1220
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.Q0 = fields.Q0.Float64()
			t.Q1 = fields.Q1.Float64()
			t.Q2 = fields.Q2.Float64()
			t.Q3 = fields.Q3.Float64()
		}
	case PrecisionFP1632:
		fields := struct {
			Q0, Q1, Q2, Q3 FP1632
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.Q0 = fields.Q0.Float64()
			t.Q1 = fields.Q1.Float64()
			t.Q2 = fields.Q2.Float64()
			t.Q3 = fields.Q3.Float64()
		}
	case PrecisionFloat64:
		err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, t)
	default:
		err = fmt.Errorf("invalid precision: %v", packet.Identifier().Precision)
	}
	if err != nil {
		return fmt.Errorf("precision %v: %w", packet.Identifier().Precision, err)
	}
	return nil
}

func (t *Quaternion) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(id.Precision.Size()*4, id)
	switch id.Precision {
	case PrecisionFloat32:
		binary.BigEndian.PutUint32(packet.Data(), math.Float32bits(float32(t.Q0)))
		binary.BigEndian.PutUint32(packet.Data()[id.Precision.Size()*1:], math.Float32bits(float32(t.Q1)))
		binary.BigEndian.PutUint32(packet.Data()[id.Precision.Size()*2:], math.Float32bits(float32(t.Q2)))
		binary.BigEndian.PutUint32(packet.Data()[id.Precision.Size()*3:], math.Float32bits(float32(t.Q3)))
	case PrecisionFP1220:
		var q0, q1, q2, q3 FP1220
		q0.FromFloat64(t.Q0)
		q1.FromFloat64(t.Q1)
		q2.FromFloat64(t.Q2)
		q3.FromFloat64(t.Q3)
		q0.toBinary(packet.Data())
		q1.toBinary(packet.Data()[id.Precision.Size():])
		q2.toBinary(packet.Data()[id.Precision.Size()*2:])
		q3.toBinary(packet.Data()[id.Precision.Size()*3:])
	case PrecisionFP1632:
		var q0, q1, q2, q3 FP1632
		q0.FromFloat64(t.Q0)
		q1.FromFloat64(t.Q1)
		q2.FromFloat64(t.Q2)
		q3.FromFloat64(t.Q3)
		q0.toBinary(packet.Data())
		q1.toBinary(packet.Data()[id.Precision.Size():])
		q2.toBinary(packet.Data()[id.Precision.Size()*2:])
		q3.toBinary(packet.Data()[id.Precision.Size()*3:])
	case PrecisionFloat64:
		binary.BigEndian.PutUint64(packet.Data(), math.Float64bits((t.Q0)))
		binary.BigEndian.PutUint64(packet.Data()[id.Precision.Size()*1:], math.Float64bits((t.Q1)))
		binary.BigEndian.PutUint64(packet.Data()[id.Precision.Size()*2:], math.Float64bits((t.Q2)))
		binary.BigEndian.PutUint64(packet.Data()[id.Precision.Size()*3:], math.Float64bits((t.Q3)))
	}
	return packet, nil
}

// DeltaV contains the delta velocity value of the SDI output in m/s.
type DeltaV = VectorXYZ

// Acceleration contains the calibrated acceleration vector in x, y, and z axes in m/s 2 .
type Acceleration = VectorXYZ

// FreeAcceleration contains the free acceleration vector in x, y, and z axes in m/s 2 .
type FreeAcceleration = VectorXYZ

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
type AccelerationHR = VectorXYZ

// DeltaQ contains the delta quaternion value of the SDI output.
type DeltaQ = Quaternion

// RateOfTurn contains the calibrated rate of turn vector in x, y, and z axes in rad/s.
type RateOfTurn = VectorXYZ

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
type RateOfTurnHR = VectorXYZ

// EulerAngles contains the three Euler angles in degrees that represent the orientation of the device.
type EulerAngles = VectorXYZ

// Temperature contains the internal temperature of the sensor in degrees Celsius.
type Temperature = Scalar

// AltitudeEllipsoid contains the altitude of the MTi-G in meters above the WGS-84 Ellipsoid.
type AltitudeEllipsoid = Scalar

// PositionECEF contains the position of the MTi-G in the Earth-Centered, Earth-Fixed (ECEF) coordinate
// system in meters.
//
// Note that position in ECEF cannot be represented in Fixed Point values because of the limited range of fixed point
// representations. Use double or float representation instead.
type PositionECEF = VectorXYZ

// VelocityXYZ contains the X, Y and Z components of the MTi-G velocity in m/s.
type VelocityXYZ = VectorXYZ

// MagneticField contains the magnetic field value in x, y, and z axes in arbitrary units.
//
// Magnetic field is normalized to 1.0 during calibration.
type MagneticField = VectorXYZ

// RotationMatrix contains the rotation matrix (DCM) that represents the orientation of the MT.
type RotationMatrix struct {
	A, B, C, D, E, F, G, H, I float64
}

func (t *RotationMatrix) UnmarshalMTData2Packet(packet MTData2Packet) error {
	var err error
	switch packet.Identifier().Precision {
	case PrecisionFloat32:
		fields := struct {
			A, B, C, D, E, F, G, H, I float32
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.A = float64(fields.A)
			t.B = float64(fields.B)
			t.C = float64(fields.C)
			t.D = float64(fields.D)
			t.E = float64(fields.E)
			t.F = float64(fields.F)
			t.G = float64(fields.G)
			t.H = float64(fields.H)
			t.I = float64(fields.I)
		}
	case PrecisionFP1220:
		fields := struct {
			A, B, C, D, E, F, G, H, I FP1220
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.A = fields.A.Float64()
			t.B = fields.B.Float64()
			t.C = fields.C.Float64()
			t.D = fields.D.Float64()
			t.E = fields.E.Float64()
			t.F = fields.F.Float64()
			t.G = fields.G.Float64()
			t.H = fields.H.Float64()
			t.I = fields.I.Float64()
		}
	case PrecisionFP1632:
		fields := struct {
			A, B, C, D, E, F, G, H, I FP1632
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.A = fields.A.Float64()
			t.B = fields.B.Float64()
			t.C = fields.C.Float64()
			t.D = fields.D.Float64()
			t.E = fields.E.Float64()
			t.F = fields.F.Float64()
			t.G = fields.G.Float64()
			t.H = fields.H.Float64()
			t.I = fields.I.Float64()
		}
	case PrecisionFloat64:
		err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, t)
	default:
		err = fmt.Errorf("invalid precision: %v", packet.Identifier().Precision)
	}
	if err != nil {
		return fmt.Errorf("precision %v: %w", packet.Identifier().Precision, err)
	}
	return nil
}

func (t *RotationMatrix) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(id.Precision.Size()*9, id)
	vals := []float64{t.A, t.B, t.C, t.D, t.E, t.F, t.G, t.H, t.I}
	switch id.Precision {
	case PrecisionFloat32:
		for i, v := range vals {
			binary.BigEndian.PutUint32(packet.Data()[id.Precision.Size()*uint8(i):], math.Float32bits(float32(v)))
		}
	case PrecisionFP1220:
		for i, v := range vals {
			var temp FP1220
			temp.FromFloat64(v)
			temp.toBinary(packet.Data()[id.Precision.Size()*uint8(i):])
		}
	case PrecisionFP1632:
		for i, v := range vals {
			var temp FP1632
			temp.FromFloat64(v)
			temp.toBinary(packet.Data()[id.Precision.Size()*uint8(i):])
		}
	case PrecisionFloat64:
		for i, v := range vals {
			binary.BigEndian.PutUint64(packet.Data()[id.Precision.Size()*uint8(i):], math.Float64bits(v))
		}
	}
	return packet, nil
}

// LatLon contains the latitude and longitude in degrees of the MTi-G position.
type LatLon struct {
	Lat, Lon float64
}

func (t *LatLon) UnmarshalMTData2Packet(packet MTData2Packet) error {
	var err error
	switch packet.Identifier().Precision {
	case PrecisionFloat32:
		fields := struct {
			Lat, Lon float32
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.Lat = float64(fields.Lat)
			t.Lon = float64(fields.Lon)
		}
	case PrecisionFP1220:
		fields := struct {
			Lat, Lon FP1220
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.Lat = fields.Lat.Float64()
			t.Lon = fields.Lon.Float64()
		}
	case PrecisionFP1632:
		fields := struct {
			Lat, Lon FP1632
		}{}
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &fields); err == nil {
			t.Lat = fields.Lat.Float64()
			t.Lon = fields.Lon.Float64()
		}
	case PrecisionFloat64:
		err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, t)
	default:
		err = fmt.Errorf("invalid precision: %v", packet.Identifier().Precision)
	}
	if err != nil {
		return fmt.Errorf("precision %v: %w", packet.Identifier().Precision, err)
	}
	return nil
}

func (t *LatLon) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(id.Precision.Size()*2, id)
	vals := []float64{t.Lat, t.Lon}
	switch id.Precision {
	case PrecisionFloat32:
		for i, v := range vals {
			binary.BigEndian.PutUint32(packet.Data()[id.Precision.Size()*uint8(i):], math.Float32bits(float32(v)))
		}
	case PrecisionFP1220:
		for i, v := range vals {
			var temp FP1220
			temp.FromFloat64(v)
			temp.toBinary(packet.Data()[id.Precision.Size()*uint8(i):])
		}
	case PrecisionFP1632:
		for i, v := range vals {
			var temp FP1632
			temp.FromFloat64(v)
			temp.toBinary(packet.Data()[id.Precision.Size()*uint8(i):])
		}
	case PrecisionFloat64:
		for i, v := range vals {
			binary.BigEndian.PutUint64(packet.Data()[id.Precision.Size()*uint8(i):], math.Float64bits(v))
		}
	}
	return packet, nil
}

// StatusByte contains the 8bit status byte which is equal to bits 0-7 of an MTData2 StatusWord packet.
type StatusByte uint8

func (t *StatusByte) String() string {
	return fmt.Sprintf("%08b", *t)
}

func (t *StatusByte) UnmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, t)
}

func (t *StatusByte) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(1, id)
	packet.Data()[0] = uint8(*t)
	return packet, nil
}

// StatusWord contains the 32bit status word.
//
// 0: Selftest
//
// This flag indicates if the MT passed theself-test according to eMTS.
// For an up-to-date result of the self-test, use the command RunSelftest.
//
// 1: Filter Valid
//
// This flag indicates if input into the orientation filter is reliable and / or complete.
// If for example the measurement range of internal sensors is exceeded, orientation output cannot be reliably
// estimated and the filter flag will drop to 0.
//
// For the MTi-G, the filter flag will also become invalid if the GPS status remains invalid for an extended period.
//
// 2: GNSS fix
//
// This flag indicates if the GNSS unit has a proper fix. The flag is only available in MTi-G units.
//
// 3-4: NoRotationUpdate Status
//
// This flag indicates the status of the no rotation update procedure in the filter after the SetNoRotation message
// has been sent.
//
//  11: Running with no rotation assumption
//  10: Rotation detected, no gyro bias estimation (sticky)
//  00: Estimation complete, no errors
//
// 5 Representative Motion
//
// (RepMo) Indicates if the MTi is in In-run Compass Calibration Representative Mode
//
// 6-7: Reserved Reserved for future use
//
// 8-19: Clip flags
//
// Indicates out of range values on sensors.
//
//  8:     Clipflag Acc X
//  9:     Clipflag Acc Y
//  10:    Clipflag Acc Z
//  11:    Clipflag Gyr X
//  12:    Clipflag Gyr Y
//  13:    Clipflag Gyr Z
//  14:    Clipflag Mag X
//  15:    Clipflag Mag Y
//  16:    Clipflag Mag Z
//  17-18: Reserved Reserved for future use
//  19:    Clipping Indication (indicates that one or more sensors are out of range)
//
// 20: Reserved Reserved for future use
//
// 21: SyncIn Marker
//
// When a SyncIn is detected, this bit will rise to 1
//
// 22: SyncOut Marker
//
// When SyncOut is active this bit will rise to 1
//
// 23-25: Filter Mode
//
// Indicates Filter Mode, currently only available for the MTi-G-710 and MTi-7:
//
//  000: Without GNSS (filter profile is in VRU mode)
//  001: Coasting mode (GNSS has been lost <60 sec ago)
//  011: With GNSS (default mode of MTi-G-710)
//
// 26-31: Reserved
//
// Reserved for future use.
type StatusWord uint32

// String returns a string representation of the status word.
func (t *StatusWord) String() string {
	return fmt.Sprintf("%032b", *t)
}

func (t *StatusWord) UnmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, t)
}

func (t *StatusWord) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(4, id)
	binary.BigEndian.PutUint32(packet.Data(), uint32(*t))
	return packet, nil
}

// UTCTime contains the timestamp expressed as the UTC time.
type UTCTime struct {
	Ns                               uint32
	Year                             uint16
	Month, Day, Hour, Minute, Second uint8
	Valid                            UTCValidity
}

// String returns the UTC time on RFC3339 (including nanoseconds) format.
func (u *UTCTime) String() string {
	return u.Time().Format(time.RFC3339Nano)
}

func (u *UTCTime) UnmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, u)
}

func (u *UTCTime) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(12, id)
	binary.BigEndian.PutUint32(packet.Data(), u.Ns)
	binary.BigEndian.PutUint16(packet.Data()[4:], u.Year)
	packet.Data()[6] = u.Month
	packet.Data()[7] = u.Day
	packet.Data()[8] = u.Hour
	packet.Data()[9] = u.Minute
	packet.Data()[10] = u.Second
	packet.Data()[11] = uint8(u.Valid)
	return packet, nil
}

// Time returns the native Go representation of the UTC time.
func (u *UTCTime) Time() time.Time {
	return time.Date(
		int(u.Year),
		time.Month(u.Month),
		int(u.Day),
		int(u.Hour),
		int(u.Minute),
		int(u.Second),
		int(u.Ns),
		time.UTC,
	)
}

func (u *UTCTime) UnmarshalTime(ts time.Time) {
	t := ts.UTC()
	u.Year = uint16(t.Year())
	u.Month = uint8(t.Month())
	u.Day = uint8(t.Day())
	u.Hour = uint8(t.Hour())
	u.Minute = uint8(t.Minute())
	u.Second = uint8(t.Second())
	u.Ns = uint32(t.Nanosecond())
}

// PacketCounter contains the packet counter.
//
// This counter is incremented with every generated MTData2 message.
type PacketCounter uint16

// String returns a string representation of the packet counter.
func (p *PacketCounter) String() string {
	return strconv.Itoa(int(*p))
}

func (p *PacketCounter) UnmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, p)
}

func (p *PacketCounter) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(2, id)
	binary.BigEndian.PutUint16(packet.Data(), uint16(*p))
	return packet, nil
}

// SampleTimeFine contains the sample time of an output expressed in 10kHz ticks.
//
// When there is no GNSS-fix in the MTi-G-710, this value is arbitrary for GNSS messages.
type SampleTimeFine uint32

// String returns a string representation of the sample time.
func (s *SampleTimeFine) String() string {
	return strconv.Itoa(int(*s))
}

func (s *SampleTimeFine) UnmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, s)
}

func (s *SampleTimeFine) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(4, id)
	binary.BigEndian.PutUint32(packet.Data(), uint32(*s))
	return packet, nil
}

// SampleTimeCoarse contains the sample time of an output expressed in seconds.
//
// When there is no GNSS-fix in the MTi-G-710, this value is arbitrary for GNSS messages.
type SampleTimeCoarse uint32

// String returns a string representation of the sample time.
func (s *SampleTimeCoarse) String() string {
	return strconv.Itoa(int(*s))
}

func (s *SampleTimeCoarse) UnmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, s)
}

func (s *SampleTimeCoarse) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(4, id)
	binary.BigEndian.PutUint32(packet.Data(), uint32(*s))
	return packet, nil
}

// BaroPressure contains the pressure as measured by the internal barometer expressed in Pascal.
type BaroPressure uint32

// String returns a string representation of the value.
func (b *BaroPressure) String() string {
	return strconv.Itoa(int(*b))
}

func (b *BaroPressure) UnmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, b)
}

func (b *BaroPressure) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(4, id)
	binary.BigEndian.PutUint32(packet.Data(), uint32(*b))
	return packet, nil
}

// GNSSPVTData contains the current GNSS position, velocity and time data.
type GNSSPVTData struct {
	// ITOW is the GPS time of week.
	//
	//  Unit: ms
	ITOW uint32

	// Year (UTC).
	//
	//  Unit: y
	Year uint16

	// Month (UTC).
	//
	//  Unit: m
	Month uint8

	// Day of the month (UTC).
	//
	//  Unit: d
	Day uint8

	// Hour of the day 0..23 (UTC).
	//
	//  Unit: h
	Hour uint8

	// Minute of hour 0..59 (UTC).
	//
	//  Unit: min
	Min uint8

	// Seconds of minute 0..60 (UTC).
	//
	//  Unit: s
	Sec uint8

	// Valid is the validity flags.
	//
	//  bit (0) = UTC Date is valid
	//  bit (1) = UTC Time of Day is valid
	//  bit (2) = UTC Time of Day has been fully resolved (i.e. no seconds uncertainty)
	Valid UTCValidity

	// TAcc is the time accuracy estimate (UTC).
	//
	//  Unit: ns
	TAcc uint32

	// Nano is the fraction of second -1e-9 .. 1e-9.
	//
	//  Unit: ns
	Nano int32

	// FixType is the GNSS fix type.
	FixType FixType

	// Flags are the fix status flags
	//
	//  bit (0) = Valid fix (within DOP and accuracy masks)
	//  bit (1) = Differential corrections are applied
	//  bit (2) = Reserved
	//  bit (3) = Reserved
	//  bit (4) = Reserved
	//  bit (5) = Heading of vehicle is valid
	Flags uint8

	// NumSV is the number of satellites used in navigation solution.
	NumSV uint8

	// Reserved1 is reserved for future use.
	Reserved1 uint8

	// Lon is the position longitude.
	//
	//  Scale: 1e-7
	//  Unit: deg
	Lon int32

	// Lat is the position latitude.
	//
	//  Scale: 1e-7
	//  Unit: deg
	Lat int32

	// Height above ellipsoid.
	//
	//  Unit: mm
	Height int32

	// HMSL is the height above mean sea level.
	//
	//  Unit: mm
	HMSL int32

	// HAcc is the horizontal accuracy estimate.
	//
	//  Unit: mm
	HAcc uint32

	// VAcc is the vertical accuracy estimate.
	//
	//  Unit: mm
	VAcc uint32

	// VelN is the NED north velocity.
	//
	//  Unit: mm/s
	VelN int32

	// VelE is the NED east velocity.
	//
	//  Unit: mm/s
	VelE int32

	// VelD is the NED down velocity.
	//
	//  Unit: mm/s
	VelD int32

	// GSpeed is the 2D ground speed.
	//
	//  Unit: mm/s
	GSpeed int32

	// HeadMot is the 2D heading of motion.
	//
	//  Scale: 1e-5
	//  Unit: deg
	HeadMot int32

	// SAcc is the speed accuracy estimate.
	//
	//  Unit: mm/s
	SAcc uint32

	// HeadAcc is the heading accuracy estimate (both motion and vehicle).
	//
	//  Unit: deg
	HeadAcc uint32

	// HeadVeh is the 2D heading of the vehicle.
	//
	//  Scale: 1e-5
	//  Unit: deg
	HeadVeh uint32

	// GDOP is the Geometric DOP.
	//
	//  Scale: 0.01
	GDOP uint16

	// PDOP is the Position DOP.
	//
	//  Scale: 0.01
	PDOP uint16

	// PDOP is the Time DOP.
	//
	//  Scale: 0.01
	TDOP uint16

	// VDOP is the Vertical DOP.
	//
	//  Scale: 0.01
	VDOP uint16

	// HDOP is the Horizontal DOP.
	//
	//  Scale: 0.01
	HDOP uint16

	// NDOP is the Northing DOP.
	//
	//  Scale: 0.01
	NDOP uint16

	// EDOP is the Easting DOP.
	//
	//  Scale: 0.01
	EDOP uint16
}

func (g *GNSSPVTData) Time() time.Time {
	return time.Date(
		int(g.Year),
		time.Month(g.Month),
		int(g.Day),
		int(g.Hour),
		int(g.Min),
		int(g.Sec),
		int(g.Nano),
		time.UTC,
	)
}

func (g *GNSSPVTData) UnmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, g)
}

func (g *GNSSPVTData) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(94, id)
	binary.BigEndian.PutUint32(packet.Data(), g.ITOW)
	binary.BigEndian.PutUint16(packet.Data()[4:], g.Year)
	packet.Data()[6] = g.Month
	packet.Data()[7] = g.Day
	packet.Data()[8] = g.Hour
	packet.Data()[9] = g.Min
	packet.Data()[10] = g.Sec
	packet.Data()[11] = uint8(g.Valid)
	binary.BigEndian.PutUint32(packet.Data()[12:], g.TAcc)
	binary.BigEndian.PutUint32(packet.Data()[16:], uint32(g.Nano))
	packet.Data()[20] = uint8(g.FixType)
	packet.Data()[21] = g.Flags
	packet.Data()[22] = g.NumSV
	packet.Data()[23] = g.Reserved1
	binary.BigEndian.PutUint32(packet.Data()[24:], uint32(g.Lon))
	binary.BigEndian.PutUint32(packet.Data()[28:], uint32(g.Lat))
	binary.BigEndian.PutUint32(packet.Data()[32:], uint32(g.Height))
	binary.BigEndian.PutUint32(packet.Data()[36:], uint32(g.HMSL))
	binary.BigEndian.PutUint32(packet.Data()[40:], g.HAcc)
	binary.BigEndian.PutUint32(packet.Data()[44:], g.VAcc)
	binary.BigEndian.PutUint32(packet.Data()[48:], uint32(g.VelN))
	binary.BigEndian.PutUint32(packet.Data()[52:], uint32(g.VelE))
	binary.BigEndian.PutUint32(packet.Data()[56:], uint32(g.VelD))
	binary.BigEndian.PutUint32(packet.Data()[60:], uint32(g.GSpeed))
	binary.BigEndian.PutUint32(packet.Data()[64:], uint32(g.HeadMot))
	binary.BigEndian.PutUint32(packet.Data()[68:], g.SAcc)
	binary.BigEndian.PutUint32(packet.Data()[72:], g.HeadAcc)
	binary.BigEndian.PutUint32(packet.Data()[76:], g.HeadVeh)
	binary.BigEndian.PutUint16(packet.Data()[80:], g.GDOP)
	binary.BigEndian.PutUint16(packet.Data()[82:], g.PDOP)
	binary.BigEndian.PutUint16(packet.Data()[84:], g.TDOP)
	binary.BigEndian.PutUint16(packet.Data()[86:], g.VDOP)
	binary.BigEndian.PutUint16(packet.Data()[88:], g.HDOP)
	binary.BigEndian.PutUint16(packet.Data()[90:], g.NDOP)
	binary.BigEndian.PutUint16(packet.Data()[92:], g.EDOP)
	return packet, nil
}

// GNSSSatInfo contains info on the currently used GNSS satellites.
type GNSSSatInfo struct {
	// ITOW is the GPS time of week.
	//
	//  Unit: ms
	ITOW uint32

	// NumSVS is the number of satellites.
	NumSVS uint8

	// Res1 is reserved for future use.
	Res1 uint8

	// Res2 is reserved for future use.
	Res2 uint8

	// Res3 is reserved for future use.
	Res3 uint8
}

func (g *GNSSSatInfo) UnmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, g)
}

func (g *GNSSSatInfo) MarshalMTData2Packet(id DataIdentifier) (MTData2Packet, error) {
	packet := NewMTData2Package(8, id)
	packet.SetIdentifier(id)
	binary.BigEndian.PutUint32(packet.Data(), g.ITOW)
	packet.Data()[4] = g.NumSVS
	packet.Data()[5] = g.Res1
	packet.Data()[6] = g.Res2
	packet.Data()[7] = g.Res3
	return packet, nil
}

type GNSSSat struct {
	// GNSSID is the GNSS identifier.
	//
	//  0 = GPS
	//  1 = SBAS
	//  2 = Galileo
	//  3 = BeiDou
	//  4 = IMES
	//  5 = QZSS
	//  6 = GLONASS
	GNSSID uint8

	// SVID is the satellite identifier.
	SVID uint8

	// CNO is the carrier to noise ratio (signal strength).
	//
	//  Unit: dBHz
	CNO uint8

	// Flags contains the satellite flags.
	//
	//  bit (0..2) = signal quality indicator
	//   0 = no signal
	//   1 = searching signal
	//   2 = signal acquired
	//   3 = signal detected but unusable
	//   4 = code locked and time synchronized
	//   5, 6, 7 = code & carrier locked; time synchronized
	//  bit (3) = SV is being used for navigation
	//  bit (4..5) = SV health flag
	//    0 = unknown
	//    1 = healthy
	//    2 = unhealthy
	//  bit (6) = differential correction data is available
	//  bit (7) = reserved
	Flags uint8
}
