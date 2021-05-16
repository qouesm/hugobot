package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var Options = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "options",
		Description: "Command for demonstrating options",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "string-option",
				Description: "String option",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "integer-option",
				Description: "Integer option",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "bool-option",
				Description: "Boolean option",
				Required:    true,
			},

			// Required options must be listed first since optional parameters
			// always come after when they're used.
			// The same concept applies to Discord's Slash-commands API

			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "channel-option",
				Description: "Channel option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user-option",
				Description: "User option",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role-option",
				Description: "Role option",
				Required:    false,
			},
		},
	},

	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		margs := []interface{}{
			// Here we need to convert raw interface{} value to wanted type.
			// Also, as you can see, here is used utility functions to convert the value
			// to particular type. Yeah, you can use just switch type,
			// but this is much simpler
			i.Data.Options[0].StringValue(),
			i.Data.Options[1].IntValue(),
			i.Data.Options[2].BoolValue(),
		}
		msgformat :=
			` Now you just learned how to use command options. Take a look to the value of which you've just entered:
			> string_option: %s
			> integer_option: %d
			> bool_option: %v
			`
		if len(i.Data.Options) >= 4 {
			margs = append(margs, i.Data.Options[3].ChannelValue(nil).ID)
			msgformat += "> channel-option: <#%s>\n"
		}
		if len(i.Data.Options) >= 5 {
			margs = append(margs, i.Data.Options[4].UserValue(nil).ID)
			msgformat += "> user-option: <@%s>\n"
		}
		if len(i.Data.Options) >= 6 {
			margs = append(margs, i.Data.Options[5].RoleValue(nil, "").ID)
			msgformat += "> role-option: <@&%s>\n"
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, we'll discuss them in "responses" part
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf(
					msgformat,
					margs...,
				),
			},
		})
	},
}
