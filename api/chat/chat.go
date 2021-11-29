package chat

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nats-io/nats.go"
)

type Chat struct {
	nc    *nats.Conn
	state []Msg
}

//go:embed tester.creds
var testerCreds string

func New() *Chat {
	return &Chat{}
}

func (c *Chat) Default(ctx context.Context) error {
	panic("not implemented")
}

type Msg struct {
	ID      int    `json:"id"`
	Message string `json:"m"`
}

// Post accepts client post message request
// Publishes message to the chat subject and stores that message to the state.
func (c *Chat) Post(ctx context.Context, msg Msg) error {
	if msg.Message == "" {
		return fmt.Errorf("no message")
	}
	if c.nc == nil {
		nc, err := natsConnect(testerCreds)
		if err != nil {
			return err
		}
		c.nc = nc
	}
	if err := c.pub("chat", msg); err != nil {
		return err
	}
	if err := c.flush(); err != nil {
		return err
	}
	c.state = append(c.state, msg)
	return nil
}

func (c *Chat) pub(subject string, msg Msg) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to marshal %s", err)
		return internalServerError
	}
	if err := c.nc.Publish(subject, buf); err != nil {
		log.Printf("publish error %s", err)
		return internalServerError
	}
	return nil
}

func (c *Chat) flush() error {
	if err := c.nc.Flush(); err != nil {
		log.Printf("flush error %s", err)
		return internalServerError
	}
	return nil
}

// State publishes to the client inbox all messages from the state.
func (c *Chat) State(ctx context.Context, inbox string) error {
	if c.nc == nil {
		nc, err := natsConnect(testerCreds)
		if err != nil {
			return err
		}
		c.nc = nc
	}

	// publish current chat state
	for _, msg := range c.state {
		if err := c.pub(inbox, msg); err != nil {
			return err
		}
	}

	// publish one empty message to signal end of stream
	if err := c.nc.Publish(inbox, nil); err != nil {
		log.Printf("publish error %s", err)
		return internalServerError
	}

	if err := c.flush(); err != nil {
		return err
	}
	return nil
}

var internalServerError = fmt.Errorf("internal server error")

func natsConnect(userJWT string) (*nats.Conn, error) {
	url := "connect.ngs.global"

	// save credentials from emebeded string variable to file
	file, err := ioutil.TempFile("", "*.creds")
	if err != nil {
		log.Printf("failed to save creds to tmp file %s", err)
		return nil, internalServerError
	}
	if _, err = fmt.Fprint(file, testerCreds); err != nil {
		log.Printf("failed to write to tmp file %s %s", file.Name(), err)
		return nil, internalServerError
	}

	conn, err := nats.Connect(url, nats.UserCredentials(file.Name()))
	if err != nil {
		log.Printf("failed to connect: %s", err)
		return nil, internalServerError
	}
	return conn, nil
}
