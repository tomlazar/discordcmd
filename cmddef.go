package discordcmd

import "github.com/bwmarrin/discordgo"

type Command struct {
	Def *discordgo.ApplicationCommand
	Run func(Context)
}
