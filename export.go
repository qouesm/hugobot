package main

import "gitlab.com/qouesm/hugobot/commands"

func exportCommands() []commands.Command {
	return []commands.Command{
		commands.ClassClear,
		// commands.Ping,
		// commands.Quietping,
		commands.ReactRoles,
		// commands.Uhoh,
	}
}
