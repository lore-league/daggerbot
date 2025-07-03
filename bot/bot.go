package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nerdwerx/daggerbot/bot/commands"
)

var Config struct {
	Debug  bool
	Guilds map[string]Guild
	Prefix string
	Token  string
}

func startup() {
	// Load configuration from environment variables or other sources
	Config.Token = os.Getenv("DISCORD_AUTH_TOKEN")
	if Config.Token == "" {
		log.Fatal("DISCORD_AUTH_TOKEN is not set")
	}

	Config.Prefix = os.Getenv("DISCORD_BOT_PREFIX")
	if Config.Prefix == "" {
		Config.Prefix = "!" // Default prefix
	}
	log.Printf("Bot Prefix set to %q", Config.Prefix)

	Config.Guilds = make(map[string]Guild)

	log.Printf("Configuration loaded")
}

func Run() error {
	// Initialize our config vars
	startup()

	// create a session
	discord, err := discordgo.New("Bot " + Config.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	if Config.Debug {
		log.Println("Debug mode enabled")
		// discord.Debug = true
	}

	// add a event handlers
	discord.AddHandler(onReady)
	discord.AddHandler(onMessage)

	// Set our permissions
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	log.Println("Bot starting up... CTRL-C to stop")

	// open session
	if err := discord.Open(); err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}

	defer func() {
		if err := discord.Close(); err != nil { // close session, after function termination
			log.Fatal("Error closing Discord session")
		}
	}()

	// keep bot running untill we're interrupted (ctrl + C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Bot gracefully shutting down...")

	return nil
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	for _, guild := range r.Guilds {
		gid := guild.ID

		guildData, err := s.Guild(gid)
		if err != nil {
			log.Printf("Error fetching guild data for %s: %v", gid, err)
			continue
		}

		Config.Guilds[gid] = Guild{
			Name:  guildData.Name,
			Roles: guildData.Roles,
		}

		log.Printf("Connected to guild: %q (ID: %s)", Config.Guilds[gid].Name, gid)

		if Config.Debug {
			log.Printf("[DEBUG] Fetched %s", Config.Guilds[gid])
		}
	}

	log.Printf("Bot is ready! Connected to %d servers", len(Config.Guilds))
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		cmd     commands.Command
		command string
		debug   = Config.Debug
		fullcmd []string
		prefix  = Config.Prefix
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
		log.Printf("[DEBUG] Channel: %v; From: %s (%s); Message: %q", channel.Name, m.Author, m.Author.DisplayName(), m.Content)
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

	server, ok := Config.Guilds[m.GuildID]
	if !ok {
		log.Printf("Message received from invalid Guild: %s", m.GuildID)
		return
	}

	log.Printf("[%s] @%s executed command %q with args %v in channel %q", server.Name, m.Author, cmd.Name, cmd.Args, channel.Name)

	if err := cmd.Handler(s, m); err != nil {
		log.Printf("Error executing command %q: %v", cmd.Name, err)
	}

}
