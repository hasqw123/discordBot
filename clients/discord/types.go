package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

type Update struct {
	Type          discordgo.MessageType
	MessageAuthor string
	Message       string
	ChannelID     string
}

var (
	botID    string
	chUpd    chan Update
	chSend   chan string
	ErrClose = errors.New("channel cLose")
)

const sizeCh = 100
