package main

import (
	"os"

	"github.com/einride/xsens-go"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	prt, err := xsensgo.DefaultSerialPort()
	if err != nil {
		logger.Fatal("could not open port", zap.Error(err))
	}

	client := xsensgo.NewClient(prt, logger)
	if err != nil {
		logger.Fatal("Got error opening new xsens %v", zap.Error(err))
		return
	}

	defer client.Close()
	err = client.Run(func(data *xsensgo.Data, err error) {
		if err != nil {
			logger.Warn("Got error from xsens data %v", zap.Error(err))
			return
		}

		if os.Getenv("RECIEVE_STATE_LOG_FROM_IMU") != "" {
			logger.Info("Received ", zap.Any("xsens data", data))
		}
	})
	if err != nil {
		logger.Error("Error in xsens Run %v", zap.Error(err))
	}
}
