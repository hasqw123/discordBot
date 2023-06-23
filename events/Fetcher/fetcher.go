package fetcher

import (
	"context"
	"discordBot/clients"
	"discordBot/events"
	"discordBot/lib/e"
	"log"
	"time"
)

type Fetcher struct {
	clients     []clients.Client
	chForUpdate chan events.Event
	chError     chan error
	isStarted   bool
}

func New(bathSize int, cln ...clients.Client) *Fetcher {
	return &Fetcher{
		clients:     cln,
		chForUpdate: make(chan events.Event, bathSize),
		chError:     make(chan error),
	}
}

func (f *Fetcher) Fetch(ctx context.Context) (events.Event, error) {

	if !f.isStarted {
		go func() {
			defer func() {
				log.Println("get fetched updates stop")
				close(f.chError)
				close(f.chForUpdate)
			}()

			for {
				select {
				case <-ctx.Done():
					newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

					f.getLastEvents(newCtx)

					cancel()

					return
				default:
					f.getEvents()
				}
			}
		}()

		f.isStarted = true
	}

	select {
	case err := <-f.chError:
		if err == events.NoEventsError {

			return events.Event{}, err
		}

		return events.Event{}, e.Wrap("can't fetch update", err)
	case event := <-f.chForUpdate:

		return event, nil
	}
}

func (f *Fetcher) getEvents() {
	for _, client := range f.clients {
		event, err := client.FetchUpdate()
		if err != nil {
			f.chError <- err

			continue
		}

		if event.IsEvent {
			f.chForUpdate <- event
		}
	}
}

func (f *Fetcher) getLastEvents(ctx context.Context) {
	var amountClients, counterError = len(f.clients), 0

	for idx, client := range f.clients {
		err := client.Close(ctx)
		switch {
		case err == events.NoEventsError:
			f.clients = append(f.clients[:idx], f.clients[idx+1:]...)

			counterError++

			if counterError == amountClients {
				f.chError <- events.NoEventsError

				return
			}
		case err != nil:
			f.chError <- err

			continue
		}
	}

	for {
		for idx, client := range f.clients {
			event, err := client.FetchUpdate()
			switch {
			case err == events.NoEventsError:
				f.clients = append(f.clients[:idx], f.clients[idx+1:]...)

				counterError++

				if counterError == amountClients {
					f.chError <- events.NoEventsError

					return
				}
			case err != nil:
				f.chError <- err

				continue
			case event.IsEvent:
				f.chForUpdate <- event
			}
		}
	}
}
