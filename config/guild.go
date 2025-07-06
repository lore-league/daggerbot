package config

import (
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

/*
 * This package provides configuration specific to individual Discord guilds
 */

type Guild struct {
	raw    *discordgo.Guild  // Raw guild data
	ID     string            // Guild ID
	Name   string            // Guild name
	Admins []*discordgo.Role // List of admin role IDs
	Config map[string]string // Guild-specific configuration
}

type guildJSON struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Admins []string          `json:"admins"`
	Config map[string]string `json:"config"`
}

func NewGuild(guild *discordgo.Guild) *Guild {
	g := &Guild{
		raw:    guild,
		ID:     guild.ID,
		Name:   guild.Name,
		Admins: make([]*discordgo.Role, 0),
		Config: make(map[string]string),
	}
	g.Config["prefix"] = Prefix // Set default prefix
	if guild.Roles != nil {
		for _, role := range guild.Roles {
			if strings.EqualFold(role.Name, "admin") || strings.EqualFold(role.Name, "administrator") {
				g.Admins = append(g.Admins, role)
			}
		}
	}
	if Debug {
		log.Printf("[DEBUG] New Guild created: %q with ID %s", g.Name, g.ID)
		log.Printf("[DEBUG] found guild roles: %v", g.RoleNames())
	}

	// Load the guild configuration from file
	if err := g.Load(); err != nil {
		log.Printf("[ERROR] Failed to load guild configuration: %v", err)
	}

	return g
}

func (g *Guild) SetRaw(raw *discordgo.Guild) {
	g.raw = raw
	g.ID = raw.ID
	g.Name = raw.Name
}

func (g *Guild) Roles() iter.Seq[*discordgo.Role] {
	return func(yield func(*discordgo.Role) bool) {
		for _, r := range g.raw.Roles {
			if !yield(r) {
				return
			}
		}
	}
}

func (g *Guild) RoleNames() []string {
	rnames := make([]string, 0, len(g.raw.Roles))
	for r := range g.Roles() {
		rnames = append(rnames, r.Name)
	}
	return rnames
}

func (g *Guild) RoleIDs() []string {
	rids := make([]string, 0, len(g.raw.Roles))
	for r := range g.Roles() {
		rids = append(rids, r.ID)
	}
	return rids
}

func (g *Guild) AdminIDs() []string {
	aid := make([]string, 0, len(g.Admins))
	for _, r := range g.Admins {
		aid = append(aid, r.ID)
	}
	return aid
}

func (g *Guild) AddAdmin(rid string) error {
	if strings.TrimSpace(rid) == "" {
		return fmt.Errorf("attempted to add a nil role as admin")
	}

	role := g.FindRoleByID(rid)
	if role == nil {
		return fmt.Errorf("role with ID %s not found in guild %q", rid, g.Name)
	}
	if !slices.Contains(g.Admins, role) {
		g.Admins = append(g.Admins, role)
		if Verbose {
			log.Printf("[VERBOSE] Added role %q (%s) as admin in guild %q", role.Name, role.ID, g.Name)
		}
		return g.Save()
	}
	if Debug {
		log.Printf("[DEBUG] Role %q (%s) is already in the admin list for %q", role.Name, role.ID, g.Name)
	}
	return nil
}

func (g *Guild) FindRoleByName(name string) *discordgo.Role {
	for r := range g.Roles() {
		if strings.EqualFold(r.Name, strings.TrimSpace(name)) {
			return r
		}
	}
	return nil
}

func (g *Guild) FindRoleByID(id string) *discordgo.Role {
	for r := range g.Roles() {
		if r.ID == strings.TrimSpace(id) {
			return r
		}
	}
	return nil
}

func (g *Guild) IsOwner(user *discordgo.User) bool {
	if user == nil {
		return false
	}
	return g.raw.OwnerID == user.ID
}

func (g *Guild) IsAdmin(member *discordgo.Member) bool {
	if member == nil || member.User == nil {
		return false
	}

	user := member.User

	// Owners are always considered admins
	if g.IsOwner(user) {
		if Debug {
			log.Printf("[DEBUG] User %q is the owner of guild %q", user.Username, g.Name)
		}
		return true // Owner is always an admin
	}

	if len(member.Roles) == 0 {
		if Debug {
			log.Printf("[DEBUG] Member %q has no roles in guild %q\n", user.Username, g.Name)
		}
		return false // No user data available
	}

	for _, r := range member.Roles {
		role := g.FindRoleByID(r)
		if role == nil {
			log.Printf("[ERR] Role with ID %q not found in guild %q", r, g.Name)
			continue // Skip if role not found -- should not happen
		}
		if slices.Contains(g.AdminIDs(), r) {
			if Verbose {
				log.Printf("[VERBOSE] User %q has admin role %q in guild %q", user.Username, role.Name, g.Name)
			}
			return true // User has an admin role
		}
	}
	return false
}

func (g *Guild) IsAdminRole(role *discordgo.Role) bool {
	if role == nil {
		return false
	}

	for _, r := range g.Admins {
		if r.ID == role.ID {
			if Debug {
				log.Printf("[DEBUG] Role %q is an admin role in guild %q", role.Name, g.Name)
			}
			return true // Role is an admin role
		}
	}
	if Debug {
		log.Printf("[DEBUG] Role %q is not an admin role in guild %q", role.Name, g.Name)
	}
	return false // Role is not an admin role
}

func (g *Guild) Load() error {
	if Debug {
		log.Printf("[DEBUG] Loading configuration for guild %q (%s)", g.Name, g.ID)
	}

	file, err := os.Open(fmt.Sprintf("guild_%s.json", g.ID))
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Configuration file for guild %s does not exist, creating new one", g.ID)
			return g.Save() // Save a new configuration file if it doesn't exist
		}
		log.Printf("Failed to open guild configuration file: %v", err)
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	decoder := json.NewDecoder(file)
	var guildData guildJSON
	if err := decoder.Decode(&guildData); err != nil {
		log.Printf("Failed to decode guild configuration file: %v", err)
		return err
	}

	for _, adminID := range guildData.Admins {
		if adminID == "" {
			log.Printf("Admin role with ID %s not found in guild %q, skipping", adminID, g.Name)
			continue // Skip if role not found -- removed sometime after save
		}
		if err := g.AddAdmin(adminID); err != nil {
			log.Printf("Failed to add admin role %s to guild %q: %v", adminID, g.Name, err)
			return err // Return error if adding admin fails
		}
	}

	g.Config = guildData.Config

	log.Printf("Successfully loaded configuration for guild %q", g.Name)

	return nil
}

func (g *Guild) Save() error {
	if Debug {
		log.Printf("[DEBUG] Saving configuration for guild %q (%s)", g.Name, g.ID)
	}

	file, _ := os.OpenFile(fmt.Sprintf("guild_%s.json", g.ID), os.O_RDWR|os.O_CREATE, 0644)
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	guildData := guildJSON{
		ID:     g.ID,
		Name:   g.Name,
		Admins: g.AdminIDs(),
		Config: g.Config,
	}

	if err := encoder.Encode(guildData); err != nil {
		log.Printf("Failed to save guild configuration: %v", err)
		return err
	}

	log.Printf("Guild %q configuration saved successfully", g.Name)

	return nil
}

func (g *Guild) String() string {
	return fmt.Sprintf("Guild: %q, Roles: %v", g.Name, g.RoleNames())
}
