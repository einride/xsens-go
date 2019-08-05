package xsens

import (
	"bufio"
	"context"
	"io"
	"time"

	"golang.org/x/xerrors"
)

// SerialPort is an interface for a serial used by the Client for communication with an Xsens device.
type SerialPort interface {
	io.ReadWriteCloser
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
}

// Client for communicating with an Xsens device.
type Client struct {
	p               SerialPort
	sc              *bufio.Scanner
	message         Message
	mtData2         MTData2
	mtData2Packet   MTData2Packet
	nextPacketIndex int
	// measurement data
	packetCounter     PacketCounter
	statusByte        StatusByte
	statusWord        StatusWord
	sampleTimeFine    SampleTimeFine
	sampleTimeCoarse  SampleTimeCoarse
	baroPressure      BaroPressure
	utcTime           UTCTime
	deltaV            VectorXYZ
	acceleration      VectorXYZ
	freeAcceleration  VectorXYZ
	accelerationHR    VectorXYZ
	deltaQ            Quaternion
	rateOfTurn        VectorXYZ
	rateOfTurnHR      VectorXYZ
	quaternion        Quaternion
	eulerAngles       VectorXYZ
	rotationMatrix    RotationMatrix
	temperature       Scalar
	altitudeEllipsoid Scalar
	positionECEF      VectorXYZ
	latLon            LatLon
	velocityXYZ       VectorXYZ
	magneticField     VectorXYZ
	gnssPVTData       GNSSPVTData
	gnssSatInfo       GNSSSatInfo
}

// NewClient returns a new client using the provided serial port for communication.
func NewClient(p SerialPort) *Client {
	sc := bufio.NewScanner(p)
	sc.Split(ScanMessages)
	return &Client{p: p, sc: sc}
}

// Close the client's serial port.
func (c *Client) Close() error {
	if err := c.p.Close(); err != nil {
		return xerrors.Errorf("xsens client: close: %w", err)
	}
	return nil
}

// Receive an Xsens message.
//
// Clears state related to a previously received message, such as scanned measurement data.
func (c *Client) Receive(ctx context.Context) error {
	// clear previous received message state
	c.message = nil
	c.mtData2 = nil
	c.mtData2Packet = nil
	c.nextPacketIndex = 0
	// receive new message
	if !c.sc.Scan() {
		if c.sc.Err() == nil {
			return xerrors.Errorf("xsens client: receive: %w", io.EOF)
		}
		return xerrors.Errorf("xsens client: receive: %w", c.sc.Err())
	}
	c.message = c.sc.Bytes()
	if c.message.Identifier() == MessageIdentifierMTData2 {
		c.mtData2 = c.message.Data()
	}
	return nil
}

// RawMessage returns the raw bytes of the last received message.
func (c *Client) RawMessage() []byte {
	return c.sc.Bytes()
}

// GoToConfig puts the Xsens device in config mode.
func (c *Client) GoToConfig(ctx context.Context) error {
	if err := c.send(ctx, NewMessage(MessageIdentifierGotoConfig, nil)); err != nil {
		return xerrors.Errorf("xsens client: go to config: %w", err)
	}
	if err := c.receiveUntil(ctx, MessageIdentifierGotoConfigAck); err != nil {
		return xerrors.Errorf("xsens client: go to config: %w", err)
	}
	return nil
}

// SetOutputConfiguration sets the Xsens device output configuration.
func (c *Client) SetOutputConfiguration(ctx context.Context, configuration OutputConfiguration) error {
	data, err := configuration.Marshal()
	if err != nil {
		return xerrors.Errorf("xsens client: set output configuration: %w", err)
	}
	if err := c.send(ctx, NewMessage(MessageIdentifierSetOutputConfiguration, data)); err != nil {
		return xerrors.Errorf("xsens client: set output configuration: %w", err)
	}
	if err := c.receiveUntil(ctx, MessageIdentifierSetOutputConfigurationAck); err != nil {
		return xerrors.Errorf("xsens client: set output configuration: %w", err)
	}
	return nil
}

// GoToMeasurement puts the Xsens device in measurement mode.
func (c *Client) GoToMeasurement(ctx context.Context) error {
	if err := c.send(ctx, NewMessage(MessageIdentifierGotoMeasurement, nil)); err != nil {
		return xerrors.Errorf("xsens client: go to measurement: %w", err)
	}
	if err := c.receiveUntil(ctx, MessageIdentifierMTData2); err != nil {
		return xerrors.Errorf("xsens client: go to config: %w", err)
	}
	return nil
}

// MessageIdentifier returns the message identifier of the last received message.
func (c *Client) MessageIdentifier() MessageIdentifier {
	return c.message.Identifier()
}

