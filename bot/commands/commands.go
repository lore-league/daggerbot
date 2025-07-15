package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/config"
)

var Commands = map[string]*Command{}

type Handler func(c *Command, s *discordgo.Session, m *discordgo.MessageCreate) error

type Data struct {
	admin   bool          // Whether the command is admin-only
	args    []string      // Arguments for the command
	guild   *config.Guild // Guild this command is registered for (optional)
	handler Handler       // Function to handle the command
}

type Command struct {
	Name        string // Name of the command
	Description string // Description of the command
	data        Data
}

func (c *Command) String() string {
	return fmt.Sprintf("%s: %s", c.Name, c.Description)
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

func (c *Command) Guild() *config.Guild {
	return c.data.guild
}

func (c *Command) Run(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if c.data.handler == nil {
		return fmt.Errorf("no handler defined for command %s", c.Name)
	}
	return c.data.handler(c, s, m)
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

// RegisterCommand registers a new command in the global commands map
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
		data: Data{
			admin:   false, // Default to non-admin
			args:    make([]string, 0),
			handler: handler,
		},
	}
}

func MessageSend(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	if len(message) > 2000 {
		log.Printf("message exceeds Discord's 2000 character limit")
	}
	if _, err := s.ChannelMessageSend(m.ChannelID, message); err != nil {
		log.Printf("failed to send message: %s", err.Error())
	}
	if config.Debug {
		log.Printf("Sent message to channel %s: %s", m.ChannelID, message)
	}
}
