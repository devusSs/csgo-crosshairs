package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/devusSs/crosshairs/api"
	"github.com/devusSs/crosshairs/api/integration"
	"github.com/devusSs/crosshairs/api/middleware"
	"github.com/devusSs/crosshairs/api/routes"
	"github.com/devusSs/crosshairs/config"
	"github.com/devusSs/crosshairs/database/postgres"
	"github.com/devusSs/crosshairs/logging"
	"github.com/devusSs/crosshairs/storage"
	"github.com/devusSs/crosshairs/updater"
	"github.com/devusSs/crosshairs/utils"
)

func main() {
	startTime := time.Now()

	printBuild := flag.Bool("v", false, "prints build information")
	cfgPath := flag.String("c", "./files/config.json", "sets config path")
	scFlag := flag.Bool("sc", false, "generates secret keys")
	debugFlag := flag.Bool("d", false, "enabled debug mode")
	dockerFlag := flag.Bool("docker", false, "enables Docker mode - uses docker.env instead of config.json file")
	disableIntegrationsFlag := flag.Bool("disable-integrations", false, "disables integrations like Twitch")
	flag.Parse()

	if !checkNetworkConnection() {
		log.Printf("[%s] No working network connection found, exiting\n", logging.ErrSign)
		os.Exit(1)
	}

	if *debugFlag {
		updater.BuildMode = "dev"

		log.Printf("[%s] You are currently using a development version\n", logging.WarnSign)
		log.Printf("[%s] Not all features may be available or working as expected\n", logging.WarnSign)

		updater.PrintBuildInfo()

		frontendUp, err := checkClientConnections()
		if err != nil {
			log.Fatalf("[%s] Error checking frontend clients: %s\n", logging.ErrSign, err.Error())
		}

		if !frontendUp {
			log.Printf("[%s] Could not find a working frontend client\n", logging.WarnSign)
		}
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

	if *disableIntegrationsFlag {
		log.Printf("[%s] You disabled integrations. Not all features may be available\n", logging.WarnSign)
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
	var cfg *config.Config

	if *dockerFlag {
		cfg, err = config.LoadEnvConfig()
		if err != nil {
			logging.WriteError(err.Error())
			os.Exit(1)
		}

		updater.PrintBuildInfo()
	} else {
		cfg, err = config.LoadConfig(*cfgPath)
		if err != nil {
			logging.WriteError(err.Error())
			os.Exit(1)
		}
	}

	if err := cfg.CheckConfig(); err != nil {
		logging.WriteError(err.Error())
		os.Exit(1)
	}

	printAsciiArt(cfg)

	api.UsingReverseProxy = cfg.UsingReverseProxy
	routes.UsingReverseProxy = cfg.UsingReverseProxy

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

	storageSvc, err := storage.NewMinioConnection(cfg)
	if err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	if err := storageSvc.CreateUserProfilePicturesBucket(); err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	if err := storageSvc.UpdateUserProfilePicture("sample.png", "./files/sample.png"); err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	_, err = storageSvc.GetUserProfilePictureLink("sample")
	if err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	// Add database.Service to middleware.
	middleware.Svc = svc

	// Generate an engineer token on startup.
	engineerToken := utils.RandomString(48)

	// Add engineer token to database.
	if err := svc.CreateNewEngineerToken(engineerToken); err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	logging.WriteInfo("Generated initial engineer token and added it to database.")

	// Setup goroutine to add another engineer token every 10 mins.
	generateNewEngineerTokenTicker := time.NewTicker(10 * time.Minute)
	go func() {
		for range generateNewEngineerTokenTicker.C {
			engineerToken := utils.RandomString(48)

			if err := svc.CreateNewEngineerToken(engineerToken); err != nil {
				logging.WriteError(err)
				os.Exit(1)
			}

			logging.WriteInfo("Generated new engineer token and added it to database")
		}
	}()

	logging.WriteInfo("Setup goroutine to re-generate engineer token")

	apiLogFile, errorLogFile, err := logging.CreateAPILogFiles()
	if err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	apiServer, err := api.NewAPIInstance(cfg, apiLogFile, errorLogFile)
	if err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	apiServer.SetupCors(cfg)

	if err := apiServer.SetupSessions(cfg); err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	if err := apiServer.SetupRedisRateLimiting(cfg); err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	apiServer.SetupRoutes(svc, storageSvc, cfg)

	// Integration initialisation
	if !*disableIntegrationsFlag {
		if err := integration.InitTwitchAuth(cfg, apiServer, fmt.Sprintf("http://%s:%d", apiServer.Host, apiServer.Port), svc); err != nil {
			logging.WriteError(err)
			os.Exit(1)
		}

		logging.WriteSuccess("Initialised Twitch authentication")
	}

	if err := apiServer.StartAPI(); err != nil {
		logging.WriteError(err)
		os.Exit(1)
	}

	// ! App exit.
	generateNewEngineerTokenTicker.Stop()

	if err := svc.CloseConnection(); err != nil {
		log.Fatalf("[%s] Error closing database connection: %s\n", logging.ErrSign, err.Error())
	}

	if err := logging.CloseLogFiles(); err != nil {
		log.Fatalf("[%s] Error closing log files: %s\n", logging.ErrSign, err.Error())
	}

	log.Printf("[%s] App ran for %.2f second(s)\n", logging.InfSign, time.Since(startTime).Seconds())
}

// Will print the domain set on config as ascii art.
func printAsciiArt(cfg *config.Config) {
	var printDomain string

	if strings.Contains(cfg.Domain, "http://") {
		printDomain = strings.Replace(cfg.Domain, "http://", "", 1)
	} else {
		printDomain = strings.Replace(cfg.Domain, "https://", "", 1)
	}

	printDomain = strings.Replace(printDomain, "api.", "", 1)

	asciiArt := figure.NewColorFigure(printDomain, "", "green", true)
	asciiArt.Print()
}

// Check if we have a working frontend.
// This might be inaccurate because we are only checking ports but it is a good indicator.
func checkClientConnections() (bool, error) {
	// Ports for React, NextJS, Vite
	potentialClientPorts := []int{3000, 5173}
	clientsOnline := 0

	httpClient := http.Client{
		Timeout: time.Second * 2,
	}

	for _, port := range potentialClientPorts {
		url := fmt.Sprintf("http://localhost:%d/", port)

		response, err := httpClient.Get(url)
		if err != nil {
			if strings.Contains(err.Error(), "connect: connection refused") {
				continue
			}

			return false, err
		}
		defer response.Body.Close()

		// No error indicates there is a working application on the port.
		// Even if we might not get a 200 status code.
		clientsOnline++

		log.Printf("[%s] Got potential working (frontend) client on port %d\n", logging.SucSign, port)

	}

	if clientsOnline == 0 {
		return false, nil
	}

	return true, nil
}

// Checks general network / internet connection.
// If we do not get a response we will exit since an API needs internet connection (duuh)
func checkNetworkConnection() bool {
	_, err := http.Get("http://clients3.google.com/generate_204")
	return err == nil
}
