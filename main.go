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
		log.Println("found command:", c.AppCommand.Name)
		appCommands = append(appCommands, c.AppCommand)
		commandHandlers[c.AppCommand.Name] = c.Handler
	}

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
		log.Println("bot is online")
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("registering commands")
	for _, g := range s.State.Guilds {
		// whitelist certain guilds for now
		if !isActiveGuild(g.ID) {
			continue
		}
		for _, v := range appCommands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, g.ID, &v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
		}
	}
	log.Println("commands registered")

	defer s.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	log.Println("bot is ready")
	<-stop

	// log.Println("unregistering commands")
	// for _, g := range s.State.Guilds {
	// 	ac, err := s.ApplicationCommands(s.State.User.ID, g.ID)
	// 	if err != nil {
	// 		log.Printf("Problem getting application commands from %v, %v", g.Name, err)
	// 		continue
	// 	}
	// 	for _, v := range ac {
	// 		// log.Println("removing: ", v.Name)
	// 		err := s.ApplicationCommandDelete(s.State.User.ID, g.ID, v.ID)
	// 		if err != nil {
	// 			log.Printf("Cannot remove '%v' command: %v", v.Name, err)
	// 		}
	// 	}
	// }
	// log.Println("commands unregistered")

	log.Println("shutting down")
}

func isActiveGuild(ID string) bool {
	for _, activeID := range activeGuilds {
		if ID == activeID {
			return true
		}
	}
	return false
}
