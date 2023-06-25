package discord

import (
	"context"
	"discordBot/events"
	"discordBot/lib/e"
	"github.com/bwmarrin/discordgo"
	"log"
	"sync"
	"time"
)

type DscClient struct {
	isCloses bool
	conn     *discordgo.Session
	ctx      context.Context
	mutex    sync.Mutex
}

const dscClient = "Discord"

var (
	chUpd chan Update
)

func New(token string, channelSize int) *DscClient {
	sess, err := discordgo.New(formatToken(token))
	if err != nil {
		log.Fatal("can't creating discord session:", err)
	}

	if err = sess.Open(); err != nil {
		log.Fatal("can't open ws conn:", err)
	}

	chUpd = make(chan Update, channelSize)

	sess.AddHandler(handler)

	return &DscClient{
		conn: sess,
		ctx:  context.Background(),
	}
}

func handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	u := Update{
		Type:          m.Type,
		MessageAuthor: m.Author.Username,
		Message:       m.Content,
		ChannelID:     m.ChannelID,
	}

	chUpd <- u
}

func (c *DscClient) sendMessage(text string, channelID string) error {
	if c.isCloses {
		c.mutex.Lock()

		_ = c.conn.Open()

		defer func() {
			_ = c.conn.Close()

			c.mutex.Unlock()
		}()
	}

	_, err := c.conn.ChannelMessageSend(channelID, text)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func formatToken(token string) string {
	return "Bot " + token
}

func (c *DscClient) FetchUpdate() (events.Event, error) {
	select {
	case <-c.ctx.Done():
		return events.Event{}, events.NoEventsError
	case u, ok := <-chUpd:
		if ok {
			event := c.event(u)

			return event, nil
		}

		if c.isCloses {
			return events.Event{}, events.NoEventsError
		}

		return events.Event{}, nil
	case <-time.After(2 * time.Second):
		return events.Event{}, nil
	}
}

func (c *DscClient) Close(ctx context.Context) error {
	c.isCloses = true
	defer func() { c.ctx = ctx }()

	if err := c.conn.Close(); err != nil {
		return e.Wrap("can't close connection discord", err)
	}
	close(chUpd)

	log.Println("discord: channel for updates and discord connection are closed")

	if len(chUpd) == 0 {
		return events.NoEventsError
	}

	return nil
}

func (c *DscClient) event(update Update) events.Event {
	updType := fetchType(update)

	res := events.Event{
		FromClient: dscClient,
		IsEvent:    true,
		Type:       updType,
		Text:       update.Message,
	}

	if updType == events.Message {
		res.Meta = events.MetaMessage{
			ChatID:        update.ChannelID,
			UserName:      update.MessageAuthor,
			ReplyToSender: c.sendMessage,
		}
	}

	return res
}

func fetchType(update Update) events.Type {
	if update.Type == 0 {
		return events.Message
	}

	return events.Unknown
}
