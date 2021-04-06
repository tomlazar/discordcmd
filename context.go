package discordcmd

import "github.com/bwmarrin/discordgo"

type Context interface {
	Ack() error
	Reply(string) error
	Embed(...*discordgo.MessageEmbed) error
	Session() *discordgo.Session
	Member() *discordgo.Member
	InteractionCreate() *discordgo.InteractionCreate
	String(key string) string
}

type context struct {
	s *discordgo.Session
	i *discordgo.InteractionCreate
}

func newContextFromInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) Context {
	return &context{s: s, i: i}
}

func (c *context) Ack() error {
	return c.s.InteractionRespond(c.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseAcknowledge,
	})
}

func (c *context) Reply(str string) error {
	return c.s.InteractionRespond(c.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: str,
		},
	})
}
func (c *context) Embed(embed ...*discordgo.MessageEmbed) error {
	return c.s.InteractionRespond(c.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: embed,
		},
	})
}

func (c *context) Session() *discordgo.Session                     { return c.s }
func (c *context) InteractionCreate() *discordgo.InteractionCreate { return c.i }
func (c *context) Member() *discordgo.Member                       { return c.i.Member }

func (c *context) String(key string) string {
	for _, kv := range c.i.Data.Options {
		if kv.Name == key {
			return kv.StringValue()
		}
	}

	return ""
}
