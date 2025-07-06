package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/config"
)

func Version(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
	if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Daggerbot Version: %s", config.Version)); err != nil {
		return fmt.Errorf("failed to send Version response: %w", err)
	}
	return nil
}

func init() {
	RegisterCommand(NewCommand("Version", "Replies with bot version information!", Version))
}
