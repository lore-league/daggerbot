package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if _, err := s.ChannelMessageSend(m.ChannelID, "Pong!"); err != nil {
		return fmt.Errorf("failed to send Pong response: %w", err)
	}
	return nil
}
