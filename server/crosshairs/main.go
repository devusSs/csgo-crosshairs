package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/devusSs/crosshairs/api"
	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/database/postgres"
	"github.com/devusSs/crosshairs/logging"
	"github.com/devusSs/crosshairs/updater"
	"github.com/devusSs/crosshairs/utils"
)

func main() {
	startTime := time.Now()

	printBuild := flag.Bool("v", false, "prints build information")
	cfgPath := flag.String("c", "./files/config.json", "sets config path")
	scFlag := flag.Bool("sc", false, "generates secret keys")
	debugFlag := flag.Bool("d", false, "enabled debug mode")
	flag.Parse()

	if *debugFlag {
		updater.BuildMode = "dev"
	}

	if *printBuild {
		updater.PrintBuildInfo()
		return
	}

	if *scFlag {
		log.Printf("[%s] Sessions secret key: \t%s\n", logging.WarnSign, utils.RandomString(24))
		log.Printf("[%s] Admin key: \t\t%s\n", logging.WarnSign, utils.RandomString(36))
		log.Printf("[%s] Make sure to add these to your config file and NEVER SHARE THEM WITH ANYONE!", logging.WarnSign)
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

	gormLogger, err := logging.CreateGormLogger()
	if err != nil {
		log.Fatalf("[%s] Error creating gorm logger: %s\n", logging.ErrSign, err.Error())
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

	svc, err := postgres.NewConnection(cfg, gormLogger)
	if err != nil {
		logging.WriteError(err.Error())
		os.Exit(1)
	}

	if err := svc.TestConnection(); err != nil {
		logging.WriteError(err.Error())
		os.Exit(1)
	}

	if err := svc.MakeMigrations(); err != nil {
		logging.WriteError(err.Error())
		os.Exit(1)
	}

	apiLogFile, err := logging.CreateAPILogFile()
	if err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	apiServer, err := api.NewAPIInstance(cfg, apiLogFile)
	if err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	if err := apiServer.SetupSessions(cfg); err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	apiServer.SetupRedisRateLimiting(cfg)

	apiServer.SetupCors(cfg)

	apiServer.SetupRoutes(svc, cfg)

	if err := apiServer.StartAPI(); err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	// ! App exit.
	if err := svc.CloseConnection(); err != nil {
		log.Fatalf("[%s] Error closing database connection: %s\n", logging.ErrSign, err.Error())
	}

	if err := logging.CloseLogFiles(); err != nil {
		log.Fatalf("[%s] Error closing log files: %s\n", logging.ErrSign, err.Error())
	}

	log.Printf("[%s] App ran for %.2f second(s)\n", logging.InfSign, time.Since(startTime).Seconds())
}
