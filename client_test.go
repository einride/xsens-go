package xsens_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.einride.tech/xsens"
	"go.einride.tech/xsens/mocks/mockserial"
	"go.einride.tech/xsens/xsensemulator"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestClient_GoToConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedGoToConfig := []byte{0xfa, 0xff, 0x30, 0x0, 0xd1}
	mtData2 := []byte{0xfa, 0xff, 0x36, 0x0, 0xcb}
	goToConfigAck := []byte{0xfa, 0xff, 0x31, 0x0, 0xd0}

	port := mockserial.NewMockPort(ctrl)
	client := xsens.NewClient(port)
	defer client.Close()

	// the client should send a GoToConfig message
	port.EXPECT().Write(expectedGoToConfig)

	// and then ignore messages other than GoToConfigAck
	port.EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(b []byte) (int, error) {
			copy(b, mtData2)
			return len(mtData2), nil
		})

	// until it receives a GoToConfigAck
	port.EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(b []byte) (int, error) {
			copy(b, goToConfigAck)
			return len(goToConfigAck), nil
		})

	// expect client to close
	port.EXPECT().Close()

	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// when requesting GoToConfig
	assert.NilError(t, client.GoToConfig(ctx))
}

func TestClient_GoToMeasurement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedGoToMeasurement := []byte{0xfa, 0xff, 0x10, 0x0, 0xf1}
	mtData2 := []byte{0xfa, 0xff, 0x36, 0x0, 0xcb}
	goToMeasurementAck := []byte{0xfa, 0xff, 0x31, 0x0, 0xd0}

	port := mockserial.NewMockPort(ctrl)
	client := xsens.NewClient(port)
	defer client.Close()

	// the client should send a GoToMeasurement message
	port.EXPECT().Write(expectedGoToMeasurement)

	// and then ignore messages other than MTData2
	port.EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(b []byte) (int, error) {
			copy(b, goToMeasurementAck)
			return len(goToMeasurementAck), nil
		})

	// until it receives MTData2
	port.EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(b []byte) (int, error) {
			copy(b, mtData2)
			return len(mtData2), nil
		})

	// expect client to close
	port.EXPECT().Close()

	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// when requesting GoToMeasurement
	assert.NilError(t, client.GoToMeasurement(ctx))
}

func TestClient_GetOutputConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	port := mockserial.NewMockPort(ctrl)
	client := xsens.NewClient(port)

	expectedReqOutputConfiguration := []byte{0xfa, 0xff, 0xc0, 0x0, 0x41}
	expected := xsens.OutputConfiguration{
		{
			DataIdentifier: xsens.DataIdentifier{
				DataType:  xsens.DataTypeLatLon,
				Precision: xsens.PrecisionFloat32,
			},
			OutputFrequency: 100,
		},
		{
			DataIdentifier: xsens.DataIdentifier{
				DataType:         xsens.DataTypeQuaternion,
				Precision:        xsens.PrecisionFloat64,
				CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
			},
			OutputFrequency: 200,
		},
	}

	reqOutputConfigurationAck := []byte{0xfa, 0xff, 0xc1, 0x8, 0x50, 0x40, 0x0, 0x64, 0x20, 0x13, 0x0, 0xc8, 0x49}

	// the client should send a ReqOutputconfiguration message
	port.EXPECT().Write(expectedReqOutputConfiguration)

	// and when it receives a ReqOutputConfigurationAck with the output configuration
	port.EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(b []byte) (int, error) {
			copy(b, reqOutputConfigurationAck)
			return len(reqOutputConfigurationAck), nil
		})

	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// it should return the parsed output configuration
	actual, err := client.GetOutputConfiguration(ctx)
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestClient_SetOutputConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	port := mockserial.NewMockPort(ctrl)
	client := xsens.NewClient(port)

	outputConfiguration := xsens.OutputConfiguration{
		{
			DataIdentifier: xsens.DataIdentifier{
				DataType:  xsens.DataTypeLatLon,
				Precision: xsens.PrecisionFloat32,
			},
			OutputFrequency: 100,
		},
		{
			DataIdentifier: xsens.DataIdentifier{
				DataType:         xsens.DataTypeQuaternion,
				Precision:        xsens.PrecisionFloat64,
				CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
			},
			OutputFrequency: 200,
		},
	}
	setOutputConfigurationAck := []byte{0xfa, 0xff, 0xc1, 0x0, 0x40}
	expectedSetOutputConfiguration := []byte{0xfa, 0xff, 0xc0, 0x8, 0x50, 0x40, 0x0, 0x64, 0x20, 0x13, 0x0, 0xc8, 0x4a}

	// the client should send a SetOutputconfiguration message with the requested output configuration
	port.EXPECT().Write(expectedSetOutputConfiguration)
	// and then it should await a SetOutputConfigurationAck message
	port.EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(b []byte) (int, error) {
			copy(b, setOutputConfigurationAck)
			return len(setOutputConfigurationAck), nil
		})

	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// when requesting to set the output configuration
	assert.NilError(t, client.SetOutputConfiguration(ctx, outputConfiguration))
}

