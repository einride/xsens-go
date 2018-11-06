package main

import (
	"go.uber.org/zap"

	"github.com/einride/xsens-go"
)

func main() {
	logger := zap.NewExample()
	serialPort, err := xsensgo.DefaultSerialPort()
	if err != nil {
		logger.Panic("Failed to open default Xsens serial port", zap.Error(err))
	}
	var data xsensgo.Data
	for {
		if err := data.Read(serialPort); err != nil {
			logger.Panic("Failed to read from Xsens serial port", zap.Error(err))
		}
		logger.Info("Data", zap.Any("data", data))
	}
}
