package discord

import (
	"discordBot/clients/discord"
	"discordBot/events"
	"discordBot/lib/e"
	"errors"
)

type Processor struct {
	dsc *discord.Client
}

type Meta struct {
	ChatID   string
	UserName string
}

var (
	ErrUnknownMetaType = errors.New("unknown meta type")
	ErrUnknownEvent    = errors.New("unknown event type ")
)

func New(client *discord.Client) *Processor {
	return &Processor{
		dsc: client,
	}
}
func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.dsc.Updates(limit)
	if err == discord.ErrClose {
		return nil, err
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.WrapIfErr("can't process message", ErrUnknownEvent)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func (p *Processor) sendHelp(channelID string) error {
	return p.dsc.SendMessage(msgHelp, channelID)
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd discord.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: upd.Message,
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.ChannelID,
			UserName: upd.MessageAuthor,
		}
	}

	return res
}

func fetchType(upd discord.Update) events.Type {
	if upd.Type != 0 {
		return events.Unknown
	}
	return events.Message
}
