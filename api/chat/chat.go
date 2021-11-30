package chat

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mantil-io/mantil.go/er"
	"github.com/nats-io/nats.go"
)

type Chat struct {
	nc    *nats.Conn
	state []Msg
}

//go:embed tester.creds
var testerCreds string

func New() *Chat {
	nc, err := natsConnect(testerCreds)
	if err != nil {
		panic(fmt.Errorf("nats connection failed %w", err))
	}
	return &Chat{
		nc: nc,
	}
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
		return er.E(er.NewBadRequestError("no message"))
	}
	if err := c.pub("chat", msg); err != nil {
		return er.E(err, er.ErrInternalServer)
	}
	if err := c.flush(); err != nil {
		return er.E(err, er.ErrInternalServer)
	}
	c.state = append(c.state, msg)
	return nil
}

// State publishes to the client inbox all messages from the state.
func (c *Chat) State(ctx context.Context, inbox string) error {
	// publish current chat state
	for _, msg := range c.state {
		if err := c.pub(inbox, msg); err != nil {
			return er.E(err, er.ErrInternalServer)
		}
	}

	// publish one empty message to signal end of stream
	if err := c.nc.Publish(inbox, nil); err != nil {
		return er.E(err, er.ErrInternalServer)
	}

	if err := c.flush(); err != nil {
		return er.E(err, er.ErrInternalServer)
	}
	return nil
}

func (c *Chat) pub(subject string, msg Msg) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return er.E(err)
	}
	if err := c.nc.Publish(subject, buf); err != nil {
		return er.E(err)
	}
	return nil
}

func (c *Chat) flush() error {
	return er.E(c.nc.Flush())
}

func natsConnect(userJWT string) (*nats.Conn, error) {
	url := "connect.ngs.global"

	credsFile, err := saveEmebedeCredsToTmpFile()
	if err != nil {
		return nil, er.E(err)
	}

	conn, err := nats.Connect(url, nats.UserCredentials(credsFile))
	if err != nil {
		return nil, er.E(err)
	}
	return conn, nil
}

func saveEmebedeCredsToTmpFile() (string, error) {
	file, err := ioutil.TempFile("", "*.creds")
	if err != nil {
		return "", er.E(err)
	}
	if _, err = fmt.Fprint(file, testerCreds); err != nil {
		return "", er.E(err)
	}
	return file.Name(), nil
}
