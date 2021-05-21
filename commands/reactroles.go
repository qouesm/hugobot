package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/qouesm/hugobot/hooks"
)

var ReactRoles = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "reactroles",
		Description: "Create a reaction roles message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "create",
				Description: "Create a reaction role message",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
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
			mRoles    []*discordgo.Role
		)

		switch i.Data.Options[0].Name {
		case "create":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{},
			})

			/* what this case does:
			parses the options,
			creates a message,
			adds reactions,
			saves message struct and roles array to json,
			call reactroles hook (handlers are there)
			*/

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
				mRoles = append(mRoles, role)
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

			// create message
			rMsg, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf(
					msgFormat,
					mArgs...,
				),
			})
			if err != nil {
				log.Println("/reactroles create; problem creating message,", err)
				panicResponse(s, i)
				return
			}

			// add reacitons
			for num := 0; num < len(i.Data.Options[0].Options)-1; num++ {
				s.MessageReactionAdd(rMsg.ChannelID, rMsg.ID, numEmoji[num].APIName())
			}

			// save to json
			var save = JsonSave{
				Msg:   rMsg,
				Roles: mRoles,
			}

			// file, err := os.Create("template.json")
			file, err := os.OpenFile("hooks/reactrolesmessages.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Println("problem reading file,", err)
				panicResponse(s, i)
				err := s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
				if err != nil {
					log.Println("Could not delete response (you got some big problems,", err)
				}
				return
			}
			defer file.Close()

			js, err := json.Marshal(save)
			if err != nil {
				log.Println("problem using marshal,", err)
				panicResponse(s, i)
				err := s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
				if err != nil {
					log.Println("Could not delete response (you got some big problems,", err)
				}
				return
			}
			file.Write(js)
			file.WriteString("\n")

			hooks.ReactRoles(s)

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

// returns discordgo Emoji struct from the corresponding int
var numEmoji = map[int]*discordgo.Emoji{
	0: {Name: "0️⃣"},
	1: {Name: "1️⃣"},
	2: {Name: "2️⃣"},
	3: {Name: "3️⃣"},
	4: {Name: "4️⃣"},
	5: {Name: "5️⃣"},
	6: {Name: "6️⃣"},
	7: {Name: "7️⃣"},
	8: {Name: "8️⃣"},
	9: {Name: "9️⃣"},
}

type JsonSave struct {
	Msg   *discordgo.Message
	Roles []*discordgo.Role
}
