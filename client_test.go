package xsens_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/einride/xsens-go"
	mockxsens "github.com/einride/xsens-go/test/mocks/xsens"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
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
	port.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		copy(b, goToConfigAck)
		return len(goToConfigAck), nil
	})
	// when requesting GoToConfig
	require.NoError(t, client.GoToConfig(ctx))
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
	port.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
		copy(b, mtData2)
		return len(mtData2), nil
	})
	// when requesting GoToMeasurement
	require.NoError(t, client.GoToMeasurement(ctx))
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
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestClient_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	port := mockxsens.NewMockSerialPort(ctrl)
	client := xsens.NewClient(port)
	err := xerrors.New("boom")
	port.EXPECT().Close().Return(err)
	require.True(t, xerrors.Is(client.Close(), err))
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
			require.NoError(t, err)
			defer func() {
				require.NoError(t, f.Close())
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
				require.NoError(t, err)
			}
			printData := func(data interface{}) {
				printf("\t%+v\n", data)
			}
			for {
				err := client.Receive(ctx)
				if xerrors.Is(err, io.EOF) {
					break
				}
				require.NoError(t, err)
				printf("%v\n", client.MessageIdentifier())
				for client.ScanMeasurementData() {
					printf("\t%v\n", client.DataType())
					switch client.DataType() {
					case xsens.DataTypePacketCounter:
						printData(client.PacketCounter())
					case xsens.DataTypeSampleTimeCoarse:
						printData(client.SampleTimeCoarse())
					case xsens.DataTypeSampleTimeFine:
						printData(client.SampleTimeFine())
					case xsens.DataTypeUTCTime:
						printData(client.UTCTime())
					case xsens.DataTypeStatusByte:
						printData(client.StatusByte())
					case xsens.DataTypeStatusWord:
						printData(client.StatusWord())
					case xsens.DataTypeEulerAngles:
						printData(client.EulerAngles())
					case xsens.DataTypeAcceleration:
						printData(client.Acceleration())
					case xsens.DataTypeDeltaV:
						printData(client.DeltaV())
					case xsens.DataTypeRateOfTurn:
						printData(client.RateOfTurn())
					case xsens.DataTypeDeltaQ:
						printData(client.DeltaQ())
					case xsens.DataTypeLatLon:
						printData(client.LatLon())
					case xsens.DataTypeAltitudeEllipsoid:
						printData(client.AltitudeEllipsoid())
					case xsens.DataTypeVelocityXYZ:
						printData(client.VelocityXYZ())
					case xsens.DataTypeGNSSPVTData:
						printData(client.GNSSPVTData())
					case xsens.DataTypeBaroPressure:
						printData(client.BaroPressure())
					case xsens.DataTypeTemperature:
						printData(client.Temperature())
					case xsens.DataTypeMagneticField:
						printData(client.MagneticField())
					case xsens.DataTypeRotationMatrix:
						printData(client.RotationMatrix())
					case xsens.DataTypeFreeAcceleration:
						printData(client.FreeAcceleration())
					case xsens.DataTypeQuaternion:
						printData(client.Quaternion())
					case xsens.DataTypeGNSSSatInfo:
						printData(client.GNSSSatInfo())
					case xsens.DataTypePositionECEF:
						printData(client.PositionECEF())
					default:
						t.Fatalf("unhandled data type: %v", client.DataType())
					}
				}
				printf("\n")
			}
			if shouldUpdateGoldenFiles() {
				require.NoError(t, ioutil.WriteFile(tt.goldenFile, actual.Bytes(), 0644))
			}
			requireGoldenFileContent(t, tt.goldenFile, actual.String())
		})
	}
}
