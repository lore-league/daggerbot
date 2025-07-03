package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Guild struct {
	Name       string            // Name of the guild
	Roles      []*discordgo.Role // Roles for each guild
	AdminRoles []*discordgo.Role // Admin roles for the guild
}

func (g Guild) String() string {
	roleNames := make([]string, len(g.Roles))
	for i, role := range g.Roles {
		roleNames[i] = fmt.Sprintf("%q", role.Name)
	}
	return fmt.Sprintf("Guild: %q, Roles: %v", g.Name, roleNames)
}
