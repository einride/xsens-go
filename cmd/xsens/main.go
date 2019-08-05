package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/einride/xsens-go"
	"github.com/einride/xsens-go/pkg/serial"
)

func main() {
	ctx := withCancelOnSignal(context.Background(), os.Interrupt)
	flags := flag.NewFlagSet("xsens", flag.ExitOnError)
	usage := func() {
		fmt.Print(`
usage:

	xsens read <port>
	xsens get-output-config <port> [-json]
	xsens set-ouptut-config <port> <config.json>

`)
		flags.PrintDefaults()
		fmt.Println()
		os.Exit(1)
	}
	flags.Usage = usage
	if len(os.Args) < 2 {
		usage()
	}
	subcommand, args := os.Args[1], os.Args[2:]
	arg := func(i int) string {
		if flags.Arg(i) == "" {
			usage()
		}
		return flags.Arg(i)
	}
	jsonFlag := flags.Bool("json", false, "use JSON output")
	baudRateFlag := flags.Int("baudRate", int(serial.BaudRate115200), "baud rate for serial communication")
	_ = flags.Parse(args)
	switch subcommand {
	case "read":
		if err := readMain(ctx, arg(0), serial.BaudRate(*baudRateFlag), *jsonFlag); err != nil {
			fmt.Println(err)
			usage()
		}
	case "get-output-config":
		if err := getOutputConfigMain(ctx, arg(0), serial.BaudRate(*baudRateFlag), *jsonFlag); err != nil {
			fmt.Println(err)
			usage()
		}
	case "set-output-config":
		if err := setOutputConfigMain(ctx, arg(0), arg(1), serial.BaudRate(*baudRateFlag)); err != nil {
			fmt.Println(err)
			usage()
		}
	default:
		usage()
	}
}

func readMain(_ context.Context, portName string, baudRate serial.BaudRate, useJSON bool) (err error) {
	port, err := serial.Open(portName, baudRate)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := port.Close(); errClose != nil {
			err = errClose
		}
	}()
	if err := port.Flush(); err != nil {
		return err
	}
	sc := xsens.NewMessageScanner(port)
	// install interrupt handler
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	// go to measurement
	if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierGotoMeasurement, nil)); err != nil {
		return err
	}
	var data xsens.MeasurementData
loop:
	for sc.Scan() {
		select {
		case <-interruptChan:
			break loop
		default:
		}
		msg := sc.Message()
		if msg.Identifier() != xsens.MessageIdentifierMTData2 {
			continue
		}
		if err := data.UnmarshalMTData2(msg); err != nil {
			return err
		}
		if useJSON {
			js, err := json.Marshal(data)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", js)
		} else {
			text, err := data.MarshalText()
			if err != nil {
				return err
			}
			fmt.Printf("%s\n\n", text)
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return nil
}

func getOutputConfigMain(_ context.Context, portName string, baudRate serial.BaudRate, useJSON bool) (err error) {
	port, err := serial.Open(portName, baudRate)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := port.Close(); errClose != nil && err == nil {
			err = errClose
		}
	}()
	if err := port.Flush(); err != nil {
		return err
	}
	sc := xsens.NewMessageScanner(port)
	// go to config
	if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierGotoConfig, nil)); err != nil {
		return err
	}
	for sc.Scan() && sc.Message().Identifier() != xsens.MessageIdentifierGotoConfig.Ack() {
		// scan for ack
	}
	if sc.Err() != nil {
		return err
	}
	// request output configuration
	if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierReqOutputConfiguration, nil)); err != nil {
		return err
	}
	for sc.Scan() && sc.Message().Identifier() != xsens.MessageIdentifierReqOutputConfiguration.Ack() {
		// scan for ack
	}
	if sc.Err() != nil {
		return err
	}
	// parse and print output configuration
	var outputConfiguration xsens.OutputConfiguration
	if err := outputConfiguration.Unmarshal(sc.Message().Data()); err != nil {
		return err
	}
	if useJSON {
		js, err := json.Marshal(outputConfiguration)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", js)
	} else {
		txt, err := outputConfiguration.MarshalText()
		if err != nil {
			return err
		}
		fmt.Printf("\n%s\n", txt)
	}
	return nil
}

func setOutputConfigMain(_ context.Context, portName string, jsonFile string, baudRate serial.BaudRate) (err error) {
	port, err := serial.Open(portName, baudRate)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := port.Close(); errClose != nil && err == nil {
			err = errClose
		}
	}()
	if err := port.Flush(); err != nil {
		return err
	}
	sc := xsens.NewMessageScanner(port)
	// parse output configuration
	js, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}
	var outputConfiguration xsens.OutputConfiguration
	if err := json.Unmarshal(js, &outputConfiguration); err != nil {
		return err
	}
	// print output configuration
	txt, err := outputConfiguration.MarshalText()
	if err != nil {
		return err
	}
	fmt.Printf("Setting output configuration:\n\n%s\n", txt)
	// go to config
	if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierGotoConfig, nil)); err != nil {
		return err
	}
	for sc.Scan() && sc.Message().Identifier() != xsens.MessageIdentifierGotoConfig.Ack() {
		// scan for ack
	}
	if sc.Err() != nil {
		return err
	}
	// set output configuration
	data, err := outputConfiguration.Marshal()
	if err != nil {
		return err
	}
	if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierSetOutputConfiguration, data)); err != nil {
		return err
	}
	for sc.Scan() && sc.Message().Identifier() != xsens.MessageIdentifierSetOutputConfiguration.Ack() {
		// scan for ack
	}
	return nil
}

func withCancelOnSignal(ctx context.Context, sig ...os.Signal) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, len(sig))
	signal.Notify(signalChan, sig...)
	go func() {
		<-signalChan
		signal.Stop(signalChan)
		cancel()
	}()
	return ctx
}
