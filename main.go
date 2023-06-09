package main

import (
	"context"
	dscClient "discordBot/clients/discord"
	"discordBot/consumer/event-consumer"
	"discordBot/events/discord"
	"flag"
	"log"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	client := dscClient.New(mustToken())
	client.SetupInterrupt(cancel, ctx)

	eventProcessor := discord.New(client)

	log.Printf("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, 100)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

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
