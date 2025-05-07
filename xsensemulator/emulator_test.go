package xsensemulator_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.einride.tech/xsens"
	"go.einride.tech/xsens/mocks/mockserial"
	"go.einride.tech/xsens/xsensemulator"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestEmulator_Convert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader, writer := io.Pipe()

	port1 := mockserial.NewMockPort(ctrl)
	port1.EXPECT().
		Read(gomock.Any()).
		Times(1).
		DoAndReturn(func(b []byte) (int, error) {
			n, err := reader.Read(b)
			if errors.Is(err, io.ErrClosedPipe) {
				return 0, io.EOF
			}
			assert.NilError(t, reader.Close())
			return n, err
		})

	c := xsens.NewClient(port1)

	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	expectedData := &xsens.LatLon{
		Lat: 1,
		Lon: 2,
	}

	var didReachEnd bool
	var g errgroup.Group

	g.Go(func() error {
		for {
			err := c.Receive(ctx)
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return err
			}
			for c.ScanMeasurementData() {
				n, ok := c.MeasurementData().(*xsens.LatLon)
				assert.Equal(t, ok, true)
				assert.Equal(t, expectedData.Lon, n.Lon)
				assert.Equal(t, expectedData.Lat, n.Lat)
				didReachEnd = true
				return nil
			}
		}
	})

	// then emulator receives
	port2 := mockserial.NewMockPort(ctrl)
	port2.EXPECT().
		Write(gomock.Any()).
		AnyTimes().
		DoAndReturn(writer.Write)

	emulator := xsensemulator.NewEmulator(port2)
	emulator.SetOutputConguration(xsens.OutputConfiguration{
		{
			DataIdentifier: xsens.DataIdentifier{
				DataType:  xsens.DataTypeLatLon,
				Precision: xsens.PrecisionFP1632,
			},
			OutputFrequency: 100,
		},
	})

	emulator.SetSendMode()
	m, err := emulator.MarshalMessage(expectedData, xsens.DataTypeLatLon)
	mes := xsens.NewMessage(xsens.MessageIdentifierMTData2, m)

	assert.NilError(t, err)
	assert.NilError(t, emulator.Transmit(mes))
	assert.NilError(t, g.Wait())
	assert.Check(t, didReachEnd)
}

func TestEmulator_Output(t *testing.T) {
	for _, tt := range []struct {
		inputFile               string
		outputConfigurationFile string
	}{
		{inputFile: "../testdata/1/output.bin", outputConfigurationFile: "../testdata/1/outputconfig.bin"},
		{inputFile: "../testdata/2/output.bin", outputConfigurationFile: "../testdata/2/outputconfig.bin"},
		{inputFile: "../testdata/3/output.bin", outputConfigurationFile: "../testdata/3/outputconfig.bin"},
		{inputFile: "../testdata/4/output.bin", outputConfigurationFile: "../testdata/4/outputconfig.bin"},
	} {
		t.Run(tt.inputFile, func(t *testing.T) {
			outputConf, err := os.ReadFile(tt.outputConfigurationFile)
			assert.NilError(t, err)
			o := xsens.OutputConfiguration{}
			assert.NilError(t, o.Unmarshal(outputConf))
			f, err := os.Open(tt.inputFile)
			assert.NilError(t, err)
			defer func() {
				assert.NilError(t, f.Close())
			}()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			deadline := time.Now().Add(100 * time.Millisecond)
			port1 := mockserial.NewMockPort(ctrl)

			ctx, cancel := context.WithDeadline(context.Background(), deadline)
			defer cancel()

			port1.EXPECT().Write(gomock.Any()).AnyTimes()
			port1.EXPECT().Read(gomock.Any()).AnyTimes().DoAndReturn(f.Read)

			port2 := mockserial.NewMockPort(ctrl)

			emulator := xsensemulator.NewEmulator(port2)
			emulator.SetOutputConguration(o)
			emulator.SetSendMode()

			var g errgroup.Group

			g.Go(func() error {
				client := xsens.NewClient(port1)
				assert.NilError(t, client.GoToMeasurement(ctx))
				for {
					err := client.Receive(ctx)
					if errors.Is(err, io.EOF) {
						break
					}
					if err != nil {
						return err
					}
					packets := []byte{}
					for client.ScanMeasurementData() {
						m, err := emulator.MarshalMessage(
							client.MeasurementData(),
							client.DataType(),
						)
						assert.NilError(t, err)
						packets = append(packets, m...)
					}
					mes := xsens.NewMessage(xsens.MessageIdentifierMTData2, packets)
					assert.Check(t, bytes.Equal(client.RawMessage(), mes))
				}
				return nil
			})
			assert.NilError(t, g.Wait())
		})
	}
}

func TestEmulator_Transmit(t *testing.T) {
	expectedData := &xsens.LatLon{
		Lat: 1,
		Lon: 2,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range []struct {
		name          string
		expectFunc    func(emulator *xsensemulator.Emulator, port *mockserial.MockPort)
		expectedError error
	}{
		{
			name: "not measurement mode (by default)",
			expectFunc: func(_ *xsensemulator.Emulator, _ *mockserial.MockPort) {
			},
			expectedError: xsensemulator.ErrNotInMeasurementMode,
		},

		{
			name: "in measurement mode",
			expectFunc: func(emulator *xsensemulator.Emulator, port *mockserial.MockPort) {
				emulator.SetSendMode()
				port.EXPECT().Write(gomock.Any())
			},
			expectedError: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// given
			port2 := mockserial.NewMockPort(ctrl)
			emulator := xsensemulator.NewEmulator(port2)
			emulator.SetOutputConguration(xsens.OutputConfiguration{
				{
					DataIdentifier: xsens.DataIdentifier{
						DataType:  xsens.DataTypeLatLon,
						Precision: xsens.PrecisionFP1632,
					},
					OutputFrequency: 100,
				},
			})

			// then expect
			tc.expectFunc(emulator, port2)

			// when
			m, err := emulator.MarshalMessage(expectedData, xsens.DataTypeLatLon)
			assert.NilError(t, err)

			msg := xsens.NewMessage(xsens.MessageIdentifierMTData2, m)

			err = emulator.Transmit(msg)
			assert.Assert(t, errors.Is(err, tc.expectedError))
		})
	}
}
