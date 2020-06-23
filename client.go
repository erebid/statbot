package statbot

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/state"
	"github.com/jackc/pgx/v4/pgxpool"
)

const schema = `CREATE TABLE IF NOT EXISTS messages (
  id bigint NOT NULL,
  sent timestamp UNIQUE,
  content text,
  channel_id bigint,
  channel_name text,
  author_id bigint,
  author_name text,
  server_id bigint,
  server_name text,
  PRIMARY KEY (id)
)`

type Client struct {
	*state.State
	db *pgxpool.Pool
}

func NewClient(state *state.State, db *pgxpool.Pool) (*Client, error) {
	if db == nil {
		return nil, errors.New("db must not be nil")
	}
	c := &Client{
		State: state,
		db:    db,
	}
	_, err := c.db.Exec(context.Background(), schema)
	if err != nil {
		return nil, err
	}
	c.AddHandler(c.handleMessage)
	return c, nil
}

func (c *Client) handleMessage(m *gateway.MessageCreateEvent) {
	if m.Author.Bot {
		return
	}
	var guildName string
	if g, err := c.Guild(m.GuildID); err == nil {
		guildName = g.Name
	} else {
		log.Println(err)
		return
	}
	var channelName string
	if ch, err := c.Channel(m.ChannelID); err == nil {
		channelName = ch.Name
	} else {
		log.Println(err)
		return
	}
	authorTag := m.Author.Username + "#" + m.Author.Discriminator
	_, err := c.db.Exec(context.Background(),
		`insert into messages(
		id, sent, content,
		channel_id, channel_name, author_id, author_name,
		server_id, server_name
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		m.ID, time.Now(), m.Content,
		m.ChannelID, channelName,
		m.Author.ID, authorTag,
		m.GuildID, guildName)
	if err != nil {
		log.Println(err)
		return
	}
}
