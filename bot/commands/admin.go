package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/config"
)

func Admin(cmd *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
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
		if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: !admin <command> [args]\nAvailable commands: `list`, `add <role_id>`, `remove <role_id>`"); err != nil {
			log.Printf("Failed sending Admin Command response: %v", err)
		}
		return nil
	}

	// Handle the config command
	switch cmd.Args[0] {
	case "list":
		if len(cmd.Args) > 1 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: !admin list"); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		// List all admin roles
		if len(guild.Admins) == 0 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "No admin roles configured for this guild."); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		adminList := "Admin roles for this guild:\n"
		for _, role := range guild.Admins {
			adminList += fmt.Sprintf("- %q (%s)\n", role.Name, role.ID)
		}
		if _, err := s.ChannelMessageSend(m.ChannelID, adminList); err != nil {
			log.Printf("Failed sending Admin Command response: %v", err)
		}
	case "add":
		if len(cmd.Args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: !admin add <role_id>"); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		roleID := cmd.Args[1]
		role, err := s.State.Role(gid, roleID)
		if err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role with ID %s not found.", roleID)); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		if guild.IsAdminRole(role) {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role %q is already an admin role.", role.Name)); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		guild.Admins = append(guild.Admins, role)
		if err := guild.Save(); err != nil {
			log.Printf("Failed saving guild configuration: %v", err)
			if _, err := s.ChannelMessageSend(m.ChannelID, "Failed to save admin role."); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return err
		}
		if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role %q added as an admin role.", role.Name)); err != nil {
			log.Printf("Failed sending Admin Command response: %v", err)
		}
		return nil
	case "remove":
		if len(cmd.Args) < 2 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "Usage: !admin remove <role_id>"); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		roleID := cmd.Args[1]
		role, err := s.State.Role(gid, roleID)
		if err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role with ID %s not found.", roleID)); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		if !guild.IsAdminRole(role) {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role %q is not an admin role.", role.Name)); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		for i, r := range guild.Admins {
			if r.ID == role.ID {
				guild.Admins = append(guild.Admins[:i], guild.Admins[i+1:]...)
				break
			}
		}
		if err := guild.Save(); err != nil {
			log.Printf("Failed saving guild configuration: %v", err)
			if _, err := s.ChannelMessageSend(m.ChannelID, "Failed to save admin role."); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return err
		}
		if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role %q removed from admin roles.", role.Name)); err != nil {
			log.Printf("Failed sending Admin Command response: %v", err)
		}
		return nil
	default:
		if _, err := s.ChannelMessageSend(m.ChannelID, "Unknown admin command. Use `!admin list`, `!admin add <role_id>`, or `!admin remove <role_id>`"); err != nil {
			log.Printf("Failed sending Admin Command response: %v", err)
		}
		return fmt.Errorf("unknown admin command: %s", cmd.Args[0])
	}

	// Admin command logic goes here

	return nil
}

func init() {
	RegisterCommand(NewCommand("Admin", "Admin commands", Admin))
}
