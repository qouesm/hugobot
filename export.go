package main

import . "gitlab.com/qouesm/qbot/commands"

func exportCommands() []Command {
	return []Command{
		Options,
		Ping,
		Quietping,
		Uhoh,
	}
}
