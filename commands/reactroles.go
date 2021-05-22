package commands

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/qouesm/hugobot/hooks"
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

					{
						Name:        "color",
						Description: "Color in format: `0xRRGGBB`",
						Required:    false,
						Type:        discordgo.ApplicationCommandOptionInteger,
					},
				},
			},

			// {
			// 	Name:        "edit",
			// 	Description: "Edit a reaction role message",
			// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
			//  Options: []*discordgo.ApplicationCommandOption{}
			// },

			// {
			// 	Name:        "delete",
			// 	Description: "Delete a reaction role message (prefered to manual deletion)",
			// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
			//  Options: []*discordgo.ApplicationCommandOption{}
			// },

		},
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{},
		})

		var roles []*discordgo.Role

		if !hasAdmin(s.State, i.Member.Roles, i.GuildID) {
			m, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{Content: "You are not allowed to use that command"})
			if err != nil {
				log.Println("problem creating followup message,", err)
			}
			time.Sleep(time.Second * 5)
			s.FollowupMessageDelete(s.State.User.ID, i.Interaction, m.ID)
			return
		}

		switch i.Data.Options[0].Name {
		case "create":
			/* what this case does:
			parses the options,
			builds a message,
			adds reactions,
			saves message struct and roles array to json,
			call reactroles hook (handlers are there)
			*/

			// create message
			var embed = discordgo.MessageEmbed{
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "\u200b", // zero width space
						Inline: true,
					},
					{
						Name:   "Role",
						Inline: true,
					},
				},
			}

			reClass, _ := regexp.Compile("(CPS|MAT|EGC)[0-9]{3}")

			// if any of the provided roles are courses we want a course field
			for _, v := range i.Data.Options[0].Options {
				if v.Name[:len(v.Name)-1] == "role" {
					role := v.RoleValue(s, i.GuildID)
					if reClass.MatchString(role.Name) {
						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:   "Course",
							Inline: true,
						})
						break
					}
				}
			}

			// roles add out of order b/c race condition; fixed with sync.WaitGroup
			var wgIns sync.WaitGroup
			wgIns.Add(1)
			lineNum := 0
			for _, v := range i.Data.Options[0].Options {
				switch v.Name {
				case "header":
					embed.Title = v.StringValue()
				case "color":
					embed.Color = int(v.IntValue())
				default: // role
					embed.Fields[0].Value += numEmoji[lineNum].APIName() + "\n" // add emoji to msg
					lineNum++                                                   // i cant do this inline?
					role := v.RoleValue(s, i.GuildID)                           // get role
					roles = append(roles, role)                                 // add role to role list
					embed.Fields[1].Value += role.Mention() + "\n"              // add role.Name to msg
					if len(embed.Fields) > 2 {
						if reClass.MatchString(role.Name) {
							embed.Fields[2].Value += courseName[role.Name]
						}
						embed.Fields[2].Value += "\n"
					}
				}
			}
			wgIns.Done()

			// i dont fully understand the api but i know that this will delete the "thinking" message even if it's sloppy
			m, _ := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{Content: "\u200b"})
			s.FollowupMessageDelete(s.State.User.ID, i.Interaction, m.ID)

			msg, err := s.ChannelMessageSendEmbed(i.ChannelID, &embed)
			if err != nil {
				log.Println("/reactroles create; problem creating message,", err)
				panicResponse(s, i)
				return
			}

			// add reacitons
			for num := 0; num < len(roles); num++ {
				s.MessageReactionAdd(msg.ChannelID, msg.ID, numEmoji[num].APIName())
			}

			// save to json
			var save = JsonSave{
				Msg:   msg,
				Roles: roles,
			}

			file, err := os.Create("hooks/messages/" + msg.ID + ".json")
			if err != nil {
				log.Println("problem creating file,", err)
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

			// start hook
			hooks.ReactRoles(s, msg.ID+".json")

		// case "edit":
		// 	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{Content: "not implimented yet"})

		// case "delete":
		// 	s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{Content: "not implimented yet"})

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

var courseName = map[string]string{
	"CPS210": "CS1: Foundations",
	"CPS310": "CS2: Data Structures",
	"CPS315": "CS3: ~~Advanced Data Structures~~",
	"CPS330": "Assembly Language",
	"CPS340": "Operating Systems",
	"CPS352": "Object Oriented Programming",
	"CPS353": "Software Engineering",
	"CPS415": "Discrete Algoritms",
	"CPS425": "Language Processing",
	"CPS485": "Projects",
	"CPS342": "Embedded Linux", 		// elective
	"CPS440": "Database Principles",	// elective
	"CPS470": "Computer Networks",		// elective
	"CPS493": "Selected Topic",			// topic courses need restructuring in the server
	"MAT251": "Calculus I",
	"MAT252": "Calculus II",
	"MAT320": "Discrete Mathmatics",
	"EGC220": "Digital Logic",
}

type JsonSave struct {
	Msg   *discordgo.Message
	Roles []*discordgo.Role
}
