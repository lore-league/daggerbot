package config

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

const Version = "0.1.0"

var (
	Debug  bool
	Guilds map[string]*Guild
	Prefix string
	Token  string
)

func AddGuild(g *discordgo.Guild) error {
	if g == nil {
		return errors.New("attempted to add a nil guild")
	}
	if _, exists := Guilds[g.ID]; exists {
		log.Printf("Guild %q already exists, updating...", g.Name)
		Guilds[g.ID].SetRaw(g)
	} else {
		log.Printf("Connected to guild: %q (%s)", g.Name, g.ID)
		Guilds[g.ID] = NewGuild(g)
	}
	return nil
}

func init() {
	log.Println("Initializing configuration...")
	Debug = false // Default debug mode is off
	Guilds = make(map[string]*Guild)
	Prefix = "!" // Default command prefix
	Token = ""   // Token should be set externally
}
