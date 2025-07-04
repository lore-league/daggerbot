package handlers

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/config"
)

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	for _, g := range r.Guilds {
		gid := g.ID

		gdata, err := s.Guild(gid)
		if err != nil {
			log.Printf("Error fetching guild data for %s: %v", gid, err)
			continue
		}

		if err := config.AddGuild(gdata); err != nil {
			log.Printf("Error adding guild %s: %v", gid, err)
		}

		if config.Debug {
			log.Printf("[DEBUG] Fetched Guild %q with Roles %s", gdata.Name, config.Guilds[gid].RoleNames())
		}
	}

	log.Printf("Bot is ready! Connected to %d servers", len(config.Guilds))
}
