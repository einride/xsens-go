package main

import (
	"fmt"
	"github.com/einride/xsens-go"
	"log"
)

func main() {
	defer xsens.Close()
	// Open the device
	err := xsens.Open()
	log.Printf("Could not open XSens device: %s\n", err)

	c := make(chan struct{})
	<-c
}
