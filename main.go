package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"asdfgh/src/cli"

	log "github.com/schollz/logger"
)

func main() {
	var err error
	// Define the debug flag
	debug := flag.Bool("debug", false, "enable debug mode")

	// Custom flag parsing to capture additional args (files)
	flag.Parse()

	files := flag.Args() // Remaining arguments after flags

	// Output debug mode status
	if *debug {
		log.SetLevel("debug")
	} else {
		log.SetLevel("info")
	}

	log.Debugf("files: %v", files)
	clis := make([]*cli.CLI, len(files))
	for i, file := range files {
		log.Debugf("Opening file: %s", file)
		clis[i], err = cli.Init(file)
		if err != nil {
			log.Errorf("Error opening file: %s", err)
			os.Exit(1)
		} else {
			log.Debugf("%+v", clis[i])
			log.Infof("Parsed %s into %d components", file, len(clis[i].TLI[0].Components))
		}
	}
	// play all of them
	for _, cli := range clis {
		if err := cli.Play(); err != nil {
			log.Errorf("Error playing TLI: %s", err)
			os.Exit(1)
		}
	}

	// wait for ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Debugf("Received interrupt signal, stopping playback")
	for _, cli := range clis {
		cli.Stop()
	}
	time.Sleep(1 * time.Second)
}