// ScanMeasurementData advances to the next measurement data packet, when the current message contains measurement data.
func (c *Client) ScanMeasurementData() bool {
	if c.message.Identifier() != MessageIdentifierMTData2 {
		return false
	}
	packet, err := c.mtData2.PacketAt(c.nextPacketIndex)
	if err != nil {
		return false
	}
	c.nextPacketIndex += len(packet)
	c.mtData2Packet = packet
	switch packet.Identifier().DataType {
	case DataTypeDeltaV:
		err = c.deltaV.unmarshalMTData2Packet(packet)
	case DataTypeAcceleration:
		err = c.acceleration.unmarshalMTData2Packet(packet)
	case DataTypeFreeAcceleration:
		err = c.freeAcceleration.unmarshalMTData2Packet(packet)
	case DataTypeAccelerationHR:
		err = c.accelerationHR.unmarshalMTData2Packet(packet)
	case DataTypeDeltaQ:
		err = c.deltaQ.unmarshalMTData2Packet(packet)
	case DataTypeRateOfTurn:
		err = c.rateOfTurn.unmarshalMTData2Packet(packet)
	case DataTypeRateOfTurnHR:
		err = c.rateOfTurnHR.unmarshalMTData2Packet(packet)
	case DataTypeQuaternion:
		err = c.quaternion.unmarshalMTData2Packet(packet)
	case DataTypeEulerAngles:
		err = c.eulerAngles.unmarshalMTData2Packet(packet)
	case DataTypeRotationMatrix:
		err = c.rotationMatrix.unmarshalMTData2Packet(packet)
	case DataTypeTemperature:
		err = c.temperature.unmarshalMTData2Packet(packet)
	case DataTypeAltitudeEllipsoid:
		err = c.altitudeEllipsoid.unmarshalMTData2Packet(packet)
	case DataTypePositionECEF:
		err = c.positionECEF.unmarshalMTData2Packet(packet)
	case DataTypeLatLon:
		err = c.latLon.unmarshalMTData2Packet(packet)
	case DataTypeVelocityXYZ:
		err = c.velocityXYZ.unmarshalMTData2Packet(packet)
	case DataTypeStatusByte:
		err = c.statusByte.unmarshalMTData2Packet(packet)
	case DataTypeStatusWord:
		err = c.statusWord.unmarshalMTData2Packet(packet)
	case DataTypeUTCTime:
		err = c.utcTime.unmarshalMTData2Packet(packet)
	case DataTypePacketCounter:
		err = c.packetCounter.unmarshalMTData2Packet(packet)
	case DataTypeSampleTimeFine:
		err = c.sampleTimeFine.unmarshalMTData2Packet(packet)
	case DataTypeSampleTimeCoarse:
		err = c.sampleTimeCoarse.unmarshalMTData2Packet(packet)
	case DataTypeBaroPressure:
		err = c.baroPressure.unmarshalMTData2Packet(packet)
	case DataTypeMagneticField:
		err = c.magneticField.unmarshalMTData2Packet(packet)
	case DataTypeGNSSPVTData:
		err = c.gnssPVTData.unmarshalMTData2Packet(packet)
	case DataTypeGNSSSatInfo:
		err = c.gnssSatInfo.unmarshalMTData2Packet(packet)
	}
	// TODO: Improve this API after removing MessageScanner
	return err == nil
}

// RawPacket returns the raw bytes of the current measurement data packet.
func (c *Client) RawPacket() []byte {
	return c.mtData2Packet
}

// DataType returns the data type of the current scanned measurement data packet.
func (c *Client) DataType() DataType {
	return c.mtData2Packet.Identifier().DataType
}

func (c *Client) PacketCounter() *PacketCounter {
	return &c.packetCounter
}

func (c *Client) SampleTimeFine() *SampleTimeFine {
	return &c.sampleTimeFine
}

func (c *Client) UTCTime() *UTCTime {
	return &c.utcTime
}

func (c *Client) StatusByte() *StatusByte {
	return &c.statusByte
}

func (c *Client) StatusWord() *StatusWord {
	return &c.statusWord
}

func (c *Client) EulerAngles() *VectorXYZ {
	return &c.eulerAngles
}

func (c *Client) Acceleration() *VectorXYZ {
	return &c.acceleration
}

func (c *Client) DeltaV() *VectorXYZ {
	return &c.deltaV
}

func (c *Client) RateOfTurn() *VectorXYZ {
	return &c.rateOfTurn
}

func (c *Client) DeltaQ() *Quaternion {
	return &c.deltaQ
}

func (c *Client) LatLon() *LatLon {
	return &c.latLon
}

func (c *Client) AltitudeEllipsoid() *Scalar {
	return &c.altitudeEllipsoid
}

func (c *Client) VelocityXYZ() *VectorXYZ {
	return &c.velocityXYZ
}

func (c *Client) GNSSPVTData() *GNSSPVTData {
	return &c.gnssPVTData
}

func (c *Client) SampleTimeCoarse() *SampleTimeCoarse {
	return &c.sampleTimeCoarse
}

func (c *Client) BaroPressure() *BaroPressure {
	return &c.baroPressure
}

func (c *Client) Temperature() *Scalar {
	return &c.temperature
}

func (c *Client) MagneticField() *VectorXYZ {
	return &c.magneticField
}

func (c *Client) RotationMatrix() *RotationMatrix {
	return &c.rotationMatrix
}

func (c *Client) FreeAcceleration() *VectorXYZ {
	return &c.freeAcceleration
}

func (c *Client) Quaternion() *Quaternion {
	return &c.quaternion
}

func (c *Client) GNSSSatInfo() *GNSSSatInfo {
	return &c.gnssSatInfo
}

func (c *Client) PositionECEF() *VectorXYZ {
	return &c.positionECEF
}

func (c *Client) send(ctx context.Context, message Message) error {
	deadline, _ := ctx.Deadline()
	if err := c.p.SetWriteDeadline(deadline); err != nil {
		return xerrors.Errorf("send %v: %w", message.Identifier(), err)
	}
	if _, err := c.p.Write(message); err != nil {
		return xerrors.Errorf("send %v: %w", message.Identifier(), err)
	}
	return nil
}

func (c *Client) receiveUntil(ctx context.Context, until MessageIdentifier) error {
	deadline, _ := ctx.Deadline()
	if err := c.p.SetReadDeadline(deadline); err != nil {
		return xerrors.Errorf("receive until %v: %w", until, err)
	}
	for {
		if err := c.Receive(ctx); err != nil {
			return xerrors.Errorf("receive until %v: %w", until, err)
		}
		if c.MessageIdentifier() != until {
			continue
		}
		return nil
	}
}