package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	Token        string
	activeGuilds []string

	appCommands     []*discordgo.ApplicationCommand
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	s *discordgo.Session
)

// init vars
func init() {
	Token = os.Getenv("QBOT")
	activeGuilds = []string{
		"510285602785198081", // qserver
		"842613819057635328", // The Boys
		// "626546254846885948",  // NPCS
	}

	commandList := exportCommands()
	for _, c := range commandList {
		log.Println("adding to lists: ", c.AppCommand.Name)
		appCommands = append(appCommands, &c.AppCommand)
		commandHandlers[c.AppCommand.Name] = c.Handler
	}

	// DEBUG
	log.Println("appCommands:     ", appCommands)
	log.Println("commandHandlers: ", commandHandlers)

	var err error
	s, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("Bad token,", err)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.Data.Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is up!")
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	for i, v := range appCommands {
		for j, g := range s.State.Guilds {
			// whitelist certain guilds for now
			if !isActiveGuild(g.ID) {
				continue
			}

			c, err := s.ApplicationCommandCreate(s.State.User.ID, g.ID, v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}

			if i == 0 {
				log.Println("Added command: ", c)
			}
			if j == 0 {
				log.Println("bot is active in: ", g.Name)
			}
		}
	}

	log.Println("All builds bot is in")
	for _, g := range s.State.Guilds {
		log.Println(g.Name)
	}

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")
}

func isActiveGuild(ID string) bool {
	for _, activeID := range activeGuilds {
		if ID == activeID {
			return true
		}
	}
	return false
}
