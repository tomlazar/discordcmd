package discordcmd

import "github.com/bwmarrin/discordgo"

type Handler struct {
	Command *discordgo.ApplicationCommand
	RunFunc func(c Context)
}

func (h *Handler) Cmd() *discordgo.ApplicationCommand {
	return h.Command
}

func (h *Handler) Run(c Context) {
	h.RunFunc(c)
}
