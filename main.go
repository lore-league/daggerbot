package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nerdwerx/daggerbot/bot"
	"github.com/nerdwerx/daggerbot/config"
)

func main() {
	if err := bot.Run(); err != nil {
		log.Fatal("Error starting bot:", err)
	}
}

func init() {
	var healthcheck bool

	flag.BoolVar(&config.Debug, "debug", false, "Enable debug mode")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose mode")
	flag.BoolVar(&healthcheck, "healthcheck", false, "Run health check")
	flag.Parse()

	if config.Debug {
		log.Println("Debug mode enabled")
		config.Verbose = true
	}

	if config.Verbose {
		log.Println("Verbose mode enabled")
	}

	if healthcheck {
		if config.Verbose {
			log.Println("health check: OK")
		}
		os.Exit(0)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Initialized: Commandline and Environment variables loaded successfully")
}
