package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Commands = map[string]*Command{}

type Handler func(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error

type Command struct {
	Name        string   // Name of the command
	Description string   // Description of the command
	Args        []string // Arguments for the command
	Handler     Handler
}

func (c Command) String() string {
	return fmt.Sprintf("%s: %s", c.Name, c.Description)
}

func RegisterCommand(command *Command) {
	name := strings.ToLower(command.Name)

	if _, exists := Commands[name]; exists {
		fmt.Printf("Command %s already exists, not registering again.\n", name)
		return
	}

	Commands[name] = command
	fmt.Printf("Registered command: %s\n", name)
}

func NewCommand(name, description string, handler Handler) *Command {
	return &Command{
		Name:        name,
		Description: description,
		Args:        make([]string, 0),
		Handler:     handler,
	}
}
