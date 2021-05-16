package commands

import "github.com/bwmarrin/discordgo"

var Quietping = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "quietping",
		Description: "Shh... ping",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Flags:   1 << 6, // only visible to caller
				Content: "Pong!",
			},
		})
	},
}
