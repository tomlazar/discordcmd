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
	s    *discordgo.Session
	i    *discordgo.InteractionCreate
	icmd *discordgo.ApplicationCommandInteractionData
}

func newContextFromInteraction(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	icmd *discordgo.ApplicationCommandInteractionData,
) Context {
	return &context{s: s, i: i, icmd: icmd}
}

func (c *context) Ack() error {
	return c.s.InteractionRespond(c.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func (c *context) Reply(str string) error {
	return c.s.InteractionRespond(c.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: str,
		},
	})
}
func (c *context) Embed(embed ...*discordgo.MessageEmbed) error {
	return c.s.InteractionRespond(c.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embed,
		},
	})
}

func (c *context) Session() *discordgo.Session                     { return c.s }
func (c *context) InteractionCreate() *discordgo.InteractionCreate { return c.i }
func (c *context) Member() *discordgo.Member                       { return c.i.Member }

func (c *context) String(key string) string {
	for _, kv := range c.icmd.Options {
		if kv.Name == key {
			return kv.StringValue()
		}
	}

	return ""
}
