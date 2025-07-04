package config

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Guild struct {
	raw        *discordgo.Guild  // Raw guild data
	Name       string            // Guild name
	Roles      []*discordgo.Role // Roles for the guild
	AdminRoles []*discordgo.Role // Admin roles for the guild
}

func NewGuild(guild *discordgo.Guild) *Guild {
	return &Guild{
		raw:        guild,
		Name:       guild.Name,
		Roles:      guild.Roles,
		AdminRoles: make([]*discordgo.Role, 0),
	}
}

func (g *Guild) SetRaw(raw *discordgo.Guild) {
	g.raw = raw
	g.Name = raw.Name
	g.Roles = raw.Roles
}

func (g *Guild) RoleNames() []string {
	rnames := make([]string, len(g.Roles))
	for i, r := range g.Roles {
		rnames[i] = fmt.Sprintf("%q", r.Name)
	}
	return rnames
}

func (g *Guild) String() string {
	return fmt.Sprintf("Guild: %q, Roles: %v", g.Name, g.RoleNames())
}
