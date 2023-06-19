package event_consumer

import (
	"discordBot/events"
	"log"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	eventCh   chan events.Event
}

func New(fetcher events.Fetcher, processor events.Processor) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		eventCh:   make(chan events.Event),
	}
}

func (c *Consumer) Start(amountHandlers int) error {
	for i := 0; i < amountHandlers; i++ {
		go c.handleEvents(c.eventCh)
	}

	for {
		gotEvent, err := c.fetcher.Fetch()

		if err != nil {
			log.Printf("[ERR] consumer %s", err.Error())

			continue
		}

		c.eventCh <- gotEvent
	}
}

func (c *Consumer) handleEvents(eventCh chan events.Event) {
	for event := range eventCh {
		if err := c.processor.Process(event); err != nil {
			log.Printf("[ERR] consumer %s", err.Error())
		}
	}
}

//func (c *Consumer) handleEvents(e events.Event) error {
//	chForHandle := make(chan events.Event)
//	wg := sync.WaitGroup{}
//
//	for i := 0; i < runtime.NumCPU(); i++ {
//		go func() {
//			for e := range chForHandle {
//				if err := c.processor.Process(e); err != nil {
//					log.Printf("can't handle: %s", err.Error())
//
//					continue
//				}
//			}
//			wg.Done()
//		}()
//		wg.Add(1)
//	}
//
//	for _, event := range e {
//		chForHandle <- event
//	}
//	close(chForHandle)
//	wg.Wait()
//
//	return nil
//}
