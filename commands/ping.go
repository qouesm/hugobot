package commands

import "github.com/bwmarrin/discordgo"

var Ping = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Ping!",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Pong!",
			},
		})
	},
}
