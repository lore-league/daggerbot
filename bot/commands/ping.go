package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Ping(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
	if _, err := s.ChannelMessageSend(m.ChannelID, "Pong!"); err != nil {
		return fmt.Errorf("failed to send Pong response: %w", err)
	}
	return nil
}

func init() {
	RegisterCommand(NewCommand("Ping", "Replies with Pong!", Ping))
}
