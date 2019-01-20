package xsens

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Scalar contains a single scalar value.
type Scalar float64

// String returns a string representation of the scalar.
func (s *Scalar) String() string {
	return strconv.FormatFloat(float64(*s), 'f', -1, 64)
}

func (s *Scalar) unmarshalMTData2Packet(packet MTData2Packet) error {
	var err error
	switch packet.Identifier().Precision {
	case PrecisionFloat32:
		var value float32
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &value); err == nil {
			*s = Scalar(value)
		}
	case PrecisionFP1220:
		var value FP1220
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &value); err == nil {
			*s = Scalar(value.Float64())
		}
	case PrecisionFP1632:
		var value FP1632
		if err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, &value); err == nil {
			*s = Scalar(value.Float64())
		}
	case PrecisionFloat64:
		err = binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, s)
	default:
		err = errors.New("invalid output format")
	}
	return errors.Wrapf(err, "output format: %v", packet.Identifier().Precision)
}

// VectorXYZ contains a vector with x, y and z-components.
type VectorXYZ struct {
	X, Y, Z float64
}

func (t *VectorXYZ) unmarshalMTData2Packet(packet MTData2Packet) error {
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
		err = errors.New("invalid output format")
	}
	return errors.Wrapf(err, "output format: %v", packet.Identifier().Precision)
}

// Quaternion contains a quaternion with q0, q1, q2 and q3-components.
type Quaternion struct {
	Q0, Q1, Q2, Q3 float64
}

func (t *Quaternion) unmarshalMTData2Packet(packet MTData2Packet) error {
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
		err = errors.New("invalid output format")
	}
	return errors.Wrapf(err, "output format: %v", packet.Identifier().Precision)
}

// RotationMatrix contains a rotation matrix (DCM).
type RotationMatrix struct {
	A, B, C, D, E, F, G, H, I float64
}

func (t *RotationMatrix) unmarshalMTData2Packet(packet MTData2Packet) error {
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
		err = errors.New("invalid output format")
	}
	return errors.Wrapf(err, "output format: %v", packet.Identifier().Precision)
}

// LatLon contains latitude and longitude components.
type LatLon struct {
	Lat, Lon float64
}

func (t *LatLon) unmarshalMTData2Packet(packet MTData2Packet) error {
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
		err = errors.New("invalid output format")
	}
	return errors.Wrapf(err, "output format: %v", packet.Identifier().Precision)
}

// StatusByte contains the 8bit status byte which is equal to bits 0-7 of an MTData2 StatusWord packet.
type StatusByte uint8

func (t *StatusByte) String() string {
	return fmt.Sprintf("%08b", *t)
}

func (t *StatusByte) unmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, t)
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
// Reserved for future use
type StatusWord uint32

// String returns a string representation of the status word.
func (t *StatusWord) String() string {
	return fmt.Sprintf("%032b", *t)
}

func (t *StatusWord) unmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, t)
}

// UTCTime contains the timestamp expressed as the UTC time.
type UTCTime struct {
	Ns                               uint32
	Year                             uint16
	Month, Day, Hour, Minute, Second uint8
	Flags                            uint8
}

// String returns the UTC time on RFC3339 (including nanoseconds) format.
func (u *UTCTime) String() string {
	return u.Time().Format(time.RFC3339Nano)
}

func (u *UTCTime) unmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, u)
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
		time.UTC)
}

// IsValid returns true if the valid flag is set on the UTC time flags.
func (u *UTCTime) IsValid() bool {
	return u.Flags&0x04 > 0
}

// PacketCounter contains the packet counter.
//
// This counter is incremented with every generated MTData2 message.
type PacketCounter uint16

// String returns a string representation of the packet counter.
func (p *PacketCounter) String() string {
	return strconv.Itoa(int(*p))
}

func (p *PacketCounter) unmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, p)
}

// SampleTimeFine contains the sample time of an output expressed in 10kHz ticks.
//
// When there is no GNSS-fix in the MTi-G-710, this value is arbitrary for GNSS messages.
type SampleTimeFine uint32

// String returns a string representation of the sample time.
func (s *SampleTimeFine) String() string {
	return strconv.Itoa(int(*s))
}

func (s *SampleTimeFine) unmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, s)
}

// SampleTimeCoarse contains the sample time of an output expressed in seconds.
//
// When there is no GNSS-fix in the MTi-G-710, this value is arbitrary for GNSS messages.
type SampleTimeCoarse uint32

// String returns a string representation of the sample time.
func (s *SampleTimeCoarse) String() string {
	return strconv.Itoa(int(*s))
}

func (s *SampleTimeCoarse) unmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, s)
}

// BaroPressure contains the pressure as measured by the internal barometer expressed in Pascal.
type BaroPressure uint32

// String returns a string representation of the value.
func (b *BaroPressure) String() string {
	return strconv.Itoa(int(*b))
}

func (b *BaroPressure) unmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, b)
}

// GNSSPVTData contains GNSS measurement data.
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
	Valid uint8

	// TAcc is the time accuracy estimate (UTC).
	//
	//  Unit: ns
	TAcc uint32

	// Nano is the fraction of second -1e-9 .. 1e-9.
	//
	//  Unit: ns
	Nano int32

	// FixType is the GNSS fix type (range 0..5).
	//
	//  0x00 = No fix
	//  0x01 = Dead reckoning only
	//  0x02 = 2D fix
	//  0x03 = 3D fix
	//  0x04 = GNSS + dead reckoning combined
	//  0x05 = Time only fix
	FixType uint8

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

func (g *GNSSPVTData) unmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, g)
}

// GNSSSatInfo contains GNSS satellite info.
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

func (g *GNSSSatInfo) unmarshalMTData2Packet(packet MTData2Packet) error {
	return binary.Read(bytes.NewReader(packet.Data()), binary.BigEndian, g)
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
	//   4 = code locked and time synchronised
	//   5, 6, 7 = code & carrier locked; time synchronised
	//  bit (3) = SV is being used for navigation
	//  bit (4..5) = SV health flag
	//    0 = unknown
	//    1 = healthy
	//    2 = unhealthy
	//  bit (6) = differential correction data is available
	//  bit (7) = reserved
	Flags uint8
}
