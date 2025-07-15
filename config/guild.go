package config

import (
	"encoding/json"
	"fmt"
	"iter"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

/*
 * This package provides configuration specific to individual Discord guilds
 */

type Config struct {
	prefix  string            // Default command prefix for the bot
	admins  []*discordgo.Role // List of admin role IDs
	gms     []*discordgo.Role // List of admin role IDs
	players []*discordgo.Role // List of admin role IDs
}

type Guild struct {
	ID     string           // Guild ID
	Name   string           // Guild name
	config *Config          // Guild-specific configuration
	guild  *discordgo.Guild // guild data
}

var ConfigKeys = []string{"prefix", "admins", "gms", "players"} // Keys used in the configuration
type ConfigMap map[string][]string                              // JSON representation of the configuration

type guildJSON struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Config ConfigMap `json:"config"`
}

func NewConfig() *Config {
	return &Config{
		prefix:  Prefix, // Default command prefix
		admins:  make([]*discordgo.Role, 0),
		gms:     make([]*discordgo.Role, 0),
		players: make([]*discordgo.Role, 0),
	}
}

func NewGuild(guild *discordgo.Guild) *Guild {
	g := &Guild{
		ID:     guild.ID,
		Name:   guild.Name,
		config: NewConfig(),
		guild:  guild,
	}

	if guild.Roles != nil {
		for _, role := range guild.Roles {
			if strings.EqualFold(role.Name, "admin") || strings.EqualFold(role.Name, "administrator") {
				g.config.admins = append(g.config.admins, role)
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

func (g *Guild) Update(guild *discordgo.Guild) {
	g.guild = guild
	g.ID = guild.ID
	g.Name = guild.Name
}

func (g *Guild) Roles() iter.Seq[*discordgo.Role] {
	return func(yield func(*discordgo.Role) bool) {
		for _, r := range g.guild.Roles {
			if !yield(r) {
				return
			}
		}
	}
}

func (g *Guild) RoleNames() []string {
	rnames := make([]string, 0, len(g.guild.Roles))
	for r := range g.Roles() {
		rnames = append(rnames, r.Name)
	}
	return rnames
}

func (g *Guild) RoleIDs() []string {
	rids := make([]string, 0, len(g.guild.Roles))
	for r := range g.Roles() {
		rids = append(rids, r.ID)
	}
	return rids
}

func (g *Guild) AdminIDs() []string {
	aid := make([]string, 0, len(g.config.admins))
	for _, r := range g.config.admins {
		aid = append(aid, r.ID)
	}
	return aid
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
		if r.ID == cleanString(id) {
			return r
		}
	}
	return nil
}

// FindRole searches for a role by either ID or name.
func (g *Guild) FindRole(tag string) *discordgo.Role {
	if regexp.MustCompile(IDRegex).MatchString(tag) {
		return g.FindRoleByID(tag)
	}
	return g.FindRoleByName(tag)
}

func (g *Guild) RolesToNames(roles []*discordgo.Role) []string {
	if len(roles) == 0 {
		return nil
	}
	names := make([]string, 0, len(roles))
	for _, r := range roles {
		if r != nil {
			names = append(names, r.Name)
		}
	}
	return names
}

func (g *Guild) IsOwner(user *discordgo.User) bool {
	if user == nil {
		return false
	}
	return g.guild.OwnerID == user.ID
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

	for _, r := range g.config.admins {
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
	var gjdata guildJSON

	if err := decoder.Decode(&gjdata); err != nil {
		log.Printf("Failed to decode guild configuration file: %v", err)
		return err
	}

	g.ID = gjdata.ID
	g.Name = gjdata.Name
	g.SetConfigMap(gjdata.Config)

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
		Config: g.GetConfigMap(),
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

func (g *Guild) Config() *Config {
	if g.config == nil {
		log.Printf("[WARN] Config is nil for guild %q, initializing with default values", g.Name)
		g.config = NewConfig() // Ensure config is initialized
	}
	return g.config
}

func (g *Guild) GetRoleConfig(key string) ([]*discordgo.Role, error) {
	switch cleanString(key) {
	case "admins", "admin":
		return g.config.admins, nil
	case "gms", "gm":
		return g.config.gms, nil
	case "players", "player":
		return g.config.players, nil
	default:
		return nil, fmt.Errorf("unrecognized role config key %q for guild %q", key, g.Name)
	}
}

func (g *Guild) SetRoleConfig(key string, values []*discordgo.Role) error {
	if key == "" || len(values) < 1 {
		return fmt.Errorf("neither config key nor value can be empty")
	}
	if Debug {
		log.Printf("[DEBUG] Setting config key %q for guild %q to %v", key, g.Name, values)
	}
	switch cleanString(key) {
	case "admins", "admin":
		g.config.admins = values
	case "gms", "gm":
		g.config.gms = values
	case "players", "player":
		g.config.players = values
	default:
		return fmt.Errorf("unrecognized config key %s for guild %q", key, g.Name)
	}
	return g.Save()
}

func (g *Guild) Prefix() string {
	return g.config.prefix
}

func (g *Guild) SetPrefix(prefix string) {
	if strings.TrimSpace(prefix) == "" {
		log.Printf("[WARN] Attempted to set an empty prefix for guild %q, using default: !", g.Name)
		g.config.prefix = Prefix // Reset to default if empty
	} else {
		g.config.prefix = strings.TrimSpace(prefix)
	}
	if Verbose {
		log.Printf("[VERBOSE] Set prefix for guild %q to: %s", g.Name, g.config.prefix)
	}
	if err := g.Save(); err != nil {
		log.Printf("[ERROR] Failed to save updated prefix for guild %q: %v", g.Name, err)
	}
}

func (g *Guild) GetConfigMap() ConfigMap {
	config := make(ConfigMap)

	config["prefix"] = []string{g.config.prefix}

	config["admins"] = make([]string, 0, len(g.config.admins))
	if len(g.config.admins) > 0 {
		adminIDs := make([]string, 0, len(g.config.admins))
		for _, role := range g.config.admins {
			adminIDs = append(adminIDs, role.ID)
		}
		config["admins"] = adminIDs
	}

	config["gms"] = make([]string, 0, len(g.config.gms))
	if len(g.config.gms) > 0 {
		gmIDs := make([]string, 0, len(g.config.gms))
		for _, role := range g.config.gms {
			gmIDs = append(gmIDs, role.ID)
		}
		config["gms"] = gmIDs
	}

	config["players"] = make([]string, 0, len(g.config.players))
	if len(g.config.players) > 0 {
		playerIDs := make([]string, 0, len(g.config.players))
		for _, role := range g.config.players {
			playerIDs = append(playerIDs, role.ID)
		}
		config["players"] = playerIDs
	}

	return config
}

func (g *Guild) SetConfigMap(config ConfigMap) {
	if Debug {
		log.Printf("[DEBUG] ConfigFromString called for guild %q with config: %v", g.Name, config)
	}

	// Initialize the config with default values
	g.config = NewConfig()

	for key, values := range config {
		if len(values) == 0 {
			log.Printf("no roles found in config for %s in guild %q, skipping", key, g.Name)
			continue // Skip if no admin roles are defined
		}

		switch cleanString(key) {
		case "prefix":
			if len(values) > 0 {
				g.config.prefix = values[0] // Set the prefix from the config
				if Debug {
					log.Printf("[DEBUG] Loaded prefix for guild %q: %s", g.Name, g.config.prefix)
				}
			} else {
				log.Printf("no prefix found in config for guild %q, using default: %s", g.Name, g.config.prefix)
			}
		case "admins", "admin":
			if Debug {
				log.Printf("[DEBUG] Found %d admin roles in config for guild %q", len(values), g.Name)
			}
			g.config.admins = g.getRoles(values)
			if len(g.config.admins) == 0 {
				log.Printf("no valid admin roles found in config for guild %q, using default admin role", g.Name)
				defaultAdminRole := g.FindRoleByName("admin")
				if defaultAdminRole == nil {
					defaultAdminRole = g.FindRoleByName("administrator")
				}
				if defaultAdminRole != nil {
					g.config.admins = append(g.config.admins, defaultAdminRole)
					log.Printf("added default admin role %q (%s) to guild %q", defaultAdminRole.Name, defaultAdminRole.ID, g.Name)
				} else {
					log.Printf("[WARN] No default admin role found in guild %q, admins list will be empty", g.Name)
				}
			}
		case "gms", "gm":
			if Debug {
				log.Printf("[DEBUG] Found %d GM roles in config for guild %q", len(values), g.Name)
			}
			g.config.gms = g.getRoles(values)
			if len(g.config.gms) == 0 {
				log.Printf("no valid GM roles found in config for guild %q, using default GM role", g.Name)
				defaultGMRole := g.FindRoleByName("gm")
				if defaultGMRole != nil {
					g.config.gms = append(g.config.gms, defaultGMRole)
					log.Printf("added default GM role %q (%s) to guild %q", defaultGMRole.Name, defaultGMRole.ID, g.Name)
				} else {
					log.Printf("[WARN] No default GM role found in guild %q, GMs list will be empty", g.Name)
				}
			}
		case "players", "player":
			if Debug {
				log.Printf("[DEBUG] Found %d Player roles in config for guild %q", len(values), g.Name)
			}
			g.config.players = g.getRoles(values)
			if len(g.config.players) == 0 {
				log.Printf("No valid Player roles found in config for guild %q, using default Player role", g.Name)
				defaultPlayerRole := g.FindRoleByName("player")
				if defaultPlayerRole != nil {
					g.config.players = append(g.config.players, defaultPlayerRole)
					log.Printf("added default Player role %q (%s) to guild %q", defaultPlayerRole.Name, defaultPlayerRole.ID, g.Name)
				} else {
					log.Printf("[WARN] No default Player role found in guild %q, Players list will be empty", g.Name)
				}
			}
		default:
			if Debug {
				log.Printf("[DEBUG] Unrecognized config key %q in guild %q, skipping", key, g.Name)
			}
		}
	}
}

func (g *Guild) ClearConfig(key string) error {
	if Debug {
		log.Printf("[DEBUG] Clearing config key %q for guild %q", key, g.Name)
	}

	if key == "prefix" {
		log.Printf("[WARN] Attempted to clear the prefix key for guild %q, resetting to default: !", g.Name)
		g.config.prefix = Prefix // Reset to default if prefix is cleared
		return g.Save()
	}

	switch key {
	case "admins", "admin":
		if Debug {
			log.Printf("[DEBUG] Clearing admin roles for guild %q", g.Name)
		}
		g.config.admins = make([]*discordgo.Role, 0) // Clear admin roles
	case "gms", "gm":
		if Debug {
			log.Printf("[DEBUG] Clearing GM roles for guild %q", g.Name)
		}
		g.config.gms = make([]*discordgo.Role, 0) // Clear GM roles
	case "players", "player":
		if Debug {
			log.Printf("[DEBUG] Clearing Player roles for guild %q", g.Name)
		}
		g.config.players = make([]*discordgo.Role, 0) // Clear Player roles
	default:
		log.Printf("[WARN] Attempted to clear unrecognized config key %q for guild %q", key, g.Name)
		return fmt.Errorf("unrecognized config key %s for guild %q", key, g.Name)
	}
	return g.Save()
}

/*
 * Private methods for Guild configuration management
 */

func (g *Guild) getRoles(values []string) []*discordgo.Role {
	roles := make([]*discordgo.Role, 0, len(values))

	if len(values) == 0 {
		if Debug {
			log.Printf("[DEBUG] No role IDs provided to GetRoles() for guild %q, returning empty list", g.Name)
		}
		return roles // Return empty list if no role IDs provided
	}

	if Debug {
		log.Printf("[DEBUG] Getting roles for guild %q with IDs: %v", g.Name, values)
	}

	// Iterate through the provided role IDs and find the corresponding roles
	for _, id := range values {
		if id == "" {
			log.Printf("found blank role ID during load for guild %q, skipping", g.Name)
			continue // Skip if role not found -- removed sometime after save
		}
		role := g.FindRoleByID(id)
		if role == nil {
			log.Printf("role with ID %s not found in guild %q during load", id, g.Name)
			continue // Skip if role not found -- removed sometime after save
		}
		roles = append(roles, role)
	}
	return roles
}

func cleanString(s string) string {
	if s == "" {
		return s
	}
	s = strings.ToLower(strings.TrimSpace(s))
	if len(s) > 2000 {
		log.Printf("cleanString: input exceeds Discord's 2000 character limit")
		return s[:2000] // Truncate to 2000 characters
	}
	return s
}
