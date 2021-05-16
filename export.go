package main

import . "gitlab.com/qouesm/hugobot/commands"

func exportCommands() []Command {
	return []Command{
		Options,
		Ping,
		Quietping,
		Uhoh,
	}
}
