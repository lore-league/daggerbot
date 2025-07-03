package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string   // Name of the command
	Description string   // Description of the command
	Args        []string // Arguments for the command
	Handler     func(s *discordgo.Session, m *discordgo.MessageCreate) error
}

func (c Command) String() string {
	return fmt.Sprintf("%s: %s", c.Name, c.Description)
}

var Commands = map[string]Command{
	"ping": {
		Name:        "Ping",
		Description: "Replies with Pong!",
		Handler:     Ping,
	},
}
