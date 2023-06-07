package event_consumer

import (
	"discordBot/clients/discord"
	"discordBot/events"
	"log"
	"runtime"
	"sync"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	bathSize  int
}

func New(fetcher events.Fetcher, processor events.Processor, bathSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		bathSize:  bathSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.bathSize)
		if err == discord.ErrClose {
			return err
		}

		if err != nil {
			log.Printf("[ERR] consumer %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {

		}
	}
}

func (c *Consumer) handleEvents(e []events.Event) error {
	chForHandle := make(chan events.Event)
	wg := sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for e := range chForHandle {
				if err := c.processor.Process(e); err != nil {
					log.Printf("can't handle: %s", err.Error())

					continue
				}
			}
			wg.Done()
		}()
		wg.Add(1)
	}

	for _, event := range e {
		chForHandle <- event
	}
	close(chForHandle)
	wg.Wait()

	return nil
}
