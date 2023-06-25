package discord

import "github.com/bwmarrin/discordgo"

type Update struct {
	Type          discordgo.MessageType
	MessageAuthor string
	Message       string
	ChannelID     string
}
