package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/alecthomas/kingpin"
	"github.com/einride/xsens-go"
	"github.com/tarm/serial"
)

func main() {
	// define command line app
	app := kingpin.New("xsens", "Xsens CLI tool.")
	// read subcommand
	readCmd := app.
		Command("read", "Read MTData2 measurement data.")
	readPortArg := readCmd.
		Arg("port", "Xsens serial port to connect to.").Required().ExistingFile()
	readJSONArg := readCmd.
		Flag("json", "Use JSON output format.").Bool()
	// get-output-config subcommand
	getOutputConfigCmd := app.
		Command("get-output-config", "Read out the current output config from the Xsens.")
	getOutputConfigPortArg := getOutputConfigCmd.
		Arg("port", "Xsens serial port to connect to.").Required().ExistingFile()
	getOutputConfigJSONFlag := getOutputConfigCmd.
		Flag("json", "Use JSON output format.").Bool()
	// set-output-config subcommand
	setOutputConfigCmd := app.
		Command("set-output-config", "Read out the current output config from the Xsens.")
	setOutputConfigPortArg := setOutputConfigCmd.
		Arg("port", "Xsens serial port to connect to.").Required().ExistingFile()
	setOutputConfigJSONFileArg := setOutputConfigCmd.
		Arg("jsonFile", "JSON output configuration to set.").Required().ExistingFile()
	// run command line app
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case readCmd.FullCommand():
		readMain(*readPortArg, *readJSONArg)
	case getOutputConfigCmd.FullCommand():
		getOutputConfigMain(*getOutputConfigPortArg, *getOutputConfigJSONFlag)
	case setOutputConfigCmd.FullCommand():
		setOutputConfigMain(*setOutputConfigPortArg, *setOutputConfigJSONFileArg)
	}
}

func readMain(portName string, useJSON bool) {
	port, err := serial.OpenPort(&serial.Config{
		Name:     portName,
		Baud:     xsens.DefaultSerialBaudRate,
		Size:     xsens.MinLengthOfMessage,
		StopBits: xsens.DefaultSerialStopBits,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := port.Close(); err != nil {
			panic(err)
		}
	}()
	if err := port.Flush(); err != nil {
		panic(err)
	}
	sc := xsens.NewMessageScanner(port)
	// install interrupt handler
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	// go to measurement
	if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierGotoMeasurement, nil)); err != nil {
		panic(err)
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
			panic(err)
		}
		if useJSON {
			js, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", js)
		} else {
			text, err := data.MarshalText()
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n\n", text)
		}
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
}

func getOutputConfigMain(portName string, useJSON bool) {
	port, err := serial.OpenPort(&serial.Config{
		Name:     portName,
		Baud:     xsens.DefaultSerialBaudRate,
		Size:     xsens.MinLengthOfMessage,
		StopBits: xsens.DefaultSerialStopBits,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := port.Close(); err != nil {
			panic(err)
		}
	}()
	if err := port.Flush(); err != nil {
		panic(err)
	}
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
	// request output configuration
	if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierReqOutputConfiguration, nil)); err != nil {
		panic(err)
	}
	for sc.Scan() && sc.Message().Identifier() != xsens.MessageIdentifierReqOutputConfiguration.Ack() {
		// scan for ack
	}
	if sc.Err() != nil {
		panic(err)
	}
	// parse and print output configuration
	var outputConfiguration xsens.OutputConfiguration
	if err := outputConfiguration.Unmarshal(sc.Message().Data()); err != nil {
		panic(err)
	}
	if useJSON {
		js, err := json.Marshal(outputConfiguration)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", js)
	} else {
		txt, err := outputConfiguration.MarshalText()
		if err != nil {
			panic(err)
		}
		fmt.Printf("\n%s\n", txt)
	}
}

func setOutputConfigMain(portName string, jsonFile string) {
	port, err := serial.OpenPort(&serial.Config{
		Name:     portName,
		Baud:     xsens.DefaultSerialBaudRate,
		Size:     xsens.MinLengthOfMessage,
		StopBits: xsens.DefaultSerialStopBits,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := port.Close(); err != nil {
			panic(err)
		}
	}()
	if err := port.Flush(); err != nil {
		panic(err)
	}
	sc := xsens.NewMessageScanner(port)
	// parse output configuration
	js, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}
	var outputConfiguration xsens.OutputConfiguration
	if err := json.Unmarshal(js, &outputConfiguration); err != nil {
		panic(err)
	}
	// print output configuration
	txt, err := outputConfiguration.MarshalText()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Setting output configuration:\n\n%s\n", txt)
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
	// set output configuration
	data, err := outputConfiguration.Marshal()
	if err != nil {
		panic(err)
	}
	if _, err := port.Write(xsens.NewMessage(xsens.MessageIdentifierSetOutputConfiguration, data)); err != nil {
		panic(err)
	}
	for sc.Scan() && sc.Message().Identifier() != xsens.MessageIdentifierSetOutputConfiguration.Ack() {
		// scan for ack
	}
}
