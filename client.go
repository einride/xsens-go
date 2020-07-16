package xsens

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"
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
	deltaV            DeltaV
	acceleration      Acceleration
	freeAcceleration  FreeAcceleration
	accelerationHR    AccelerationHR
	deltaQ            DeltaQ
	rateOfTurn        RateOfTurn
	rateOfTurnHR      RateOfTurnHR
	quaternion        Quaternion
	eulerAngles       EulerAngles
	rotationMatrix    RotationMatrix
	temperature       Temperature
	altitudeEllipsoid AltitudeEllipsoid
	positionECEF      PositionECEF
	latLon            LatLon
	velocityXYZ       VelocityXYZ
	magneticField     MagneticField
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
		return fmt.Errorf("xsens client: close: %w", err)
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
	deadline, _ := ctx.Deadline()
	if err := c.p.SetReadDeadline(deadline); err != nil {
		return fmt.Errorf("xsens client: receive: %w", err)
	}
	if !c.sc.Scan() {
		if c.sc.Err() == nil {
			return fmt.Errorf("xsens client: receive: %w", io.EOF)
		}
		return fmt.Errorf("xsens client: receive: %w", c.sc.Err())
	}
	c.message = c.sc.Bytes()
	if err := c.message.Validate(); err != nil {
		return fmt.Errorf("xsens client: receive: %w", c.sc.Err())
	}
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
		return fmt.Errorf("xsens client: go to config: %w", err)
	}
	if err := c.receiveUntil(ctx, MessageIdentifierGotoConfigAck); err != nil {
		return fmt.Errorf("xsens client: go to config: %w", err)
	}
	return nil
}

// SetOutputConfiguration sets the Xsens device output configuration.
func (c *Client) SetOutputConfiguration(ctx context.Context, configuration OutputConfiguration) error {
	data, err := configuration.Marshal()
	if err != nil {
		return fmt.Errorf("xsens client: set output configuration: %w", err)
	}
	if err := c.send(ctx, NewMessage(MessageIdentifierSetOutputConfiguration, data)); err != nil {
		return fmt.Errorf("xsens client: set output configuration: %w", err)
	}
	if err := c.receiveUntil(ctx, MessageIdentifierSetOutputConfigurationAck); err != nil {
		return fmt.Errorf("xsens client: set output configuration: %w", err)
	}
	return nil
}

// GetOutputConfiguration returns the Xsens output configuration.
func (c *Client) GetOutputConfiguration(ctx context.Context) (OutputConfiguration, error) {
	if err := c.send(ctx, NewMessage(MessageIdentifierReqOutputConfiguration, nil)); err != nil {
		return nil, fmt.Errorf("xsens client: get output configuration: %w", err)
	}
	if err := c.receiveUntil(ctx, MessageIdentifierReqOutputConfigurationAck); err != nil {
		return nil, fmt.Errorf("xsens client: get output configuration: %w", err)
	}
	var result OutputConfiguration
	if err := result.Unmarshal(c.message.Data()); err != nil {
		return nil, fmt.Errorf("xsens client: get output configuration: %w", err)
	}
	return result, nil
}

// GoToMeasurement puts the Xsens device in measurement mode.
func (c *Client) GoToMeasurement(ctx context.Context) error {
	if err := c.send(ctx, NewMessage(MessageIdentifierGotoMeasurement, nil)); err != nil {
		return fmt.Errorf("xsens client: go to measurement: %w", err)
	}
	if err := c.receiveUntil(ctx, MessageIdentifierMTData2); err != nil {
		return fmt.Errorf("xsens client: go to config: %w", err)
	}
	return nil
}

// MessageIdentifier returns the message identifier of the last received message.
func (c *Client) MessageIdentifier() MessageIdentifier {
	return c.message.Identifier()
}

