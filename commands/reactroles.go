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

					{
						Name:        "role4",
						Description: "4th role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionRole,
					},

					{
						Name:        "role5",
						Description: "5th role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionRole,
					},

					{
						Name:        "role6",
						Description: "6th role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionRole,
					},

					{
						Name:        "role7",
						Description: "7th role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionRole,
					},

					{
						Name:        "role8",
						Description: "8th role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionRole,
					},

					{
						Name:        "role9",
						Description: "9th role",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionRole,
					},

					{
						Name:        "role10",
						Description: "10th role",
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
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Type: discordgo.InteractionResponseChannelMessageWithSource,
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{},
			})

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
					msgFormat += numEmoji[line-1].Name // emoji from 0-9
					msgFormat += "<@&%s> "             // role mentionable
					msgFormat += "%s"                  // role name
					msgFormat += "\n"
				default:
					log.Println("unexpected type: ", t)
				}
			}

			rMsg, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf(
					msgFormat,
					mArgs...,
				),
			})
			if err != nil {
				log.Println("/reactroles create; problem creating message,", err)
				return
			}

			// s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{Content: "1a"})
			// s.InteractionResponseDelete(s.State.User.ID, i.Interaction)

			// TODO: next part; a proper formatted message has been created.
			// the message (rMsg) now needs to be populated with reactions
			// and the actual role add/remove function must be setup and stored.

			log.Println("adding reacts")
			err = s.MessageReactionAdd(rMsg.ChannelID, rMsg.ID, numEmoji[0].ID)
			if err != nil {
				log.Println("Could not add reaction", numEmoji[0].Name, ",", err)
			}
			// for num := 0; num < len(i.Data.Options[0].Options); num++ {
			// 	log.Println("react,", num)
			// 	s.MessageReactionAdd(rMsg.ChannelID, rMsg.ID, numEmoji[num].ID)
			// }

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

var numEmoji = map[int]discordgo.Emoji{
	0: {Name: ":zero:"},
	1: {Name: ":one:"},
	2: {Name: ":two:"},
	3: {Name: ":three:"},
	4: {Name: ":four:"},
	5: {Name: ":five:"},
	6: {Name: ":six:"},
	7: {Name: ":seven:"},
	8: {Name: "eight:"},
	9: {Name: ":nine:"},
}
