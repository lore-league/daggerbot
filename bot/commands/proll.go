package commands

import (
	"github.com/bwmarrin/discordgo"
)

func proll(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
	return MessagePrivateSend(s, m, parseRoll(c.Args(), "your"))
}

func init() {
	RegisterCommand(NewCommand("PRoll", "Privately replies with Roll!", proll))
}