// MeasurementData returns the last scanned measurement data.
func (c *Client) MeasurementData() MeasurementData {
	switch c.mtData2Packet.Identifier().DataType {
	case DataTypeDeltaV:
		return &c.deltaV
	case DataTypeAcceleration:
		return &c.acceleration
	case DataTypeFreeAcceleration:
		return &c.freeAcceleration
	case DataTypeAccelerationHR:
		return &c.accelerationHR
	case DataTypeDeltaQ:
		return &c.deltaQ
	case DataTypeRateOfTurn:
		return &c.rateOfTurn
	case DataTypeRateOfTurnHR:
		return &c.rateOfTurnHR
	case DataTypeQuaternion:
		return &c.quaternion
	case DataTypeEulerAngles:
		return &c.eulerAngles
	case DataTypeRotationMatrix:
		return &c.rotationMatrix
	case DataTypeTemperature:
		return &c.temperature
	case DataTypeAltitudeEllipsoid:
		return &c.altitudeEllipsoid
	case DataTypePositionECEF:
		return &c.positionECEF
	case DataTypeLatLon:
		return &c.latLon
	case DataTypeVelocityXYZ:
		return &c.velocityXYZ
	case DataTypeStatusByte:
		return &c.statusByte
	case DataTypeStatusWord:
		return &c.statusWord
	case DataTypeUTCTime:
		return &c.utcTime
	case DataTypePacketCounter:
		return &c.packetCounter
	case DataTypeSampleTimeFine:
		return &c.sampleTimeFine
	case DataTypeSampleTimeCoarse:
		return &c.sampleTimeCoarse
	case DataTypeBaroPressure:
		return &c.baroPressure
	case DataTypeMagneticField:
		return &c.magneticField
	case DataTypeGNSSPVTData:
		return &c.gnssPVTData
	case DataTypeGNSSSatInfo:
		return &c.gnssSatInfo
	default:
		return nil
	}
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
	// TODO: Improve this API after removing MessageScanner
	data := c.MeasurementData()
	if data == nil {
		return false
	}
	if err := data.unmarshalMTData2Packet(c.mtData2Packet); err != nil {
		return false
	}
	return true
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

func (c *Client) EulerAngles() *EulerAngles {
	return &c.eulerAngles
}

func (c *Client) Acceleration() *Acceleration {
	return &c.acceleration
}

func (c *Client) DeltaV() *DeltaV {
	return &c.deltaV
}

func (c *Client) RateOfTurn() *RateOfTurn {
	return &c.rateOfTurn
}

func (c *Client) DeltaQ() *DeltaQ {
	return &c.deltaQ
}

func (c *Client) LatLon() *LatLon {
	return &c.latLon
}

func (c *Client) AltitudeEllipsoid() *AltitudeEllipsoid {
	return &c.altitudeEllipsoid
}

func (c *Client) VelocityXYZ() *VelocityXYZ {
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

func (c *Client) Temperature() *Temperature {
	return &c.temperature
}

func (c *Client) MagneticField() *MagneticField {
	return &c.magneticField
}

func (c *Client) RotationMatrix() *RotationMatrix {
	return &c.rotationMatrix
}

func (c *Client) FreeAcceleration() *FreeAcceleration {
	return &c.freeAcceleration
}

func (c *Client) Quaternion() *Quaternion {
	return &c.quaternion
}

func (c *Client) GNSSSatInfo() *GNSSSatInfo {
	return &c.gnssSatInfo
}

func (c *Client) PositionECEF() *PositionECEF {
	return &c.positionECEF
}

func (c *Client) send(ctx context.Context, message Message) error {
	deadline, _ := ctx.Deadline()
	if err := c.p.SetWriteDeadline(deadline); err != nil {
		return fmt.Errorf("send %v: %w", message.Identifier(), err)
	}
	if _, err := c.p.Write(message); err != nil {
		return fmt.Errorf("send %v: %w", message.Identifier(), err)
	}
	return nil
}

func (c *Client) receiveUntil(ctx context.Context, until MessageIdentifier) error {
	for {
		if err := c.Receive(ctx); err != nil {
			return fmt.Errorf("receive until %v: %w", until, err)
		}
		if c.MessageIdentifier() != until {
			continue
		}
		return nil
	}
}
