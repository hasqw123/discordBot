package Processor

import (
	"discordBot/events"
	"discordBot/lib/e"
)

type Processor struct {
}

func New() *Processor {
	return &Processor{}
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", events.ErrUnknownEvent)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	metaMsg, err := metaMessage(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, metaMsg, event.FromClient); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func metaMessage(event events.Event) (events.MetaMessage, error) {
	res, ok := event.Meta.(events.MetaMessage)
	if !ok {
		return events.MetaMessage{}, e.Wrap("can't get meta message", events.ErrUnknownMetaType)
	}

	return res, nil
}
