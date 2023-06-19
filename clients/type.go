package clients

import "discordBot/events"

type Client interface {
	FetchUpdate() (events.Event, error)
}
