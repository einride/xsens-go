// +build Client

package xsensgo

import (
	"log"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests are used online, which requires a connected Client.
// Used for debugging and some verification.

func TestReadmsgs(t *testing.T) {
	test := assert.New(t)
	prt, err := NewClient()
	test.Nil(err)
	defer prt.Close()

	err = prt.readMessages(func(data XsensData, err error) {
		log.Printf("got this %v, %v", data, err)
	})
	test.Nil(err)
}

func TestRun(t *testing.T) {
	test := assert.New(t)
	x, err := NewClient()
	test.Nil(err)
	err = x.Run(func(data XsensData, err error) {
		log.Printf("Got this %v, %v", data, err)
	})
	test.Nil(err)
}

func TestHeadingCalc(t *testing.T) {
	test := assert.New(t)
	x, err := NewClient()
	test.Nil(err)
	err = x.Run(func(data XsensData, err error) {
		heading := data.Heading()
		log.Printf("Heading is %v", heading)
		// Set north as reference
		wanted := math.Atan2(1, 0) * 180 / math.Pi
		log.Printf("wanted: %v", wanted)
		// Check angle error (reference - actual)
		headingError := wanted - heading
		switch {
		case headingError > 180.0:
			headingError -= 360.0
		case headingError < -180.0:
			headingError += 360
		}
		// When Client is facing north, headingError should be 0
		log.Printf("Heading error is %v", headingError)
	})
}
