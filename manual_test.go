// +build Client

package xsensgo

import (
	"go.uber.org/zap"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests are used online, which requires a connected Client.
// Used for debugging and some verification.

func TestReadmsgs(t *testing.T) {
	test := assert.New(t)

	logger := zap.NewExample()
	prt, err := DefaultSerialPort()
	test.Nil(err)
	client := NewClient(prt, logger)
	test.Nil(err)
	defer client.Close()

	err = client.readMessages(func(data *Data, err error) {
		logger.Info("got this", zap.Any("data", data), zap.Error(err))
	})
	test.Nil(err)
}

func TestRun(t *testing.T) {
	test := assert.New(t)

	logger := zap.NewExample()
	prt, err := DefaultSerialPort()
	test.Nil(err)

	client := NewClient(prt, logger)
	test.Nil(err)
	defer client.Close()

	err = client.Run(func(data *Data, err error) {
		logger.Info("Got this", zap.Any("data", data), zap.Error(err))
	})
	test.Nil(err)
}

func TestHeadingCalc(t *testing.T) {
	test := assert.New(t)

	logger := zap.NewExample()
	prt, err := DefaultSerialPort()
	test.Nil(err)

	client := NewClient(prt, logger)
	test.Nil(err)
	defer client.Close()

	err = client.Run(func(data *Data, err error) {
		heading := data.Heading()
		logger.Info("Heading is", zap.Float64("heading", heading))
		// Set north as reference
		wanted := math.Atan2(1, 0) * 180 / math.Pi
		logger.Info("wanted:", zap.Float64("wanted",wanted))
		// Check angle error (reference - actual)
		headingError := wanted - heading
		switch {
		case headingError > 180.0:
			headingError -= 360.0
		case headingError < -180.0:
			headingError += 360
		}
		// When Client is facing north, headingError should be 0
		logger.Info("Heading error is", zap.Float64("headingError",headingError))
	})
}
