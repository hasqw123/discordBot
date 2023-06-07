package discord

import (
	"log"
	"strings"
)

const (
	HelpCmd = "/help"
)

func (p *Processor) doCmd(text string, ChatID string, userName string) error {
	text = strings.TrimSpace(text)

	switch text {
	case HelpCmd:
		return p.sendHelp(ChatID)
	default:
		log.Printf("%s wtite mesage: %s", userName, text)
		return nil
	}
}
