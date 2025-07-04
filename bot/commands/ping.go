package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if _, err := s.ChannelMessageSend(m.ChannelID, "Pong!"); err != nil {
		return fmt.Errorf("failed to send Pong response: %w", err)
	}
	return nil
}

func init() {
	if _, exists := Commands["ping"]; exists {
		return
	} else {
		log.Println("Registering Ping command")
		Commands["ping"] = Command{
			Name:        "Ping",
			Description: "Replies with Pong!",
			Handler:     Ping,
		}
	}
}
