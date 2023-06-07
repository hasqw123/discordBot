package discord

import (
	"context"
	"discordBot/lib/e"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Client struct {
	conn *discordgo.Session
}

func New(token string) *Client {
	sess, err := discordgo.New(formatToken(token))
	if err != nil {
		log.Fatal("can't creating discord session:", err)
	}

	u, err := sess.User("@me")
	if err != nil {
		log.Fatal("can't get botId: ", err)
	}
	botID = u.ID
	chUpd = make(chan Update, sizeCh)
	chSend = make(chan string, sizeCh)

	if err = sess.Open(); err != nil {
		log.Fatal("can't open ws conn:", err)
	}

	sess.AddHandler(handler)

	return &Client{
		conn: sess,
	}
}

func handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botID {
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

func (c *Client) SendMessage(text string, channelID string) error {
	_, err := c.conn.ChannelMessageSend(channelID, text)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) CloseConn() error {
	if err := c.conn.Close(); err != nil {
		return e.WrapIfErr("can't close", err)
	}
	close(chUpd)
	return nil
}

func (c *Client) Updates(limit int) ([]Update, error) {
	firstValue, open := <-chUpd
	if !open {
		return nil, ErrClose
	}

	res := make([]Update, 0, limit)
	res = append(res, firstValue)

	t := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-t.C:
			return res, nil
		case u := <-chUpd:
			res = append(res, u)
			if len(res) == limit {
				t.Stop()
				return res, nil
			}
		}
	}

}

func (c *Client) SetupInterrupt(cancel context.CancelFunc, ctx context.Context) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
	}()

	go func() {
		<-ctx.Done()
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err := c.CloseConn()
		if err != nil {
			log.Println("err to close :", err)
		}
		log.Println("conn close")
	}()
}

func formatToken(token string) string {
	return "Bot " + token
}
