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
	} else {
		defer xsens.Close()
	}

	c := make(chan struct{})
	<-c
}
