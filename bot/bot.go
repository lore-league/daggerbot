package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/bot/handlers"
	"github.com/nerdwerx/daggerbot/config"
)

func startup() {
	// Load configuration from environment variables or other sources
	config.Token = os.Getenv("DISCORD_AUTH_TOKEN")
	if config.Token == "" {
		log.Fatal("DISCORD_AUTH_TOKEN is not set")
	}
	config.Prefix = os.Getenv("DISCORD_BOT_PREFIX")
	if config.Prefix == "" {
		config.Prefix = "!" // Default prefix
	}
	log.Printf("Bot Prefix set to %q", config.Prefix)
	log.Printf("Configuration loaded")
}

func Run() error {
	// Initialize our config vars
	startup()

	// create a session
	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	if config.Debug {
		log.Println("Debug mode enabled")
		// discord.Debug = true
	}

	// add a event handlers
	discord.AddHandler(handlers.OnReady)
	discord.AddHandler(handlers.OnMessage)

	// Set our permissions
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	log.Println("Bot starting up... CTRL-C to stop")

	// open session
	if err := discord.Open(); err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	defer func() {
		if err := discord.Close(); err != nil { // close session, after function termination
			log.Fatal("Error closing Discord session")
		}
	}()

	// keep bot running untill we're interrupted (ctrl + C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Bot gracefully shutting down...")

	return nil
}
