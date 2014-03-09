package asmiostat

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type optionset struct {
	RunWebServer bool
	IostatArgs   []string
}

func runStdout(options optionset) {
	//fmt.Println("-Started")
	deverrors := make(chan error)
	if err := startBackgroundSubs(deverrors); err != nil {
		fmt.Printf("Error getting device mappings: %s\n", err)
		os.Exit(-2)
	}

	iostat := make(chan string)
	iostaterrors := make(chan error)
	if err := startIostat(options, iostat, iostaterrors); err != nil {
		fmt.Printf("Error starting iostat: %s\n", err)
		os.Exit(-3)
	}

	startStdoutput(iostat)

	//fmt.Println("-initialization done. waiting for any errors")
	select {
	case err := <-deverrors:
		fmt.Printf("Error getting device mappings: %s\n", err)
		os.Exit(-2)
	case err := <-iostaterrors:
		if err != io.EOF {
			fmt.Printf("Error running iostat: %s\n", err)
			os.Exit(-3)
		} else {
			os.Exit(0)
		}
	}
	fmt.Printf("Done\n")
}

func Main() {
	options := optionset{false, []string{}}

	f := flag.NewFlagSet("asmiostat", flag.ExitOnError)
	f.BoolVar(&options.RunWebServer, "w", false, "Runs as a webserver instead of stdout")
	f.Parse(os.Args[1:])
	options.IostatArgs = f.Args()
	//fmt.Printf("-iostat args: %v\n", options.IostatArgs)

	//startBackgroundProcess()

	if options.RunWebServer {
		fmt.Printf("Not implemented yet\n")
	} else {
		runStdout(options)
	}
}
