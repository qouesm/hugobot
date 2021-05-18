package commands

import (
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var ReactRoles = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "reactroles",
		Description: "Create a reaction roles message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "create",
				Description: "Create a reaction role message",
				// Type:			discordgo.ApplicationCommandOptionSubCommandGroup,
				Type: discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "header",
						Description: "Header of message",
						Required:    true,
						Type:        discordgo.ApplicationCommandOptionString,
					},

					{
						Name:        "role1",
						Description: "1st role",
						Required:    true,
						Type:        discordgo.ApplicationCommandOptionRole,
					},

					{
						Name:        "role2",
						Description: "2nd role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionRole,
					},

					{
						Name:        "role3",
						Description: "3rd role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionRole,
					},
				},
			},

			{
				Name:        "edit",
				Description: "Edit a reaction role message",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var (
			msgFormat string
			cArgs     []interface{} // used internally
			mArgs     []interface{} // used for msgFormat
		)

		switch i.Data.Options[0].Name {
		case "create":
			cArgs = []interface{}{
				i.Data.Options[0].Options[0].StringValue(), // header
			}
			mArgs = []interface{}{
				i.Data.Options[0].Options[0].StringValue(),
			}

			// add roles
			// roles add out of order b/c race condition; fixed with sync.WaitGroup
			var wgIns sync.WaitGroup
			wgIns.Add(1)
			for count := 1; count < len(i.Data.Options[0].Options); count++ {
				role := i.Data.Options[0].Options[count].RoleValue(s, i.GuildID)
				cArgs = append(cArgs, role)
				mArgs = append(mArgs, role.ID, role.Name)
			}
			wgIns.Done()

			for line, v := range cArgs {
				switch t := v.(type) {
				case string:
					msgFormat += "```\n" + `%s` + "\n```\n" // first line is header inside ``
				case *discordgo.Role:
					msgFormat += numEmoji[line-1] // emoji from 0-9
					msgFormat += "<@&%s> "        // role @
					msgFormat += "%s"             // role name
					msgFormat += "\n"
				default:
					log.Println("unexpected type: ", t)
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: fmt.Sprintf(
						msgFormat,
						mArgs...,
					),
				},
			})

			// TODO: next part; a proper formatted message has been created.
			// the message now needs to be populated with reactions
			// and the actual role add/remove function must be setup and stored.

		case "edit":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Content: "Not implimented yet",
				},
			})

		default:
			log.Println("how tf")
		}
	},
}

var numEmoji = map[int]string{
	0: ":zero:",
	1: ":one:",
	2: ":two:",
	3: ":three:",
	4: ":four:",
	5: ":five:",
	6: ":six:",
	7: ":seven:",
	8: ":eight:",
	9: ":nine:",
}
