package main

import (
	"fmt"
	"github.com/einride/xsens-go"
)

func main() {

	// Open the device
	err := xsens.Open()
	if nil != err {
		fmt.Printf("Could not open XSens device: %s\n", err)
		return
	}
	defer xsens.Close()

	xsens.Run(func(data xsens.XsensData, err error) {
		fmt.Printf("got some data %v, %v", data, err)
	})

	c := make(chan struct{})
	<-c
}
