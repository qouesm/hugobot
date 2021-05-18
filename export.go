package main

import "gitlab.com/qouesm/hugobot/commands"

func exportCommands() []commands.Command {
	return []commands.Command{
		commands.ClassClear,
		commands.Options,
		commands.Ping,
		commands.Quietping,
		commands.ReactRoles,
		commands.Uhoh,
	}
}
