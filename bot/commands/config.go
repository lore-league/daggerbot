package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/config"
)

func Config(cmd *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
	gid := m.GuildID
	if gid == "" {
		if _, err := s.ChannelMessageSend(m.ChannelID, "This command can only be used on a server"); err != nil {
			log.Printf("failed sending Config Command response: %v", err)
		}
		return nil
	}

	guild, ok := config.Guilds[gid]
	if !ok { // This shouldn't happen
		return fmt.Errorf("guild %s was not registered before use", gid)
	}

	if !guild.IsAdmin(m.Member) {
		log.Println("User is not an admin, denying access to Config command")
		if _, err := s.ChannelMessageSend(m.ChannelID, "You must be an admin to use this command!"); err != nil {
			log.Printf("Failed sending Config Command response: %v", err)
		}
		return nil
	}

	if len(cmd.Args) < 1 {
		if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: !config <command> [args]\nAvailable commands: `get`, `set`, `list`, `clear`, `help`"); err != nil {
			log.Printf("Failed sending Config Command response: %v", err)
		}
		return nil
	}

	// Handle the config command
	switch cmd.Args[0] {

	case "get":
		if len(cmd.Args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: !config get <key>"); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		key := cmd.Args[1]
		value, exists := guild.Config[key]
		if !exists {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Config key `%s` does not exist", key)); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Config `%s`: %v", key, value)); err != nil {
			log.Printf("Failed sending Config Command response: %v", err)
		}
		return nil

	case "set":
		if len(cmd.Args) < 3 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: !config set <key> <value>"); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		key := cmd.Args[1]
		value := strings.TrimSpace(strings.Join(cmd.Args[2:], " "))
		if guild.Config == nil {
			guild.Config = make(map[string]string)
		}
		if key == "" || value == "" {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Key and value cannot be empty"); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		guild.Config[key] = value
		if err := guild.Save(); err != nil {
			log.Printf("Failed to save config for guild %s: %v", gid, err)
			if _, err := s.ChannelMessageSend(m.ChannelID, "Failed to save config"); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return err
		}
		if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Config `%s` set to `%s`", key, value)); err != nil {
			log.Printf("Failed sending Config Command response: %v", err)
		}
		return nil

	case "list":
		if len(guild.Config) == 0 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "No config variables set for this guild"); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		var response string
		for key, value := range guild.Config {
			response += fmt.Sprintf("`%s: %s`\n", key, value)
		}
		if _, err := s.ChannelMessageSend(m.ChannelID, response); err != nil {
			log.Printf("Failed sending Config Command response: %v", err)
		}
		return nil

	case "clear":
		if len(cmd.Args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: !config clear <key>"); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		key := cmd.Args[1]
		if _, exists := guild.Config[key]; !exists {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Config key `%s` does not exist", key)); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		delete(guild.Config, key)
		if err := guild.Save(); err != nil {
			log.Printf("Failed to save config for guild %s: %v", gid, err)
			if _, err := s.ChannelMessageSend(m.ChannelID, "Failed to save config"); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return err
		}
		if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Config key `%s` cleared", key)); err != nil {
			log.Printf("Failed sending Config Command response: %v", err)
		}
		return nil

	default:
		helpMessage := "Config Command Help:\n" +
			"`!config get <key>` - Retrieves the value of a config key" +
			"`!config set <key> <value>` - Sets a config key to a value" +
			"`!config list` - Lists all config keys and their values" +
			"`!config clear <key>` - Clears a config key" +
			"`!config help` - Displays this help message"
		if _, err := s.ChannelMessageSend(m.ChannelID, helpMessage); err != nil {
			log.Printf("Failed sending Config Command response: %v", err)
		}
		return nil
	}
}

func init() {
	RegisterCommand(NewCommand("Config", "Sets or Retrieves config variables (admin only)", Config))
}
