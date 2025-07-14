package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/bot/commands"
	"github.com/nerdwerx/daggerbot/config"
)

func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		fullcmd = make([]string, 0)
		message = m.Content
		my      = s.State.User
	)

	// Ignore my own messages
	if m.Author.ID == my.ID {
		return
	}

	if m.Member.User == nil {
		m.Member.User = m.Author // Ensure m.Member.User is set to the message author
	}

	guild, ok := config.Guilds[m.GuildID]
	if !ok {
		log.Printf("Message received from invalid Guild: %s", m.GuildID)
		return
	}
	prefix := guild.Config["prefix"]

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Failed to fetch channel %s: %v", m.ChannelID, err)
		return
	}

	for _, mention := range m.Mentions {
		if mention.ID == my.ID {
			// Extract the command from the message
			if config.Debug {
				log.Printf("[DEBUG] Mention found: %v", mention)
			}
			if idx := strings.Index(message, prefix); idx != -1 {
				// If the prefix is found, split the message
				fullcmd = strings.Split(strings.TrimPrefix(message[idx:], prefix), " ")
			}
		}
	}

	if len(fullcmd) == 0 {
		if after, ok := strings.CutPrefix(message, prefix); ok {
			fullcmd = strings.Split(after, " ")
		} else {
			// If no mention or prefix found, return
			if config.Debug {
				log.Printf("[DEBUG] No command found in message: %q", message)
			}
			return
		}
	}

	command := strings.TrimSpace(fullcmd[0])
	if command == "" {
		return
	}

	cmd, ok := commands.Commands[strings.ToLower(command)]
	if !ok {
		if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sorry, I don't understand %q.", command)); err != nil {
			log.Printf("Failed sending Unknown Command response: %v", err)
		}
		if config.Verbose {
			log.Printf("[VERBOSE] Command %q not found in commands map", command)
		}
		return
	}

	cmd.SetGuild(guild) // Set the guild for the command

	// Inject any args into the command
	if len(fullcmd) > 1 {
		cmd.SetArgs(fullcmd[1:])
	}

	if cmd.Admin() && !guild.IsAdmin(m.Member) {
		log.Printf("[%s] user @%s (%s) is not an admin, denying access to %s command", guild.Name, m.Author.DisplayName(), m.Author, cmd.Name())
		if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("you must be an admin to use this command, %s", m.Author)); err != nil {
			log.Printf("Failed sending Admin Command response: %v", err)
		}
		return
	}

	log.Printf("[%s] @%s (%s) executing command %q with args %v in channel %q", guild.Name, m.Author.DisplayName(), m.Author, cmd.Name(), cmd.Args(), channel.Name)

	if err := cmd.Run(s, m); err != nil {
		log.Printf("Error executing command %q: %v", cmd.Name(), err)
	}

}
