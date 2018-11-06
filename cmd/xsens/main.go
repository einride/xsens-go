package main

import (
	"github.com/einride/xsens-go"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	serialPort, err := xsens.DefaultSerialPort()
	if err != nil {
		logger.Panic("Failed to open default Xsens serial port", zap.Error(err))
	}
	var data xsens.Data
	for {
		if err := data.Read(serialPort); err != nil {
			logger.Panic("Failed to read from Xsens serial port", zap.Error(err))
		}
		logger.Info("Data", zap.Any("data", data))
	}
}
