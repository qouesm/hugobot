package commands

import "github.com/bwmarrin/discordgo"

var ClassClear = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "classclear",
		Description: "ADMIN: Remove all class roles from all server members",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: i.Member.Nick,
			},
		})
	},
}
