package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
			mRoles    []*discordgo.Role
		)

		switch i.Data.Options[0].Name {
		case "create":
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
				mRoles = append(mRoles, role)
			}
			wgIns.Done()

			log.Println("printing mRoles")
			for i, v := range mRoles {
				log.Println(i, v.Name)
			}
			log.Println("lets try gettings ints from the map")
			log.Println(emojiNum[&discordgo.Emoji{Name: "2️⃣"}])
			log.Println(emojiNum)
			for k, v := range emojiNum {
				log.Println(k, v)
			}

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

			for num := 0; num < len(i.Data.Options[0].Options)-1; num++ {
				s.MessageReactionAdd(rMsg.ChannelID, rMsg.ID, numEmoji[num].APIName())
			}

			s.AddHandler(func(s *discordgo.Session, mr *discordgo.MessageReactionAdd) {
				if mr.UserID == s.State.User.ID {
					return
				}
				if mr.MessageID == rMsg.ID {
					log.Println("add:", mr.Emoji.APIName())

					number := emojiNum[&mr.Emoji]
					log.Println(number)
					role := mRoles[number]

					log.Println(number, mr.Emoji.APIName(), &mr.Emoji.Name, role.Name)
					err := s.GuildMemberRoleAdd(mr.GuildID, mr.UserID, role.ID)
					if err != nil {
						log.Println("Couldn't add role:", mr.Emoji.User.Username, ",", err)
					}
				}
			})

			s.AddHandler(func(s *discordgo.Session, mr *discordgo.MessageReactionRemove) {
				if mr.UserID == s.State.User.ID {
					return
				}
				if mr.MessageID == rMsg.ID {
					log.Println("del:", mr.Emoji.APIName())

					number := emojiNum[&mr.Emoji]
					log.Println(number)
					role := mRoles[number]

					log.Println(number, mr.Emoji.APIName(), &mr.Emoji.Name, role.Name)
					err := s.GuildMemberRoleRemove(mr.GuildID, mr.UserID, role.ID)
					if err != nil {
						log.Println("Couldn't del role:", mr.Emoji.User.Username, ",", err)
					}
				}
			})

			file, err := os.Create("template.json")
			if err != nil {
				log.Println("problem creating file,", err)
			}
			defer file.Close()
			js, err := json.Marshal(rMsg)
			if err != nil {
				log.Println("problem using marshal,", err)
			}
			file.Write(js)

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

var emojiNum = map[*discordgo.Emoji]int{
	{Name: "0️⃣"}: 0,
	{Name: "1️⃣"}: 1,
	{Name: "2️⃣"}: 2,
	{Name: "3️⃣"}: 3,
	{Name: "4️⃣"}: 4,
	{Name: "5️⃣"}: 5,
	{Name: "6️⃣"}: 6,
	{Name: "7️⃣"}: 7,
	{Name: "8️⃣"}: 8,
	{Name: "9️⃣"}: 9,
}
