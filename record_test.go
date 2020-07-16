// +build record

package xsens_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/einride/xsens-go"
	"github.com/einride/xsens-go/pkg/serial"
	"gotest.tools/v3/assert"
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
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypePacketCounter,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeSampleTimeFine,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeUTCTime,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeStatusWord,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeEulerAngles,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeAcceleration,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeDeltaV,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeRateOfTurn,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeDeltaQ,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeLatLon,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeAltitudeEllipsoid,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeVelocityXYZ,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeGNSSPVTData,
					},
					OutputFrequency: 4,
				},
			},
		},
		{
			outputFile:       "testdata/2/output.bin",
			outputConfigFile: "testdata/2/outputconfig.bin",
			outputConfig: xsens.OutputConfiguration{
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypePacketCounter,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeSampleTimeFine,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeUTCTime,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeStatusByte,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeEulerAngles,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFP1220,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeAcceleration,
						Precision: xsens.PrecisionFP1220,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeDeltaV,
						Precision: xsens.PrecisionFP1220,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeRateOfTurn,
						Precision: xsens.PrecisionFP1220,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeDeltaQ,
						Precision: xsens.PrecisionFP1220,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeLatLon,
						Precision: xsens.PrecisionFP1220,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeAltitudeEllipsoid,
						Precision: xsens.PrecisionFP1220,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeVelocityXYZ,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFP1220,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeGNSSPVTData,
					},
					OutputFrequency: 4,
				},
			},
		},
		{
			outputFile:       "testdata/3/output.bin",
			outputConfigFile: "testdata/3/outputconfig.bin",
			outputConfig: xsens.OutputConfiguration{
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypePacketCounter,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeSampleTimeFine,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeUTCTime,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeStatusByte,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeEulerAngles,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFP1632,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeAcceleration,
						Precision: xsens.PrecisionFP1632,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeDeltaV,
						Precision: xsens.PrecisionFP1632,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeRateOfTurn,
						Precision: xsens.PrecisionFP1632,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeDeltaQ,
						Precision: xsens.PrecisionFP1632,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeLatLon,
						Precision: xsens.PrecisionFP1632,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeAltitudeEllipsoid,
						Precision: xsens.PrecisionFP1632,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeVelocityXYZ,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFP1632,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeGNSSPVTData,
					},
					OutputFrequency: 4,
				},
			},
		},
		{
			outputFile:       "testdata/4/output.bin",
			outputConfigFile: "testdata/4/outputconfig.bin",
			outputConfig: xsens.OutputConfiguration{
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypePacketCounter,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeSampleTimeFine,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeUTCTime,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeStatusByte,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeEulerAngles,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFloat64,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeAcceleration,
						Precision: xsens.PrecisionFloat64,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeDeltaV,
						Precision: xsens.PrecisionFloat64,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeRateOfTurn,
						Precision: xsens.PrecisionFloat64,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeDeltaQ,
						Precision: xsens.PrecisionFloat64,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeLatLon,
						Precision: xsens.PrecisionFloat64,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeAltitudeEllipsoid,
						Precision: xsens.PrecisionFloat64,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeVelocityXYZ,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFloat64,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeGNSSPVTData,
					},
					OutputFrequency: 4,
				},
			},
		},
		{
			outputFile:       "testdata/5/output.bin",
			outputConfigFile: "testdata/5/outputconfig.bin",
			outputConfig: xsens.OutputConfiguration{
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypePacketCounter,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeSampleTimeFine,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeSampleTimeCoarse,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeUTCTime,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeStatusWord,
					},
					OutputFrequency: xsens.MaxOutputFrequency,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeRotationMatrix,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeFreeAcceleration,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:         xsens.DataTypeQuaternion,
						CoordinateSystem: xsens.CoordinateSystemEastNorthUp,
						Precision:        xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeTemperature,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypePositionECEF,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeMagneticField,
						Precision: xsens.PrecisionFloat32,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeBaroPressure,
					},
					OutputFrequency: 100,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeGNSSPVTData,
					},
					OutputFrequency: 4,
				},
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType: xsens.DataTypeGNSSSatInfo,
					},
					OutputFrequency: 4,
				},
			},
		},
	} {
		t.Run(tt.outputFile, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			// ensure output folder present
			assert.NilError(t, os.MkdirAll(filepath.Dir(tt.outputFile), 0o774))
			// write output config
			outputConfigData, err := tt.outputConfig.Marshal()
			assert.NilError(t, err)
			assert.NilError(t, ioutil.WriteFile(tt.outputConfigFile, outputConfigData, 0o644))
			// open Xsens port
			port, err := serial.Open("/dev/ttyUSB0", serial.BaudRate115200)
			assert.NilError(t, err)
			client := xsens.NewClient(port)
			defer func() {
				assert.NilError(t, client.Close())
			}()
			assert.NilError(t, client.GoToConfig(ctx))
			assert.NilError(t, client.SetOutputConfiguration(ctx, tt.outputConfig))
			assert.NilError(t, client.GoToMeasurement(ctx))
			out, err := os.Create(tt.outputFile)
			assert.NilError(t, err)
			for i := 0; i < numMessages; {
				assert.NilError(t, client.Receive(ctx))
				msg := client.RawMessage()
				i++
				_, err = out.Write(msg)
				assert.NilError(t, err)
			}
		})
	}
}
