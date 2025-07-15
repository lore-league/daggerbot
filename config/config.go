package config

import (
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

const Version = "0.1.1"

const IDRegex = `^\d{17,19}$` // Matches Discord IDs (17-19 digits)

var (
	Debug   bool
	Guilds  map[string]*Guild
	Prefix  string
	Token   string
	Verbose bool
)

func RegisterGuild(g *discordgo.Guild) error {
	if g == nil {
		return errors.New("attempted to add a nil guild")
	}
	if _, exists := Guilds[g.ID]; exists {
		log.Printf("updating %q", g.Name)
		Guilds[g.ID].Update(g)
	} else {
		log.Printf("registered guild: %q (%s)", g.Name, g.ID)
		Guilds[g.ID] = NewGuild(g)
	}
	return nil
}

func SaveGuilds() error {
	for _, g := range Guilds {
		if err := g.Save(); err != nil {
			log.Printf("error saving guild %q: %v", g.Name, err)
			return err
		}
	}
	log.Println("all guild configurations saved successfully")
	return nil
}

func init() {
	log.Println("initializing configuration...")
	Guilds = make(map[string]*Guild)
	Prefix = "!" // Default command prefix
	Token = ""   // Token should be set externally
}
