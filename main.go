package main

import (
	"discordBot/clients/telegram"
	event_consumer "discordBot/consumer/event-consumer"
	fetcher "discordBot/events/Fetcher"
	"discordBot/events/Processor"
	"flag"
	"log"
	"runtime"
)

func main() {

	//ctx, cancel := context.WithCancel(context.Background())

	//client := dscClient.New(mustToken())
	//client.SetupInterrupt(cancel, ctx)

	//eventProcessor := discord.New(client)

	//log.Printf("service started")

	//consumer := event_consumer.New(eventProcessor, eventProcessor, 100)

	//if err := consumer.Start(); err != nil {
	//	log.Fatal("service is stopped", err)
	//}

	//TODO: не доделал
	client := telegram.New("api.telegram.org", "6072205028:AAHmtmZo_9mdxvkxyDQ7HGsGoBOBnHV7jT8", 100)
	fethcer := fetcher.New(client)
	processor := Processor.New()
	consumer := event_consumer.New(fethcer, processor)
	consumer.Start(runtime.NumCPU())
}

func mustToken() string {

	token := flag.String(
		"dsc-bot-token",
		"",
		"token for access to discord bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
