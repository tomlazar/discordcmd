package discordcmd

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	"github.com/go-kit/kit/log"
)

var (
	ErrDuplicateCmd          = errors.New("a command with the same name cannot be added twice")
	ErrInvalidApplicationCmd = errors.New("the application command is not valid")
)

type Commands struct {
	log        log.Logger
	handlers   map[string]CmdDef
	commandIds map[string][]string
}

func NewCommands(l log.Logger) *Commands {
	if l == nil {
		l = log.NewNopLogger()
	}

	return &Commands{
		log:        l,
		commandIds: map[string][]string{},
	}
}

func (c *Commands) TryAddCmd(h CmdDef) error {
	cmd := h.Cmd()
	if cmd == nil {
		return errors.Wrap(ErrInvalidApplicationCmd, "Cmd() must not return null")
	}
	if cmd.Name == "" {
		return errors.Wrap(ErrInvalidApplicationCmd, "Cmd() must have a valid name")
	}
	if c.handlers == nil {
		c.handlers = map[string]CmdDef{}
	}

	_, ok := c.handlers[cmd.Name]
	if ok {
		return errors.Wrap(ErrInvalidApplicationCmd, "two commands with the same name cannot be registered")
	}

	c.handlers[cmd.Name] = h
	return nil
}

func (c *Commands) AddCmd(h CmdDef) {
	if err := c.TryAddCmd(h); err != nil {
		panic(err.Error())
	}
}

func (c *Commands) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := c.handlers[i.Data.Name]; ok {
		h.Run(newContextFromInteraction(s, i))
	}
}

func (c *Commands) OnGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	l := log.With(c.log, "guild", g.Guild.Name)

	l.Log("msg", "adding guild commands")

	// delete duplicates
	cmds, err := s.ApplicationCommands(s.State.User.ID, g.Guild.ID)
	if err != nil {
		l.Log("err", err, "msg", "cannot list application commands")
		return
	}

	for _, cmd := range cmds {
		l.Log("err", err, "msg", "deleting application cmd ", "name", cmd.Name)
		err = s.ApplicationCommandDelete(s.State.User.ID, g.Guild.ID, cmd.ID)
		if err != nil {
			l.Log("err", err, "msg", "Cannot register the commands handler")
		}
	}

	// ensure new commands
	for k, v := range c.handlers {
		id, err := s.ApplicationCommandCreate(s.State.User.ID, g.Guild.ID, v.Cmd())
		if err != nil {
			l.Log("err", err, "msg", "Cannot register the commands handler")
		}
		l.Log("msg", "added handler", "name", k)

		c.commandIds[g.Guild.ID] = append(c.commandIds[g.Guild.ID], id.ID)
	}
}

func (c *Commands) TryRegisterHandler(s *discordgo.Session) error {
	s.AddHandler(c.OnGuildCreate)
	s.AddHandler(c.OnInteractionCreate)

	return nil
}

func (c *Commands) Close(s *discordgo.Session) error {
	for guild, ids := range c.commandIds {
		for _, id := range ids {
			err := s.ApplicationCommandDelete(s.State.User.ID, guild, id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Commands) RegisterHandler(s *discordgo.Session) {
	if err := c.TryRegisterHandler(s); err != nil {
		panic(err.Error())
	}
}
