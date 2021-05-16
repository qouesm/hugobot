package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	AppCommand discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
