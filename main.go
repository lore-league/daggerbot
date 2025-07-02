package daggerbot

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Create a new Discord session using the provided bot token.
	// Replace "authentication token" with your actual bot token.
	// Make sure to keep your token secure and not share it publicly.
	discord, err := discordgo.New("Bot " + "authentication token")
	if err != nil {
		// If there is an error creating the session, log it and exit.
		println("error creating Discord session,", err.Error())
		os.Exit(1)
		return
	}

	// placeholder for future functionality
	_ = discord
}
