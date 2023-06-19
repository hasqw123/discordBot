package main

import (
	"context"
	"discordBot/clients/telegram"
	"discordBot/consumer/event-consumer"
	fetcher "discordBot/events/Fetcher"
	"discordBot/events/Processor"
	"flag"
	"log"
	"os/signal"
	"syscall"
)

// TODO: это все нужно будет в аконфиг
const (
	tgBotHost                 = "api.telegram.org"
	batchSize                 = 100
	amountGoRoutineForHandler = 6
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client := telegram.New(tgBotHost, "6072205028:AAHmtmZo_9mdxvkxyDQ7HGsGoBOBnHV7jT8", batchSize)

	ftr := fetcher.New(batchSize, client)

	processor := Processor.New()

	log.Println("service started")

	consumer := event_consumer.New(ctx, ftr, processor)
	err := consumer.Start(amountGoRoutineForHandler)
	if err != nil {
		log.Println(err)
	}

	log.Println("exit...")
	log.Println("service stopped")
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
