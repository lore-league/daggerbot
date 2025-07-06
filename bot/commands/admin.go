package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Admin(cmd *Command, s *discordgo.Session, m *discordgo.MessageCreate) error {
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
		if len(cmd.Guild.Admins) == 0 {
			if _, err := s.ChannelMessageSend(m.ChannelID, "No admin roles configured for this guild."); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		adminList := "Admin roles for this guild:\n"
		for _, role := range cmd.Guild.Admins {
			adminList += fmt.Sprintf("- %s [%s]\n", role.Name, role.ID)
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
		role, err := s.State.Role(cmd.Guild.ID, roleID)
		if err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role with ID %s not found.", roleID)); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		if cmd.Guild.IsAdminRole(role) {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role %q is already an admin role.", role.Name)); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		cmd.Guild.Admins = append(cmd.Guild.Admins, role)
		if err := cmd.Guild.Save(); err != nil {
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
		role, err := s.State.Role(cmd.Guild.ID, roleID)
		if err != nil {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role with ID %s not found.", roleID)); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		if !cmd.Guild.IsAdminRole(role) {
			if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role %q is not an admin role.", role.Name)); err != nil {
				log.Printf("Failed sending Admin Command response: %v", err)
			}
			return nil
		}
		for i, r := range cmd.Guild.Admins {
			if r.ID == role.ID {
				cmd.Guild.Admins = append(cmd.Guild.Admins[:i], cmd.Guild.Admins[i+1:]...)
				break
			}
		}
		if err := cmd.Guild.Save(); err != nil {
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
	cmd := NewCommand("Admin", "Admin commands", Admin)
	cmd.Admin = true
	RegisterCommand(cmd)
}
