package commands

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var ClassClear = Command{
	AppCommand: discordgo.ApplicationCommand{
		Name:        "classclear",
		Description: "ADMIN: Remove all class roles from all server members",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		curGuild, err := s.Guild(i.GuildID)
		if err != nil {
			log.Println("Problem getting guild struct from Interaction,", err)
			panicResponse(s, i)
			return
		}

		// if "command caller" does not have role "Admin", return
		// if !hasAdmin(i.Member.Roles) {
		if !hasAdmin(i.Member.Roles, s.State, curGuild.ID) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionApplicationCommandResponseData{
					Flags:   1 << 6, // only visible to caller
					Content: "You are not allowed to use that command",
				},
			})
			return
		}

		curGuild.Members, err = s.GuildMembers(curGuild.ID, "0", 200)
		if err != nil {
			log.Println("Problem getting members from current guild,", err)
			panicResponse(s, i)
		}

		// match department call name + 3 numbers
		reClass, err := regexp.Compile("(CPS|MAT|EGC)[0-9]{3}")
		if err != nil {
			log.Println("Bad regex,", err)
			panicResponse(s, i)
			return
		}

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

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Roles removed successfully",
			},
		})
	},
}

func panicResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "Something went wrong",
		},
	})
}

// passing the state and guild ID is sloppy but Member.roles only returns role ID's and not the structs

func hasAdmin(slice []string, state *discordgo.State, gID string) bool {
	for _, v := range slice {
		r, _ := state.Role(gID, v)
		if strings.EqualFold(r.Name, "admin") {
			return true
		}
	}
	return false
}
