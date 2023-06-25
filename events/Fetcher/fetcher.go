package fetcher

import (
	"context"
	"discordBot/clients"
	"discordBot/events"
	"discordBot/lib/e"
	"log"
	"reflect"
	"sync"
	"time"
)

type Fetcher struct {
	clients              []clients.Client
	chForEvent           chan events.Event
	chError              chan error
	isStarted            bool
	counterNoEventsError int
	wg                   sync.WaitGroup
}

func New(bathSize int, cln ...clients.Client) *Fetcher {
	return &Fetcher{
		clients:    cln,
		chForEvent: make(chan events.Event, bathSize),
		chError:    make(chan error, 1),
	}
}

func (f *Fetcher) Fetch(ctx context.Context) (events.Event, error) {
	if !f.isStarted {
		for i := 0; i < len(f.clients); i++ {
			client := f.clients[i]

			f.wg.Add(1)
			go func(client clients.Client) {
				defer func() {
					log.Printf("fetcher: getEvents for %s stopped", reflect.TypeOf(client).String())
					f.wg.Done()
				}()

				for {
					select {
					case <-ctx.Done():
						newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

						err := client.Close(newCtx)
						switch err {
						case events.NoEventsError:
							f.chError <- err

							cancel()

							return
						case nil:
							f.getLastEvents(client)

							cancel()

							return
						default:
							f.chError <- err

							continue
						}

					default:
						event, err := client.FetchUpdate()
						if err != nil {
							f.chError <- err

							continue
						}

						if event.IsEvent {
							f.chForEvent <- event
						}
					}
				}
			}(client)
		}

		f.isStarted = true
	}

	select {
	case err := <-f.chError:
		if err == events.NoEventsError {
			f.counterNoEventsError++

			if f.counterNoEventsError == len(f.clients) {
				f.wg.Wait()

				close(f.chForEvent)
				log.Println("fetcher: chForEvent from fetcher closed")

				return events.Event{}, e.Wrap("got shutdown signal from all client", err)
			}

			return events.Event{}, e.Wrap("got shutdown signal from client", err)
		}

		return events.Event{}, e.Wrap("can't fetch update", err)
	case event, opened := <-f.chForEvent:
		if !opened {
			close(f.chError)
			log.Println("fetcher: chForError from fetcher closed")

			return events.Event{}, events.NoEventsError
		}

		return event, nil
	}
}

func (f *Fetcher) getLastEvents(client clients.Client) {
	for {
		event, err := client.FetchUpdate()
		switch err {
		case events.NoEventsError:
			f.chError <- err

			return
		case nil:
			if event.IsEvent {
				f.chForEvent <- event
			}

			continue
		default:
			f.chError <- err

			continue
		}
	}
}
