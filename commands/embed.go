package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var Embed = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "embed",
		Description: "create embed message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "title",
				Description: "Title",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
			},

			{
				Name:        "description",
				Description: "Description",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
			},

			{
				Name:        "url",
				Description: "URL",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionString,
			},

			{
				Name:        "color",
				Description: "Color format: `0xRRGGBB",
				Required:    false,
				Type:        discordgo.ApplicationCommandOptionInteger,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		if len(i.Data.Options) == 0 {
			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{Content: "I can't create an embed message with no content"})
			return
		}

		var embed discordgo.MessageEmbed

		for _, v := range i.Data.Options {
			switch v.Name {
			case "title":
				embed.Title = v.StringValue()
			case "description":
				embed.Description = v.StringValue()
			case "url":
				embed.Description = v.StringValue()
			case "color":
				embed.Color = int(v.IntValue())
			default:
				log.Println("bad type")
			}
		}

		_, err := s.ChannelMessageSendEmbed(i.ChannelID, &embed)
		if err != nil {
			log.Println("problem send embed message,", err)
		}

		err = s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
		if err != nil {
			log.Println("problem deleting defer message,", err)
		}
	},
}
