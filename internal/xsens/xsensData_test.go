package xsens

import (
	"bytes"
	"encoding/binary"
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

func Test1632(t *testing.T) {
	buf := bytes.NewReader([]byte{0x24, 0x39, 0x58, 0x10, 0x00, 0x03})
	var fp xsens1632
	err := binary.Read(buf, binary.BigEndian, &fp)
	assert.Nil(t, err)

	assert.Equal(t, int16(0x0003), fp.Integer)
	assert.Equal(t, uint32(0x24395810), fp.Fraction)
	assert.InDelta(t, fp.ToFloat(), 3.1415, 0.00001)
}

const (
	east  = 0
	north = 90
	//west                = 180
	westNeg             = -180
	south               = -90
	northEast           = 45
	northWest           = 135
	southEast           = -45
	southWest           = -135
	magneticDeclination = 3.71
)

func TestHeading(t *testing.T) {
	// Heading returns a value which has the magnetic north as reference, as this is not the same as geographical north,
	// the magnetic declination should be compensated for to get correctness in a geographical sense.
	//East
	dataEast := Data{
		Latlng: XDILatLng{Lat: 0, Lng: 0},
		Euler:  XDIEulerAngles{Roll: 0, Pitch: 0, Yaw: 0},
		Mag:    XDIMagneticXYZ{0, 1, 0}, // East --> Heading == 0
	}
	heading := Data.Heading(dataEast)
	log.Printf("Heading is: %v", heading-magneticDeclination)
	assert.Equal(t, east+magneticDeclination, heading)

	// North
	dataNorth := Data{
		Latlng: XDILatLng{Lat: 0, Lng: 0},
		Euler:  XDIEulerAngles{Roll: 0, Pitch: 0, Yaw: 0},
		Mag:    XDIMagneticXYZ{1, 0, 0}, // North --> Heading == 90
	}

	heading = Data.Heading(dataNorth)
	log.Printf("Heading is: %v", heading-magneticDeclination)
	assert.Equal(t, north+magneticDeclination, heading)

	// West
	dataWest := Data{
		Latlng: XDILatLng{Lat: 0, Lng: 0},
		Euler:  XDIEulerAngles{Roll: 0, Pitch: 0, Yaw: 0},
		Mag:    XDIMagneticXYZ{0, -1, 0}, // West --> Heading == -180||180
	}

	heading = Data.Heading(dataWest)
	log.Printf("Heading is: %v", heading-magneticDeclination)
	assert.Equal(t, westNeg+magneticDeclination, heading)

	// South
	dataSouth := Data{
		Latlng: XDILatLng{Lat: 0, Lng: 0},
		Euler:  XDIEulerAngles{Roll: 0, Pitch: 0, Yaw: 0},
		Mag:    XDIMagneticXYZ{-1, 0, 0}, // South --> Heading == -90
	}

	heading = Data.Heading(dataSouth)
	log.Printf("Heading is: %v", heading-magneticDeclination)
	assert.Equal(t, south+magneticDeclination, heading)

	// NorthEast
	dataNorthEast := Data{
		Latlng: XDILatLng{Lat: 0, Lng: 0},
		Euler:  XDIEulerAngles{Roll: 0, Pitch: 0, Yaw: 0},
		Mag:    XDIMagneticXYZ{1, 1, 0}, // NorthEast --> Heading == 45
	}

	heading = Data.Heading(dataNorthEast)
	log.Printf("Heading is: %v", heading-magneticDeclination)
	assert.Equal(t, northEast+magneticDeclination, heading)

	// NorthWest
	dataNorthWest := Data{
		Latlng: XDILatLng{Lat: 0, Lng: 0},
		Euler:  XDIEulerAngles{Roll: 0, Pitch: 0, Yaw: 0},
		Mag:    XDIMagneticXYZ{1, -1, 0}, // NorthWest --> Heading == 135
	}

	heading = Data.Heading(dataNorthWest)
	log.Printf("Heading is: %v", heading-magneticDeclination)
	assert.Equal(t, northWest+magneticDeclination, heading)

	// SouthEast
	dataSouthEast := Data{
		Latlng: XDILatLng{Lat: 0, Lng: 0},
		Euler:  XDIEulerAngles{Roll: 0, Pitch: 0, Yaw: 0},
		Mag:    XDIMagneticXYZ{-1, 1, 0}, // SouthEast --> Heading == -45
	}

	heading = Data.Heading(dataSouthEast)
	log.Printf("Heading is: %v", heading-magneticDeclination)
	assert.Equal(t, southEast+magneticDeclination, heading)

	// SouthWest
	dataSouthWest := Data{
		Latlng: XDILatLng{Lat: 0, Lng: 0},
		Euler:  XDIEulerAngles{Roll: 0, Pitch: 0, Yaw: 0},
		Mag:    XDIMagneticXYZ{-1, -1, 0}, // SouthWest --> Heading == -135
	}

	heading = Data.Heading(dataSouthWest)
	log.Printf("Heading is: %v", heading-magneticDeclination)
	assert.Equal(t, southWest+magneticDeclination, heading)

}
