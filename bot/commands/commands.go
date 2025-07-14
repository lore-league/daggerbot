package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/config"
)

var Commands = map[string]*Command{}

type Handler func(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error

type Data struct {
	admin       bool          // Whether the command is admin-only
	args        []string      // Arguments for the command
	description string        // Description of the command
	guild       *config.Guild // Guild this command is registered for (optional)
	name        string        // Name of the command
}

type Command struct {
	data    Data
	handler Handler // Function to handle the command
}

func (c *Command) String() string {
	return fmt.Sprintf("%s: %s", c.data.name, c.data.description)
}

func (c *Command) Admin() bool {
	return c.data.admin
}

func (c *Command) Args() []string {
	if c.data.args == nil {
		return []string{}
	}

	return c.data.args
}

func (c *Command) Description() string {
	return c.data.description
}

func (c *Command) Guild() *config.Guild {
	return c.data.guild
}

func (c *Command) Run(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if c.handler == nil {
		return fmt.Errorf("no handler defined for command %s", c.data.name)
	}
	return c.handler(c, s, m)
}

func (c *Command) Name() string {
	return c.data.name
}

func (c *Command) SetAdmin() {
	c.data.admin = true
}

func (c *Command) SetArgs(args []string) {
	c.data.args = args
}

func (c *Command) SetGuild(guild *config.Guild) {
	c.data.guild = guild
}

func (c *Command) SetHandler(handler Handler) {
	c.handler = handler
}

// RegisterCommand registers a new command in the global commands map
func RegisterCommand(command *Command) {
	name := strings.ToLower(command.Name())

	if _, exists := Commands[name]; exists {
		fmt.Printf("Command %s already exists, not registering again.\n", name)
		return
	}

	Commands[name] = command
	fmt.Printf("Registered command: %s\n", name)
}

func NewCommand(name, description string, handler Handler) *Command {
	return &Command{
		data: Data{
			name:        name,
			description: description,
			args:        make([]string, 0),
			admin:       false, // Default to non-admin
		},
		handler: handler,
	}
}
