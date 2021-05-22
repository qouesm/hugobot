package commands

import (
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var ClassClear = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "classclear",
		Description: "ADMIN: Remove all class roles from all server members",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Println("classclear")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{},
		})

		curGuild, err := s.Guild(i.GuildID)
		if err != nil {
			log.Println("Problem getting guild struct from Interaction,", err)
			panicResponse(s, i)
			return
		}

		// if "command caller" does not have role "Admin", return
		if !hasAdmin(s.State, i.Member.Roles, curGuild.ID) {
			log.Println("no admin")
			_, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{Content: "You are not allowed to use that command"})
			if err != nil{
				log.Println("problem creating follow up message,", err)
				return
			}
			return
		}
		log.Println("admin")

		curGuild.Members, err = s.GuildMembers(curGuild.ID, "0", 200)
		if err != nil {
			log.Println("Problem getting members from current guild,", err)
			panicResponse(s, i)
		}

		// match department call name + 3 numbers
		reClass, _ := regexp.Compile("(CPS|MAT|EGC)[0-9]{3}")

		for _, member := range curGuild.Members {
			for _, roleID := range member.Roles {
				role, err := s.State.Role(curGuild.ID, roleID)
				if err != nil {
					log.Println("Problem getting role struct,", err)
					continue
				}
				if reClass.MatchString(role.Name) {
					err := s.GuildMemberRoleRemove(curGuild.ID, member.User.ID, roleID)
					if err != nil {
						log.Println("Problem removing role,", err)
					}
				}
			}
		}

		_, err = s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
			Content: "Roles removed successfully",
		})
		if err != nil {
			log.Println("problem creating followup,", err)
		}

	},
}
