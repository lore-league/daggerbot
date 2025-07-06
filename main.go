package main

import (
	"flag"
	"log"

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
	flag.BoolVar(&config.Debug, "debug", false, "Enable debug mode")
	flag.BoolVar(&config.Verbose, "verbose", false, "Enable verbose mode")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Initialized: Commandline and Environment variables loaded successfully")
}
