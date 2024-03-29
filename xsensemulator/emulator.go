package xsensemulator

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"go.einride.tech/xsens"
)

var (
	ErrNotInMeasurementMode     = errors.New("not in measurement mode")
	ErrNotInOutputConfiguration = errors.New("not in output configuration")
)

type Emulator struct {
	port io.ReadWriteCloser
	w    *bufio.Writer
	sc   *bufio.Scanner

	mutex                 sync.Mutex
	outputConf            xsens.OutputConfiguration
	lastMessageIdentifier xsens.MessageIdentifier
}

func NewEmulator(p io.ReadWriteCloser) *Emulator {
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
		// Give a chance to quit upon context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
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
		case xsens.MessageIdentifierSetOutputConfiguration:
			if err := e.outputConf.Unmarshal(m.Data()); err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			e.mutex.Lock()
			e.lastMessageIdentifier = xsens.MessageIdentifierSetOutputConfiguration
			e.mutex.Unlock()
			_, err := e.port.Write(
				xsens.NewMessage(xsens.MessageIdentifierSetOutputConfigurationAck, nil),
			)
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
		case xsens.MessageIdentifierGotoMeasurement:
			e.mutex.Lock()
			e.lastMessageIdentifier = xsens.MessageIdentifierMTData2
			e.mutex.Unlock()
			_, err := e.port.Write(xsens.NewMessage(xsens.MessageIdentifierMTData2, nil))
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
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

func (e *Emulator) MarshalMessage(
	measurement xsens.MeasurementData,
	dataType xsens.DataType,
) ([]byte, error) {
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
		return nil, ErrNotInOutputConfiguration
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
