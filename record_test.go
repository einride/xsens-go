// +build record

package xsens_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/einride/xsens-go"
	"github.com/stretchr/testify/require"
	"github.com/tarm/serial"
)

func TestRecord_TestData(t *testing.T) {
	const numMessages = 100
	for _, tt := range []struct {
		outputFile       string
		outputConfigFile string
		outputConfig     xsens.OutputConfiguration
	}{
		{
			outputFile:       "testdata/1/output.bin",
			outputConfigFile: "testdata/1/outputconfig.bin",
			outputConfig: xsens.OutputConfiguration{
				{
					DataType:        xsens.DataTypePacketCounter,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeSampleTimeFine,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeUTCTime,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeStatusWord,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:         xsens.DataTypeEulerAngles,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFloat32,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeAcceleration,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeDeltaV,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeRateOfTurn,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeDeltaQ,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeLatLon,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeAltitudeEllipsoid,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:         xsens.DataTypeVelocityXYZ,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFloat32,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeGNSSPVTData,
					OutputFrequency: 4,
				},
			},
		},
		{
			outputFile:       "testdata/2/output.bin",
			outputConfigFile: "testdata/2/outputconfig.bin",
			outputConfig: xsens.OutputConfiguration{
				{
					DataType:        xsens.DataTypePacketCounter,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeSampleTimeFine,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeUTCTime,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeStatusByte,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:         xsens.DataTypeEulerAngles,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFP1220,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeAcceleration,
					Precision:       xsens.PrecisionFP1220,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeDeltaV,
					Precision:       xsens.PrecisionFP1220,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeRateOfTurn,
					Precision:       xsens.PrecisionFP1220,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeDeltaQ,
					Precision:       xsens.PrecisionFP1220,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeLatLon,
					Precision:       xsens.PrecisionFP1220,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeAltitudeEllipsoid,
					Precision:       xsens.PrecisionFP1220,
					OutputFrequency: 100,
				},
				{
					DataType:         xsens.DataTypeVelocityXYZ,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFP1220,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeGNSSPVTData,
					OutputFrequency: 4,
				},
			},
		},
		{
			outputFile:       "testdata/3/output.bin",
			outputConfigFile: "testdata/3/outputconfig.bin",
			outputConfig: xsens.OutputConfiguration{
				{
					DataType:        xsens.DataTypePacketCounter,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeSampleTimeFine,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeUTCTime,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeStatusByte,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:         xsens.DataTypeEulerAngles,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFP1632,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeAcceleration,
					Precision:       xsens.PrecisionFP1632,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeDeltaV,
					Precision:       xsens.PrecisionFP1632,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeRateOfTurn,
					Precision:       xsens.PrecisionFP1632,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeDeltaQ,
					Precision:       xsens.PrecisionFP1632,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeLatLon,
					Precision:       xsens.PrecisionFP1632,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeAltitudeEllipsoid,
					Precision:       xsens.PrecisionFP1632,
					OutputFrequency: 100,
				},
				{
					DataType:         xsens.DataTypeVelocityXYZ,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFP1632,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeGNSSPVTData,
					OutputFrequency: 4,
				},
			},
		},
		{
			outputFile:       "testdata/4/output.bin",
			outputConfigFile: "testdata/4/outputconfig.bin",
			outputConfig: xsens.OutputConfiguration{
				{
					DataType:        xsens.DataTypePacketCounter,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeSampleTimeFine,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeUTCTime,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeStatusByte,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:         xsens.DataTypeEulerAngles,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFloat64,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeAcceleration,
					Precision:       xsens.PrecisionFloat64,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeDeltaV,
					Precision:       xsens.PrecisionFloat64,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeRateOfTurn,
					Precision:       xsens.PrecisionFloat64,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeDeltaQ,
					Precision:       xsens.PrecisionFloat64,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeLatLon,
					Precision:       xsens.PrecisionFloat64,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeAltitudeEllipsoid,
					Precision:       xsens.PrecisionFloat64,
					OutputFrequency: 100,
				},
				{
					DataType:         xsens.DataTypeVelocityXYZ,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFloat64,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeGNSSPVTData,
					OutputFrequency: 4,
				},
			},
		},
		{
			outputFile:       "testdata/5/output.bin",
			outputConfigFile: "testdata/5/outputconfig.bin",
			outputConfig: xsens.OutputConfiguration{
				{
					DataType:        xsens.DataTypePacketCounter,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeSampleTimeFine,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeSampleTimeCoarse,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeUTCTime,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:        xsens.DataTypeStatusWord,
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataType:         xsens.DataTypeRotationMatrix,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFloat32,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeFreeAcceleration,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:         xsens.DataTypeQuaternion,
					CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
					Precision:        xsens.PrecisionFloat32,
					OutputFrequency:  100,
				},
				{
					DataType:        xsens.DataTypeTemperature,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypePositionECEF,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeMagneticField,
					Precision:       xsens.PrecisionFloat32,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeBaroPressure,
					OutputFrequency: 100,
				},
				{
					DataType:        xsens.DataTypeGNSSPVTData,
					OutputFrequency: 4,
				},
				{
					DataType:        xsens.DataTypeGNSSSatInfo,
					OutputFrequency: 4,
				},
			},
		},
	} {
		t.Run(tt.outputFile, func(t *testing.T) {
			// ensure output folder present
			require.NoError(t, os.MkdirAll(filepath.Dir(tt.outputFile), 0774))
			// write output config
			outputConfigData, err := tt.outputConfig.Marshal()
			require.NoError(t, err)
			require.NoError(t, ioutil.WriteFile(tt.outputConfigFile, outputConfigData, 0644))
			// open Xsens port
			port, err := serial.OpenPort(&serial.Config{
				Name:     "/dev/ttyUSB0",
				Baud:     xsens.DefaultSerialBaudRate,
				Size:     xsens.MinLengthOfMessage,
				StopBits: xsens.DefaultSerialStopBits,
			})
			require.NoError(t, err)
			defer func() {
				require.NoError(t, port.Close())
			}()
			sc := xsens.NewMessageScanner(port)
			// go to config
			if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierGotoConfig, nil)); err != nil {
				panic(err)
			}
			for sc.Scan() && sc.Message().Identifier() != xsens.MessageIdentifierGotoConfig.Ack() {
				// scan for ack
			}
			if sc.Err() != nil {
				panic(err)
			}
			// set output config
			outputConfigMessage := xsens.NewMessage(xsens.MessageIdentifierSetOutputConfiguration, outputConfigData)
			if _, err := port.Write(outputConfigMessage); err != nil {
				panic(err)
			}
			for sc.Scan() && sc.Message().Identifier() != xsens.MessageIdentifierSetOutputConfiguration.Ack() {
				// scan for ack
			}
			// go to measurement
			if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierGotoMeasurement, nil)); err != nil {
				panic(err)
			}
			out, err := os.Create(tt.outputFile)
			require.NoError(t, err)
			for i := 0; i < numMessages; {
				require.True(t, sc.Scan())
				msg := sc.Message()
				require.NoError(t, msg.Validate())
				if msg.Identifier() != xsens.MessageIdentifierMTData2 {
					continue
				}
				i++
				_, err = out.Write(msg)
				require.NoError(t, err)
			}
			require.NoError(t, sc.Err())
		})
	}
}

var _ = xsens.OutputConfiguration{
	{
		DataType:        xsens.DataTypePacketCounter,
		OutputFrequency: xsens.MaxOutputFrequency,
	},
	{
		DataType:        xsens.DataTypeSampleTimeFine,
		OutputFrequency: xsens.MaxOutputFrequency,
	},
	{
		DataType:        xsens.DataTypeUTCTime,
		OutputFrequency: xsens.MaxOutputFrequency,
	},
	{
		DataType:        xsens.DataTypeStatusWord,
		OutputFrequency: xsens.MaxOutputFrequency,
	},
	{
		DataType:         xsens.DataTypeEulerAngles,
		CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
		Precision:        xsens.PrecisionFloat32,
		OutputFrequency:  100,
	},
}
