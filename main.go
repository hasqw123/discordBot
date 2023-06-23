package main

import (
	"context"
	"discordBot/clients/telegram"
	"discordBot/consumer/event-consumer"
	fetcher "discordBot/events/Fetcher"
	"discordBot/events/Processor"
	"discordBot/lib/configs"
	"discordBot/lib/e"
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

func main() {
	config := loadConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	tgClient := telegram.New(config.TgBot.Host, config.TgBot.Token, config.BatchSize)

	ftr := fetcher.New(config.BatchSize, tgClient)

	processor := Processor.New()

	log.Println("service started")

	consumer := event_consumer.New(ctx, ftr, processor)
	err := consumer.Start(config.AmountHandler)
	if err != nil {
		log.Println(err)
	}

	log.Println("exit...")
	log.Println("service stopped")
}

func loadConfig() configs.Config {
	config := configs.New()

	yamlConfig, err := os.Open("lib/configs/app.yaml")
	if err != nil {
		log.Fatal("can't load config", err)
	}
	defer func() { _ = yamlConfig.Close() }()

	data, err := io.ReadAll(yamlConfig)
	if err != nil {
		log.Fatal("can't load config", err)
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(e.Wrap("can't load config", err).Error())
	}

	if err = configCheck(config); err != nil {
		log.Fatal(e.Wrap("can't load config", err).Error())
	}

	return config
}

func configCheck(config interface{}) error {
	err := errors.New("all config fields are required")

	value := reflect.ValueOf(config)
	typeValue := value.Type()
	amount := typeValue.NumField()

	for i := 0; i < amount; i++ {
		v := value.Field(i)
		kv := v.Kind()

		switch kv {
		case reflect.Int:
			if v.Interface() == 0 {
				return err
			}
		case reflect.String:
			if v.Interface() == "" {
				return err
			}
		default:
			if err := configCheck(v.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
