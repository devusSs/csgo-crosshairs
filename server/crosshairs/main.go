package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/logging"
	"github.com/devusSs/crosshairs/updater"
)

func main() {
	startTime := time.Now()

	printBuild := flag.Bool("v", false, "prints build information")
	cfgPath := flag.String("c", "./files/config.json", "sets config path")
	flag.Parse()

	if *printBuild {
		updater.PrintBuildInfo()
		return
	}

	if err := logging.CreateDefaultLogsDirectory(); err != nil {
		log.Fatalf("[%s] Error creating logs directory: %s\n", logging.ErrSign, err.Error())
	}

	if err := logging.CreateAppLogFile(); err != nil {
		log.Fatalf("[%s] Error creating app logger: %s\n", logging.ErrSign, err.Error())
	}

	if err := logging.CreateErrorLogFile(); err != nil {
		log.Fatalf("[%s] Error creating error logger: %s\n", logging.ErrSign, err.Error())
	}

	// ! It is safe to use WriteX methods from here.

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		logging.WriteError(err.Error())
		os.Exit(1)
	}

	if err := cfg.CheckConfig(); err != nil {
		logging.WriteError(err.Error())
		os.Exit(1)
	}

	// ! App exit.
	if err := logging.CloseLogFiles(); err != nil {
		log.Fatalf("[%s] Error closing log files: %s\n", logging.ErrSign, err.Error())
	}

	log.Printf("[%s] App ran for %.2f second(s)\n", logging.InfSign, time.Since(startTime).Seconds())
}
