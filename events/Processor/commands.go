package Processor

import (
	"discordBot/events"
	"discordBot/lib/e"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	HelpCmd = "/help"
	GetRate = "/getRate"
)

func (p *Processor) doCmd(text string, metaMsg events.MetaMessage, fromClient string) error {
	text = strings.TrimSpace(text)

	switch text {
	case HelpCmd:
		return p.sendHelp(metaMsg)
	case GetRate:
		return p.getRate(metaMsg)
	default:
		log.Printf("%s, %s wtite mesage: %s", fromClient, metaMsg.UserName, text)

		return nil
	}
}

func (p *Processor) sendHelp(metaMsg events.MetaMessage) error {
	return metaMsg.ReplyToSender(msgHelp, metaMsg.ChatID)
}

func (p *Processor) getRate(msg events.MetaMessage) (err error) {
	defer func() { err = e.WrapIfErr("can't getRate", err) }()

	var rate rate

	resp, err := http.Get("https://quote.ru/api/v1/ticker/72413")
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &rate); err != nil {
		return err
	}

	res := fmt.Sprintf("Текущий курс доллара к рублю по ЦБ: %v 😱", rate.Data.Ticker.LastPrice)

	if err = msg.ReplyToSender(res, msg.ChatID); err != nil {
		return err
	}

	return nil
}

// TODO: rework
type rate struct {
	Data struct {
		Ticker struct {
			Type      string  `json:"type"`
			LastPrice float64 `json:"lastPrice"`
		} `json:"ticker"`
	} `json:"data"`
}
