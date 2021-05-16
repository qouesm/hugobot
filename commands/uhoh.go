package commands

import "github.com/bwmarrin/discordgo"

var Uhoh = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "uhoh",
		Description: "Uh oh!",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Stinky!",
			},
		})
	},
}
