package xsensemulator

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"go.einride.tech/xsens"
)

var ErrNotInMeasurementMode = errors.New("not in measurement mode")

type UDPSerialPort struct {
	OriginConn      *net.UDPConn
	DestinationAddr *net.UDPAddr
}

func NewUDPSerialPort(origin, destination string) (*UDPSerialPort, error) {
	udpOriginAddr, err := net.ResolveUDPAddr("udp", origin)
	if err != nil {
		return nil, fmt.Errorf("new udp serial port: %w", err)
	}
	udpDestinationAddr, err := net.ResolveUDPAddr("udp", destination)
	if err != nil {
		return nil, fmt.Errorf("new udp serial port: %w", err)
	}
	originConn, err := net.ListenUDP("udp", udpOriginAddr)
	if err != nil {
		return nil, fmt.Errorf("new udp serial port: %w", err)
	}
	return &UDPSerialPort{
		OriginConn:      originConn,
		DestinationAddr: udpDestinationAddr,
	}, nil
}

func (t *UDPSerialPort) Read(p []byte) (int, error) {
	n, _, err := t.OriginConn.ReadFromUDP(p)
	return n, err
}

func (t *UDPSerialPort) Write(p []byte) (n int, err error) {
	return t.OriginConn.WriteToUDP(p, t.DestinationAddr)
}

func (t *UDPSerialPort) Close() error {
	return t.OriginConn.Close()
}

func (t *UDPSerialPort) SetReadDeadline(t2 time.Time) error {
	return t.OriginConn.SetReadDeadline(t2)
}

func (t *UDPSerialPort) SetWriteDeadline(t2 time.Time) error {
	return t.OriginConn.SetWriteDeadline(t2)
}

type Emulator struct {
	port xsens.SerialPort
	w    *bufio.Writer
	sc   *bufio.Scanner

	mutex                 sync.Mutex
	outputConf            xsens.OutputConfiguration
	lastMessageIdentifier xsens.MessageIdentifier
}

func NewEmulator(p xsens.SerialPort) *Emulator {
	sc := bufio.NewScanner(p)
	sc.Split(xsens.ScanMessages)
	return &Emulator{
		w:    bufio.NewWriter(p),
		sc:   sc,
		port: p,
	}
}

func (e *Emulator) Close() error {
	if err := e.port.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	return nil
}

func (e *Emulator) SetOutputConguration(configuration xsens.OutputConfiguration) {
	e.outputConf = configuration
}

func (e *Emulator) SetSendMode() {
	e.lastMessageIdentifier = xsens.MessageIdentifierMTData2
}

func (e *Emulator) Receive(ctx context.Context) error {
	for {
		deadline, ok := ctx.Deadline()
		if !ok {
			return fmt.Errorf("no deadline")
		}
		if err := e.port.SetReadDeadline(deadline); err != nil {
			return fmt.Errorf("xsens client: receive: %w", err)
		}
		if !e.sc.Scan() {
			if e.sc.Err() != nil {
				return fmt.Errorf("receive: %w", e.sc.Err())
			}
			return fmt.Errorf("receive: %w", io.EOF)
		}
		var m xsens.Message = e.sc.Bytes()
		if err := m.Validate(); err != nil {
			return fmt.Errorf("receive: %w", err)
		}
		switch m.Identifier() {
		case xsens.MessageIdentifierGotoConfig:
			e.mutex.Lock()
			e.lastMessageIdentifier = xsens.MessageIdentifierGotoConfig
			e.mutex.Unlock()
			_, err := e.port.Write(xsens.NewMessage(xsens.MessageIdentifierGotoConfigAck, nil))
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			continue
		case xsens.MessageIdentifierSetOutputConfiguration:
			if err := e.outputConf.Unmarshal(m.Data()); err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			e.mutex.Lock()
			e.lastMessageIdentifier = xsens.MessageIdentifierSetOutputConfiguration
			e.mutex.Unlock()
			_, err := e.port.Write(xsens.NewMessage(xsens.MessageIdentifierSetOutputConfigurationAck, nil))
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			continue
		case xsens.MessageIdentifierGotoMeasurement:
			e.mutex.Lock()
			e.lastMessageIdentifier = xsens.MessageIdentifierMTData2
			e.mutex.Unlock()
			_, err := e.port.Write(xsens.NewMessage(xsens.MessageIdentifierMTData2, nil))
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			continue
		}

		// Give a chance to quit upon context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func (e *Emulator) Transmit(m xsens.Message) error {
	if e.lastMessageIdentifier != xsens.MessageIdentifierMTData2 {
		return fmt.Errorf("transmit: %w", ErrNotInMeasurementMode)
	}
	if err := m.Validate(); err != nil {
		return fmt.Errorf("transmit: %w", err)
	}
	if _, err := e.port.Write(m); err != nil {
		return fmt.Errorf("transmit: %w", err)
	}
	return nil
}

func (e *Emulator) MarshalMessage(measurement xsens.MeasurementData, dataType xsens.DataType) ([]byte, error) {
	var id xsens.DataIdentifier
	var isSet bool
	for _, d := range e.outputConf {
		if d.DataType != dataType {
			continue
		}
		isSet = true
		id = d.DataIdentifier
	}
	if !isSet {
		return nil, errors.New("not in output configuration")
	}
	packetData, err := measurement.MarshalMTData2Packet(id)
	if err != nil {
		return nil, fmt.Errorf("transmit: %w", err)
	}
	return packetData, nil
}

func (e *Emulator) LastMessageIdentifier() xsens.MessageIdentifier {
	var id xsens.MessageIdentifier
	e.mutex.Lock()
	id = e.lastMessageIdentifier
	e.mutex.Unlock()
	return id
}
