package commands

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var RestartHandlers = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "restarthandlers",
		Description: "temporary command to restart reactrole message handlers",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{},
		})

		js, err := os.ReadFile("template.json")
		if err != nil {
			log.Println("problem reading file", err)
		}
		var rMsg *discordgo.Message
		// rMsg, err = s.State.Message()
		err = json.Unmarshal(js, &rMsg)
		if err != nil {
			log.Println("problem using unmarshal")
		}

		log.Println(rMsg.ID)
		log.Println(rMsg.Content)
	},
}
