package xsens

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
)

type Emulator struct {
	port io.ReadWriteCloser
	w    *bufio.Writer
	sc   *bufio.Scanner

	mutex                 sync.Mutex
	outputConf            OutputConfiguration
	lastMessageIdentifier MessageIdentifier
}

func NewEmulator(p io.ReadWriteCloser) *Emulator {
	sc := bufio.NewScanner(p)
	sc.Split(ScanMessages)
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

func (e *Emulator) SetOutputConguration(configuration OutputConfiguration) {
	e.outputConf = configuration
}

func (e *Emulator) SetSendMode() {
	e.lastMessageIdentifier = MessageIdentifierMTData2
}

func (e *Emulator) Receive() error {
	for {
		if !e.sc.Scan() {
			if e.sc.Err() != nil {
				return fmt.Errorf("receive: %w", e.sc.Err())
			}
			return fmt.Errorf("receive: %w", io.EOF)
		}
		var m Message = e.sc.Bytes()
		if err := m.Validate(); err != nil {
			return fmt.Errorf("receive: %w", err)
		}
		switch m.Identifier() {
		case MessageIdentifierGotoConfig:
			_, err := e.port.Write(NewMessage(MessageIdentifierGotoConfigAck, nil))
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			log.Println("GOTO MessageIdentifierGotoConfig")
			continue
		case MessageIdentifierSetOutputConfiguration:
			log.Println("GOTO MessageIdentifierSetOutputConfiguration")
			if err := e.outputConf.Unmarshal(m.Data()); err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			e.mutex.Lock()
			e.mutex.Unlock()
			_, err := e.port.Write(NewMessage(MessageIdentifierSetOutputConfigurationAck, nil))
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			continue
		case MessageIdentifierGotoMeasurement:
			e.mutex.Lock()
			e.lastMessageIdentifier = MessageIdentifierMTData2
			e.mutex.Unlock()
			_, err := e.port.Write(NewMessage(MessageIdentifierMTData2, nil))
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			continue
		}
	}
}

func (e *Emulator) Transmit(m Message) error {
	if e.lastMessageIdentifier != MessageIdentifierMTData2 {
		return nil
	}
	if err := m.Validate(); err != nil {
		return fmt.Errorf("transmit: %w", err)
	}
	if _, err := e.port.Write(m); err != nil {
		return fmt.Errorf("transmit: %w", err)
	}
	return nil
}

func (e *Emulator) MarshalMessage(measurement MeasurementData, dataType DataType) ([]byte, error) {
	var id DataIdentifier
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
	packetData, err := measurement.marshalMTData2Packet(id)
	if err != nil {
		return nil, fmt.Errorf("transmit: %w", err)
	}
	return packetData, nil
}
