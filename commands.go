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
	log  log.Logger
	cmds map[string]*Command
	ids  map[string][]string
}

func NewCommands(l log.Logger) *Commands {
	if l == nil {
		l = log.NewNopLogger()
	}

	return &Commands{
		log: l,
		ids: map[string][]string{},
	}
}

func (c *Commands) TryAddCmd(h *Command) error {
	if h.Def == nil {
		return errors.Wrap(ErrInvalidApplicationCmd, "Cmd() must not return null")
	}
	if h.Def.Name == "" {
		return errors.Wrap(ErrInvalidApplicationCmd, "Cmd() must have a valid name")
	}
	if c.cmds == nil {
		c.cmds = map[string]*Command{}
	}

	_, ok := c.cmds[h.Def.Name]
	if ok {
		return errors.Wrap(ErrInvalidApplicationCmd, "two commands with the same name cannot be registered")
	}

	c.cmds[h.Def.Name] = h
	return nil
}

func (c *Commands) AddCmd(h *Command) {
	if err := c.TryAddCmd(h); err != nil {
		panic(err.Error())
	}
}

func (c *Commands) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := c.cmds[i.Data.Name]; ok {
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
	for k, v := range c.cmds {
		id, err := s.ApplicationCommandCreate(s.State.User.ID, g.Guild.ID, v.Def)
		if err != nil {
			l.Log("err", err, "msg", "Cannot register the commands handler")
		}
		l.Log("msg", "added handler", "name", k)

		c.ids[g.Guild.ID] = append(c.ids[g.Guild.ID], id.ID)
	}
}

func (c *Commands) TryRegisterHandler(s *discordgo.Session) error {
	s.AddHandler(c.OnGuildCreate)
	s.AddHandler(c.OnInteractionCreate)

	return nil
}

func (c *Commands) Close(s *discordgo.Session) error {
	for guild, ids := range c.ids {
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
