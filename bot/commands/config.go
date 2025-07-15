package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/config"
)

func Config(cmd *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
	var (
		args  = cmd.Args()
		guild = cmd.Guild()
	)

	if len(args) < 1 {
		MessageSend(s, m, "Usage: !config <command> [args]\nAvailable commands: `get`, `set`, `list`, `clear`, `help`")
		return nil
	}

	// Handle the config command
	switch args[0] {

	case "get":
		if len(args) < 2 {
			MessageSend(s, m, "Usage: !config get <key>")
			return nil
		}
		key := args[1]
		value, exists := guild.GetConfigMap()[key]
		if !exists {
			MessageSend(s, m, fmt.Sprintf("Config key `%s` does not exist", key))
			return nil
		}
		MessageSend(s, m, fmt.Sprintf("Config `%s`: %v", key, value))
		return nil

	case "set":
		if len(args) < 3 {
			MessageSend(s, m, "Usage: !config set <key> <value>")
			return nil
		}
		key := args[1]
		values := strings.Split(strings.Join(args[2:], ","), ",")
		if key == "" || len(values) == 0 {
			MessageSend(s, m, "Key and value cannot be empty")
			return nil
		}
		cmap := guild.GetConfigMap()
		if _, ok := cmap[key]; !ok {
			MessageSend(s, m, fmt.Sprintf("Key must exist and must be one of %v", config.ConfigKeys))
			return nil
		}

		// Special handling for prefix
		if strings.TrimSpace(strings.ToLower(key)) == "prefix" {
			if len(values) != 1 || len(values[0]) == 0 {
				MessageSend(s, m, "Prefix must be a single non-empty string")
				return nil
			}
			guild.SetPrefix(strings.TrimSpace(values[0]))
			MessageSend(s, m, fmt.Sprintf("Prefix set to `%s`", guild.Prefix()))
			return nil
		}

		// Handle roles
		newValues := make([]*discordgo.Role, 0, len(values))
		for _, v := range values {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			// Check to see if the value is a valid ID or a role name
			if role := guild.FindRole(v); role != nil {
				newValues = append(newValues, role) // Use the role ID if a role is found
			} else {
				MessageSend(s, m, fmt.Sprintf("Could not find a role associated with %s", v))
				continue
			}
		}

		if err := guild.SetRoleConfig(key, newValues); err != nil {
			log.Printf("Failed to set config for guild %q: %v", guild.Name, err)
			MessageSend(s, m, "Failed to set config")
			return err
		}

		MessageSend(s, m, fmt.Sprintf("Config `%s` set to `%s`", key, strings.Join(guild.RolesToNames(newValues), ", ")))
		return nil

	case "list":
		var response string

		for _, key := range config.ConfigKeys {
			if strings.ToLower(strings.TrimSpace(key)) == "prefix" {
				response += fmt.Sprintf("`%s: %s`\n", key, guild.Prefix())
			} else {
				roles, err := guild.GetRoleConfig(key)
				if err != nil {
					log.Printf("Error retrieving roles for key %s: %v", key, err)
					response += fmt.Sprintf("`%s: (error retrieving roles)`\n", key)
					continue
				}
				if len(roles) == 0 {
					response += fmt.Sprintf("`%s: (no roles assigned)`\n", key)
					continue
				}
				response += fmt.Sprintf("`%s: %s`\n", key, strings.Join(guild.RolesToNames(roles), ", "))
			}
		}

		MessageSend(s, m, response)
		return nil

	case "clear":
		if len(args) < 2 {
			MessageSend(s, m, "Usage: !config clear <key>")
			return nil
		}
		key := args[1]
		if _, exists := guild.GetConfigMap()[key]; !exists {
			MessageSend(s, m, fmt.Sprintf("Config key `%s` does not exist", key))
			return nil
		}
		if key == "prefix" {
			MessageSend(s, m, "You cannot clear the `prefix` key")
			return nil
		}
		if err := guild.ClearConfig(key); err != nil {
			log.Printf("Failed to clear config for guild %q: %v", guild.Name, err)
			MessageSend(s, m, "Failed to clear config")
			return err
		}
		MessageSend(s, m, fmt.Sprintf("Config key `%s` cleared", key))
		return nil

	default:
		keys := make([]string, 0, len(config.ConfigKeys))
		for _, key := range config.ConfigKeys {
			keys = append(keys, fmt.Sprintf("`%s`", strings.TrimSpace(key)))
		}
		MessageSend(s, m, "Config Command Help:\n"+
			"`!config get <key>` - Retrieves the value of a config key\n"+
			"`!config set <key> <value>` - Sets a config key to a value\n"+
			"`!config list` - Lists all config keys and their values\n"+
			"`!config clear <key>` - Clears a config key\n"+
			"`!config help` - Displays this help message"+
			"\n\nAvailable config keys: "+strings.Join(keys, ", "))
		return nil
	}
}

func init() {
	cmd := NewCommand("Config", "Sets or retrieves config variables", Config)
	cmd.SetAdmin()
	RegisterCommand(cmd)
}
