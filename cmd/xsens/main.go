package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"go.bug.st/serial"
	"go.einride.tech/xsens"
	"golang.org/x/sync/errgroup"
)

const DefaultBaudRate = 115200

func main() {
	ctx := withCancelOnSignal(context.Background(), os.Interrupt)
	flags := flag.NewFlagSet("xsens", flag.ExitOnError)
	jsonFlag := flags.Bool("json", false, "use JSON output")
	baudRateFlag := flags.Int("baudRate", DefaultBaudRate, "baud rate for serial communication")
	configTimeoutFlag := flags.Duration("configTimeout", time.Second, "timeout for config operations")
	usage := func() {
		fmt.Print(`
usage:

	xsens read [-baudRate <int>] <port>
	xsens get-output-config [-baudRate <int>] [-json] [-configTimeout <duration>] <port>
	xsens set-ouptut-config [-baudRate <int>] [-configTimeout <duration>] <port> <config.json>

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
	_ = flags.Parse(args)
	portName := arg(0)
	port, err := serial.Open(portName, &serial.Mode{BaudRate: *baudRateFlag})
	if err != nil {
		fmt.Println(err)
		usage()
	}
	client := xsens.NewClient(port)
	g, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithCancel(ctx)
	g.Go(func() error {
		<-ctx.Done()
		return client.Close()
	})
	switch subcommand {
	case "read":
		g.Go(func() error {
			defer cancel()
			return readMain(ctx, client)
		})
	case "get-output-config":
		g.Go(func() error {
			defer cancel()
			return getOutputConfigMain(ctx, client, *configTimeoutFlag, *jsonFlag)
		})
	case "set-output-config":
		g.Go(func() error {
			defer cancel()
			return setOutputConfigMain(ctx, client, arg(1), *configTimeoutFlag)
		})
	default:
		usage()
	}
	if err := g.Wait(); err != nil && !strings.Contains(err.Error(), "closed") {
		fmt.Println(err)
		usage()
	}
}

func readMain(ctx context.Context, client *xsens.Client) error {
	if err := client.GoToMeasurement(ctx); err != nil {
		return err
	}
	for {
		fmt.Println()
		fmt.Println(client.MessageIdentifier())
		for client.ScanMeasurementData() {
			fmt.Printf("\t%v\n", client.DataType())
			fmt.Printf("\t%+v\n", client.MeasurementData())
		}
		if err := client.Receive(ctx); err != nil {
			if strings.Contains(err.Error(), "closed") {
				return nil
			}
			return err
		}
	}
}

func getOutputConfigMain(ctx context.Context, client *xsens.Client, timeout time.Duration, useJSON bool) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	if err := client.GoToConfig(ctx); err != nil {
		return err
	}
	outputConfiguration, err := client.GetOutputConfiguration(ctx)
	if err != nil {
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

func setOutputConfigMain(ctx context.Context, client *xsens.Client, jsonFile string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	// parse output configuration
	js, err := os.ReadFile(jsonFile)
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
	if err := client.GoToConfig(ctx); err != nil {
		return err
	}
	return client.SetOutputConfiguration(ctx, outputConfiguration)
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
