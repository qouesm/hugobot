package commands

import "github.com/bwmarrin/discordgo"

type Command struct {
	appCommand discordgo.ApplicationCommand
	handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
