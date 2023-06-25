package Processor

import (
	"discordBot/events"
	"log"
	"strings"
)

const (
	HelpCmd = "/help"
)

func (p *Processor) doCmd(text string, metaMsg events.MetaMessage, fromClient string) error {
	text = strings.TrimSpace(text)

	switch text {
	case HelpCmd:
		return p.sendHelp(metaMsg)
	default:
		log.Printf("%s, %s wtite mesage: %s", fromClient, metaMsg.UserName, text)

		return nil
	}
}

func (p *Processor) sendHelp(metaMsg events.MetaMessage) error {
	return metaMsg.ReplyToSender(msgHelp, metaMsg.ChatID)
}
