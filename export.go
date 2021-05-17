package main

import . "gitlab.com/qouesm/hugobot/commands"

func exportCommands() []Command {
	return []Command{
		ClassClear,
		Options,
		Ping,
		Quietping,
		Uhoh,
	}
}
