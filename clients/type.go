package clients

import (
	"context"
	"discordBot/events"
)

type Client interface {
	FetchUpdate() (events.Event, error)
	Close(ctx context.Context) error
}
