package fetcher

import (
	"discordBot/clients"
	"discordBot/events"
	"discordBot/lib/e"
)

type Fetcher struct {
	clients     []clients.Client
	chForUpdate chan events.Event
	chError     chan error
	isStarted   bool
}

func New(cln ...clients.Client) *Fetcher {
	return &Fetcher{
		clients:     cln,
		chForUpdate: make(chan events.Event, 50), //TODO: подумать как передовать размер канала сюда
		chError:     make(chan error),            //TODO: и здесь тоже подумать
	}
}

func (f *Fetcher) Fetch() (events.Event, error) {
	//TODO: нужно еще проудмать этот момент

	if !f.isStarted {
		go func() {
			for {
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
		}()

		f.isStarted = true
	}

	select {
	case err := <-f.chError:
		return events.Event{}, e.Wrap("can't fetch update", err)
	case event := <-f.chForUpdate:
		return event, nil
	}
}
