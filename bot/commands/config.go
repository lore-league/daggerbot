package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Config(cmd *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
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
		value, exists := cmd.Guild.Config[key]
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
		if cmd.Guild.Config == nil {
			cmd.Guild.Config = make(map[string]string)
		}
		if key == "" || value == "" {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Key and value cannot be empty"); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		cmd.Guild.Config[key] = value
		if err := cmd.Guild.Save(); err != nil {
			log.Printf("Failed to save config for guild %q: %v", cmd.Guild.Name, err)
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
		if len(cmd.Guild.Config) == 0 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "No config variables set for this guild"); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		var response string
		for key, value := range cmd.Guild.Config {
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
		if _, exists := cmd.Guild.Config[key]; !exists {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Config key `%s` does not exist", key)); err != nil {
				log.Printf("Failed sending Config Command response: %v", err)
			}
			return nil
		}
		delete(cmd.Guild.Config, key)
		if err := cmd.Guild.Save(); err != nil {
			log.Printf("Failed to save config for guild %q: %v", cmd.Guild.Name, err)
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
			"`!config get <key>` - Retrieves the value of a config key\n" +
			"`!config set <key> <value>` - Sets a config key to a value\n" +
			"`!config list` - Lists all config keys and their values\n" +
			"`!config clear <key>` - Clears a config key\n" +
			"`!config help` - Displays this help message"
		if _, err := s.ChannelMessageSend(m.ChannelID, helpMessage); err != nil {
			log.Printf("Failed sending Config Command response: %v", err)
		}
		return nil
	}
}

func init() {
	cmd := NewCommand("Config", "Sets or retrieves config variables", Config)
	cmd.Admin = true
	RegisterCommand(cmd)
}
