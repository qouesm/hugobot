package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	s               *discordgo.Session
	Token           string
	activeGuilds    []string
	appCommands     []discordgo.ApplicationCommand
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

// init vars
func init() {
	// using my own QBOT to test for now; new bot will come soon
	Token = os.Getenv("QBOT")
	activeGuilds = []string{
		"510285602785198081", // qserver
		"842613819057635328", // The Boys
		// "626546254846885948",  // NPCS
	}

	commandList := exportCommands()
	for _, c := range commandList {
		log.Println("adding to lists: ", c.AppCommand.Name)
		appCommands = append(appCommands, c.AppCommand)
		commandHandlers[c.AppCommand.Name] = c.Handler
	}
	log.Println("aC:", appCommands)

	// DEBUG
	// log.Println("appCommands:     ", appCommands)
	// log.Println("commandHandlers: ", commandHandlers)

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

	log.Println(appCommands)

	for _, g := range s.State.Guilds {
		for _, v := range appCommands {
			log.Println(g.Name)
			log.Println(v)
			// whitelist certain guilds for now
			if !isActiveGuild(g.ID) {
				continue
			}

			_, err := s.ApplicationCommandCreate(s.State.User.ID, g.ID, &v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
		}
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
