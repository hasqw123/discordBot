package events

import "errors"

type Processor interface {
	Process(e Event) error
}

type Fetcher interface {
	Fetch() (Event, error)
}

type Type int

const (
	Unknown Type = iota
	Message
)

var (
	ErrUnknownEvent    = errors.New("unknown event type ")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

type Event struct {
	IsEvent bool
	Type    Type
	Text    string
	Meta    interface{}
}

type MetaMessage struct {
	ChatID        string
	UserName      string
	ReplyToSender func(text string, chatID string) error
}
