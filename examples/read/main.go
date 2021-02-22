package main

import (
	"context"
	"flag"
	"log"
	"os"

	"go.einride.tech/xsens"
	"go.einride.tech/xsens/serial"
)

func main() {
	ctx := context.Background()
	log.SetFlags(0)
	port := flag.String("port", "", "serial port to read from")
	baudRateFlag := flag.Int("baudRate", int(serial.BaudRate115200), "baud rate for serial port")
	flag.Parse()
	if *port == "" {
		flag.Usage()
		os.Exit(1)
	}
	// Open serial port.
	serialPort, err := serial.Open(*port, serial.BaudRate(*baudRateFlag))
	if err != nil {
		log.Fatal(err)
	}
	client := xsens.NewClient(serialPort)
	// Perform GoToMeasurement sequence.
	if err := client.GoToMeasurement(ctx); err != nil {
		log.Panic(err)
	}
	for {
		// Scan through all packets in the current MTData2 message.
		log.Println(client.MessageIdentifier())
		for client.ScanMeasurementData() {
			log.Printf("\t%v", client.DataType())
			log.Printf("\t%+v", client.MeasurementData())
		}
		// Receive next MTData2 message.
		if err := client.Receive(ctx); err != nil {
			log.Panic(err)
		}
	}
}
