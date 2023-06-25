package events

import (
	"context"
	"errors"
)

type Processor interface {
	Process(e Event) error
}

type Fetcher interface {
	Fetch(ctx context.Context) (Event, error)
}

type Type int

const (
	Unknown Type = iota
	Message
)

var (
	ErrUnknownEvent    = errors.New("unknown event type ")
	ErrUnknownMetaType = errors.New("unknown meta type")
	NoEventsError      = errors.New("events are over")
)

type Event struct {
	FromClient string
	IsEvent    bool
	Type       Type
	Text       string
	Meta       interface{}
}

type MetaMessage struct {
	ChatID        string
	UserName      string
	ReplyToSender func(text string, chatID string) error
}
