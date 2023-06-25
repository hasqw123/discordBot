package event_consumer

import (
	"context"
	"discordBot/events"
	"log"
	"sync"
)

type Consumer struct {
	baseCtx   context.Context
	fetcher   events.Fetcher
	processor events.Processor
	eventCh   chan events.Event
}

func New(ctx context.Context, fetcher events.Fetcher, processor events.Processor) Consumer {
	return Consumer{
		baseCtx:   ctx,
		fetcher:   fetcher,
		processor: processor,
		eventCh:   make(chan events.Event),
	}
}

func (c *Consumer) Start(amountHandlers int) error {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for i := 0; i < amountHandlers; i++ {
		wg.Add(1)
		go c.handleEvent(c.eventCh, &wg, i+1)
	}

	for {
		gotEvent, err := c.fetcher.Fetch(c.baseCtx)
		switch {
		case err == events.NoEventsError:
			close(c.eventCh)

			return err
		case err != nil:
			log.Printf("[ERR] consumer %s", err.Error())

			continue
		case gotEvent.IsEvent:
			c.eventCh <- gotEvent

			continue
		}
	}
}

func (c *Consumer) handleEvent(eventCh chan events.Event, wg *sync.WaitGroup, numberHandler int) {
	defer wg.Done()

	for event := range eventCh {
		if err := c.processor.Process(event); err != nil {
			log.Printf("[ERR] consumer %s", err.Error())
		}
	}

	log.Printf("handler %d finished", numberHandler)
}
