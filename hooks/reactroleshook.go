package hooks

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func ReactRoles(s *discordgo.Session) {
	// read back from json
	// file, err := os.ReadFile("template.json")
	file, err := os.Open("hooks/reactrolesmessages.json")
	if err != nil {
		log.Println("problem reading file", err)
		return
	}

	// split into lines
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// add the hooks for each line (each message)
	for scanner.Scan() {
		line := scanner.Bytes()

		// json => go struct
		var save = new(JsonSave)
		err = json.Unmarshal(line, &save)
		if err != nil {
			log.Println("problem using unmarshal", err)
			continue
		}

		// role adding hook
		s.AddHandler(func(s *discordgo.Session, mr *discordgo.MessageReactionAdd) {
			if mr.UserID == s.State.User.ID {
				return
			}
			if mr.MessageID == save.Msg.ID {
				number := emojiNum[mr.Emoji.APIName()]
				role := save.Roles[number]
				err := s.GuildMemberRoleAdd(mr.GuildID, mr.UserID, role.ID)
				if err != nil {
					user, err1 := s.User(mr.UserID)
					if err1 != nil {
						log.Println("Couln't get user's struct,", err1)
						log.Println("Couldn't add role,", err)
						return
					}
					log.Println("Couldn't add role:", user.Username, ",", err)
					return
				}

				dm, err := s.UserChannelCreate(mr.UserID)
				if err != nil {
					log.Println("could not create dm channel,", err)
					return
				}
				g, err := s.Guild(mr.GuildID)
				if err != nil {
					log.Println("could not get current guild,", err)
					return
				}
				_, err = s.ChannelMessageSend(dm.ID, g.Name+": ADDED @"+role.Name)
				if err != nil {
					log.Println("could not dm user,", err)
					return
				}
			}
		})

		// role deletion hook
		s.AddHandler(func(s *discordgo.Session, mr *discordgo.MessageReactionRemove) {
			if mr.UserID == s.State.User.ID {
				return
			}
			if mr.MessageID == save.Msg.ID {
				number := emojiNum[mr.Emoji.APIName()]
				role := save.Roles[number]
				err := s.GuildMemberRoleRemove(mr.GuildID, mr.UserID, role.ID)
				if err != nil {
					user, err1 := s.User(mr.UserID)
					if err1 != nil {
						log.Println("Couln't get user's struct,", err1)
						log.Println("Couldn't add role,", err)
						return
					}
					log.Println("Couldn't del role:", user.Username, ",", err)
					return
				}

				dm, err := s.UserChannelCreate(mr.UserID)
				if err != nil {
					log.Println("could not create dm channel,", err)
					return
				}
				g, err := s.Guild(mr.GuildID)
				if err != nil {
					log.Println("could not get current guild,", err)
					return
				}
				_, err = s.ChannelMessageSend(dm.ID, g.Name+": REMOVED @"+role.Name)
				if err != nil {
					log.Println("could not dm user,", err)
					return
				}
			}
		})
	}
}

// returns int from the corresponding unicode emoji
var emojiNum = map[string]int{
	"0️⃣": 0,
	"1️⃣": 1,
	"2️⃣": 2,
	"3️⃣": 3,
	"4️⃣": 4,
	"5️⃣": 5,
	"6️⃣": 6,
	"7️⃣": 7,
	"8️⃣": 8,
	"9️⃣": 9,
}

type JsonSave struct {
	Msg   *discordgo.Message
	Roles []*discordgo.Role
}
