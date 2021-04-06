package discordcmd

import "github.com/bwmarrin/discordgo"

type CmdDef interface {
	Cmd() *discordgo.ApplicationCommand
	Run(ctx Context)
}
