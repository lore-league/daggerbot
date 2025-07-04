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
		cmd     commands.Command
		command string
		debug   = config.Debug
		fullcmd []string
		prefix  = config.Prefix
		my      = s.State.User
		ok      bool
	)

	// Ignore my own messages
	if m.Author.ID == my.ID {
		return
	}

	message := m.Content
	channel, _ := s.Channel(m.ChannelID)

	if debug {
		log.Printf("[DEBUG] Channel: %v; Author: %s; Message: %q", channel.Name, m.Author, m.Content)
	}

	for _, mention := range m.Mentions {
		if mention.ID == my.ID {
			// Extract the command from the message
			if debug {
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
			if debug {
				log.Printf("[DEBUG] No command found in message: %q", message)
			}
			return
		}
	}

	command = strings.TrimSpace(fullcmd[0])
	if command == "" {
		return
	}

	if cmd, ok = commands.Commands[command]; !ok {
		if _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sorry, I don't understand %q.", command)); err != nil {
			log.Printf("Failed sending Unknown Command response: %v", err)
		}
		if debug {
			log.Printf("[DEBUG] Command %q not found in commands map", command)
		}
		return
	}

	// Inject any args into the command
	if len(fullcmd) > 1 {
		cmd.Args = fullcmd[1:]
	} else {
		cmd.Args = []string{}
	}

	guild, ok := config.Guilds[m.GuildID]
	if !ok {
		log.Printf("Message received from invalid Guild: %s", m.GuildID)
		return
	}

	log.Printf("[%s] @%s (%s) executed command %q with args %v in channel %q", guild.Name, m.Author.DisplayName(), m.Author, cmd.Name, cmd.Args, channel.Name)

	if err := cmd.Handler(s, m); err != nil {
		log.Printf("Error executing command %q: %v", cmd.Name, err)
	}

}