func TestClient_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	port := mockserial.NewMockPort(ctrl)
	client := xsens.NewClient(port)

	err := errors.New("boom")
	port.EXPECT().Close().Return(err)

	assert.Assert(t, errors.Is(client.Close(), err))
}

func TestClient_ScanMeasurementData(t *testing.T) {
	for _, tt := range []struct {
		inputFile  string
		goldenFile string
	}{
		{inputFile: "testdata/1/output.bin", goldenFile: "testdata/1/output.client.golden"},
		{inputFile: "testdata/2/output.bin", goldenFile: "testdata/2/output.client.golden"},
		{inputFile: "testdata/3/output.bin", goldenFile: "testdata/3/output.client.golden"},
		{inputFile: "testdata/4/output.bin", goldenFile: "testdata/4/output.client.golden"},
		{inputFile: "testdata/5/output.bin", goldenFile: "testdata/5/output.client.golden"},
	} {
		t.Run(tt.inputFile, func(t *testing.T) {
			f, err := os.Open(tt.inputFile)
			assert.NilError(t, err)
			defer func() {
				assert.NilError(t, f.Close())
			}()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			port := mockserial.NewMockPort(ctrl)
			deadline := time.Now().Add(100 * time.Millisecond)
			ctx, cancel := context.WithDeadline(context.Background(), deadline)
			defer cancel()
			port.EXPECT().Read(gomock.Any()).AnyTimes().DoAndReturn(f.Read)
			client := xsens.NewClient(port)
			var actual bytes.Buffer
			printf := func(format string, args ...interface{}) {
				_, err := fmt.Fprintf(&actual, format, args...)
				assert.NilError(t, err)
			}
			for {
				err := client.Receive(ctx)
				if errors.Is(err, io.EOF) {
					break
				}
				assert.NilError(t, err)
				printf("%v\n", client.MessageIdentifier())
				for client.ScanMeasurementData() {
					printf("\t%v\n", client.DataType())
					printf("\t%+v\n", client.MeasurementData())
				}
				printf("\n")
			}
			if shouldUpdateGoldenFiles() {
				assert.NilError(t, os.WriteFile(tt.goldenFile, actual.Bytes(), 0o600))
			}
			requireGoldenFileContent(t, tt.goldenFile, actual.String())
		})
	}
}

func TestUDPEmulator(t *testing.T) {
	addrEmulator := "127.0.0.1:24001"
	addrClient := "127.0.0.1:24002"

	timeout := 100 * time.Millisecond

	connEmulator, err := xsensemulator.NewUDPSerialPort(
		addrEmulator,
		addrClient,
		xsensemulator.WithTimeout(timeout),
	)
	assert.NilError(t, err)
	defer func() {
		assert.NilError(t, connEmulator.Close())
	}()

	connClient, err := xsensemulator.NewUDPSerialPort(
		addrClient,
		addrEmulator,
		xsensemulator.WithTimeout(timeout),
	)
	assert.NilError(t, err)

	emu := xsensemulator.NewEmulator(connEmulator)

	outputConf := xsens.OutputConfiguration{
		xsens.OutputConfigurationSetting{DataIdentifier: xsens.DataIdentifier{DataType: xsens.DataTypeUTCTime}},
		xsens.OutputConfigurationSetting{DataIdentifier: xsens.DataIdentifier{DataType: xsens.DataTypeEulerAngles}},
		xsens.OutputConfigurationSetting{DataIdentifier: xsens.DataIdentifier{DataType: xsens.DataTypeVelocityXYZ}},
		xsens.OutputConfigurationSetting{DataIdentifier: xsens.DataIdentifier{DataType: xsens.DataTypeAcceleration}},
		xsens.OutputConfigurationSetting{DataIdentifier: xsens.DataIdentifier{DataType: xsens.DataTypeFreeAcceleration}},
		xsens.OutputConfigurationSetting{DataIdentifier: xsens.DataIdentifier{DataType: xsens.DataTypeRateOfTurn}},
	}

	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	var g errgroup.Group

	g.Go(func() error {
		err := emu.Receive(ctx)
		if !strings.Contains(err.Error(), "timeout") {
			return err
		}
		return nil
	})

	g.Go(func() error {
		client := xsens.NewClient(connClient)

		assert.NilError(t, client.SetOutputConfiguration(ctx, outputConf))
		assert.NilError(t, client.GoToConfig(ctx))

		assert.NilError(t, connClient.Close())
		return nil
	})

	assert.NilError(t, g.Wait())
}
