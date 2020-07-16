package xsens_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/einride/xsens-go"
	mockxsens "github.com/einride/xsens-go/test/mocks/xsens"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestClient_GoToConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	port := mockxsens.NewMockSerialPort(ctrl)
	client := xsens.NewClient(port)
	expectedGoToConfig := []byte{0xfa, 0xff, 0x30, 0x0, 0xd1}
	mtData2 := []byte{0xfa, 0xff, 0x36, 0x0, 0xcb}
	goToConfigAck := []byte{0xfa, 0xff, 0x31, 0x0, 0xd0}
	deadline := time.Unix(1, 2)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	// the client should send a GoToConfig message
	port.EXPECT().SetWriteDeadline(deadline)
	port.EXPECT().Write(expectedGoToConfig)
	// and then ignore messages other than GoToConfigAck
	port.EXPECT().SetReadDeadline(deadline)
	port.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		copy(b, mtData2)
		return len(mtData2), nil
	})
	// until it receives a GoToConfigAck
	port.EXPECT().SetReadDeadline(deadline)
	port.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		copy(b, goToConfigAck)
		return len(goToConfigAck), nil
	})
	// when requesting GoToConfig
	assert.NilError(t, client.GoToConfig(ctx))
}

func TestClient_GoToMeasurement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	port := mockxsens.NewMockSerialPort(ctrl)
	client := xsens.NewClient(port)
	expectedGoToMeasurement := []byte{0xfa, 0xff, 0x10, 0x0, 0xf1}
	mtData2 := []byte{0xfa, 0xff, 0x36, 0x0, 0xcb}
	goToConfigAck := []byte{0xfa, 0xff, 0x31, 0x0, 0xd0}
	deadline := time.Unix(1, 2)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	// the client should send a GoToMeasurement message
	port.EXPECT().SetWriteDeadline(deadline)
	port.EXPECT().Write(expectedGoToMeasurement)
	// and then ignore messages other than MTData2
	port.EXPECT().SetReadDeadline(deadline)
	port.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		copy(b, goToConfigAck)
		return len(goToConfigAck), nil
	})
	// until it receives MTData2
	port.EXPECT().SetReadDeadline(deadline)
	port.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		copy(b, mtData2)
		return len(mtData2), nil
	})
	// when requesting GoToMeasurement
	assert.NilError(t, client.GoToMeasurement(ctx))
}

func TestClient_GetOutputConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	port := mockxsens.NewMockSerialPort(ctrl)
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
	deadline := time.Unix(1, 2)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	// the client should send a ReqOutputconfiguration message
	port.EXPECT().SetWriteDeadline(deadline)
	port.EXPECT().Write(expectedReqOutputConfiguration)
	// and when it receives a ReqOutputConfigurationAck with the output configuration
	port.EXPECT().SetReadDeadline(deadline)
	port.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		copy(b, reqOutputConfigurationAck)
		return len(reqOutputConfigurationAck), nil
	})
	// it should return the parsed output configuration
	actual, err := client.GetOutputConfiguration(ctx)
	assert.NilError(t, err)
	assert.DeepEqual(t, expected, actual)
}

func TestClient_SetOutputConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	port := mockxsens.NewMockSerialPort(ctrl)
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
	deadline := time.Unix(1, 2)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	// the client should send a SetOutputconfiguration message with the requested output configuration
	port.EXPECT().SetWriteDeadline(deadline)
	port.EXPECT().Write(expectedSetOutputConfiguration)
	// and then it should await a SetOutputConfigurationAck message
	port.EXPECT().SetReadDeadline(deadline)
	port.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		copy(b, setOutputConfigurationAck)
		return len(setOutputConfigurationAck), nil
	})
	// when requesting to set the output configuration
	assert.NilError(t, client.SetOutputConfiguration(ctx, outputConfiguration))
}

func TestClient_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	port := mockxsens.NewMockSerialPort(ctrl)
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
		tt := tt
		t.Run(tt.inputFile, func(t *testing.T) {
			f, err := os.Open(tt.inputFile)
			assert.NilError(t, err)
			defer func() {
				assert.NilError(t, f.Close())
			}()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			port := mockxsens.NewMockSerialPort(ctrl)
			deadline := time.Unix(1, 2)
			ctx, cancel := context.WithDeadline(context.Background(), deadline)
			defer cancel()
			port.EXPECT().SetReadDeadline(deadline).AnyTimes()
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
				assert.NilError(t, ioutil.WriteFile(tt.goldenFile, actual.Bytes(), 0o644))
			}
			requireGoldenFileContent(t, tt.goldenFile, actual.String())
		})
	}
}
