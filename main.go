package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/qouesm/hugobot/hooks"
)

var (
	s               *discordgo.Session
	Token           string
	appCommands     []discordgo.ApplicationCommand
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

// init vars
func init() {
	go herokuListen()
	Token = os.Getenv("HUGOBOT")
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
	s.UpdateGameStatus(4, "⚙ bot is starting...")

	defer s.Close()

	stop := make(chan os.Signal)

	log.Println("restarting hooks")
	messages, err := ioutil.ReadDir("hooks/messages")
	if err != nil {
		log.Panicln("COULD NOT START REACT HOOKS,", err)
	}
	for _, msg := range messages {
		hooks.ReactRoles(s, msg.Name())
	}
	log.Println("started hooks")

	log.Println("registering commands")
	for _, g := range s.State.Guilds {
		if g.Name == "new paltz cs" {  // this shit won't go away
			err := s.State.GuildRemove(g)
			if err != nil {
				log.Println("fuck")
			}
		}
		for _, v := range appCommands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, g.ID, &v)
			if err != nil {
				log.Printf("Cannot create '%v' command in guild '%v': %v", v.Name, g.Name, err)
			}
		}
	}
	log.Println("commands registered")

	s.UpdateGameStatus(4, "cooler than @Hugo the Hawk")

	signal.Notify(stop, os.Interrupt)
	log.Println("bot is ready")
	<-stop

	log.Println("unregistering commands")
	for _, g := range s.State.Guilds {
		ac, err := s.ApplicationCommands(s.State.User.ID, g.ID)
		if err != nil {
			log.Printf("Problem getting application commands from %v, %v", g.Name, err)
			continue
		}
		for _, v := range ac {
			err := s.ApplicationCommandDelete(s.State.User.ID, g.ID, v.ID)
			if err != nil {
				log.Printf("Cannot remove '%v' command: %v", v.Name, err)
			}
		}
	}
	log.Println("commands unregistered")

	log.Println("shutting down")
}

func herokuListen() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT not set")
		return
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
