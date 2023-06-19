package telegram

import (
	"context"
	"discordBot/events"
	"discordBot/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type TgClient struct {
	host     string
	basePath string
	offset   int
	limit    int
	client   http.Client
	updts    []Update
	ctx      context.Context
	isClosed bool
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string, bathSize int) *TgClient {
	return &TgClient{
		host:     host,
		limit:    bathSize,
		basePath: newBasePath(token),
		client:   http.Client{},
		updts:    make([]Update, 0, bathSize),
		ctx:      context.Background(),
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *TgClient) updates(offset, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *TgClient) sendMessage(text string, chatId string) (err error) {
	defer func() { err = e.WrapIfErr("can't send message", err) }()

	q := url.Values{}
	q.Add("chat_id", chatId)
	q.Add("text", text)

	if _, err := c.doRequest(sendMessageMethod, q); err != nil {
		return err
	}

	return nil
}

func (c *TgClient) doRequest(method string, q url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()
	u := url.URL{
		Scheme:   "https",
		Host:     c.host,
		Path:     path.Join(c.basePath, method),
		RawQuery: q.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *TgClient) FetchUpdate() (events.Event, error) {
	select {
	case <-c.ctx.Done():
		return events.Event{}, events.NoEventsError
	default:
		if len(c.updts) != 0 {
			update := c.event(c.updts[0])
			c.updts = c.updts[1:]

			return update, nil
		}

		if c.isClosed {
			return events.Event{}, events.NoEventsError
		}

		updates, err := c.updates(c.offset, c.limit)
		if err != nil {
			return events.Event{}, e.Wrap("can't fetchUpdate", err)
		}

		if len(updates) == 0 {
			return events.Event{}, nil
		}

		update := c.event(updates[0])
		c.updts = updates[1:]
		c.offset = updates[len(updates)-1].ID + 1

		return update, nil
	}
}

func (c *TgClient) Close(ctx context.Context) error {
	c.isClosed = true
	c.ctx = ctx
	return nil
}

func (c *TgClient) event(update Update) events.Event {
	updType := fetchType(update)

	res := events.Event{
		IsEvent: true,
		Type:    updType,
		Text:    fetchText(update),
	}

	if updType == events.Message {
		res.Meta = events.MetaMessage{
			ChatID:        strconv.Itoa(update.Message.Chat.ID),
			UserName:      update.Message.From.UserName,
			ReplyToSender: c.sendMessage,
		}
	}

	return res
}

func fetchText(update Update) string {
	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}

func fetchType(update Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}

	return events.Message
}
