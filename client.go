package gabby

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Client is a Gabby client, which connects to a Gabby server
type Client struct {
	id           uuid.UUID
	domain       *Domain
	key          *Key    // see serverGabber note...
	serverGabber *Gabber // note that a Gabby server is a Gabber, just like any other Gabby client is
	ws           *websocket.Conn
	connected    bool
	doneChan     chan struct{}
	closeChan    chan chan struct{}
	sync.RWMutex
}

// ErrAlreadyConnected means we're already connected
var ErrAlreadyConnected = errors.New("already connected")

// NewClient returns a new Gabby client
func NewClient(id uuid.UUID, domain *Domain) *Client {
	c := &Client{
		id:     id,
		domain: domain,
	}
	return c
}

// Domain returns the Gabby domain
func (c *Client) Domain() *Domain {
	c.RLock()
	defer c.RUnlock()
	return c.domain
}

// Close closes the client's connection to the server
// This blocks until the connection is cleanly closed.
func (c *Client) Close() {
	c.Lock()
	defer c.Unlock()
	if !c.connected {
		return
	}
	if c.closeChan == nil {
		return
	}
	ackChan := make(chan struct{})
	c.closeChan <- ackChan
	<-ackChan
}

// Connect connects to a Gabby server
func (c *Client) Connect() error {
	c.Lock()
	defer c.Unlock()
	if c.connected {
		return ErrAlreadyConnected
	}

	c.Domain().updateDomainDetails() // ???: we might not want to always do this, for example on a reconnect

	u := c.Domain().wsURL.String()
	if u == "" {
		return ErrNoRecords
	}

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		return err
	}

	c.ws = ws
	c.connected = true

	go c.runWS()

	return nil
}

// runWS runs the Websocket connection inbound data and handles errors
func (c *Client) runWS() {
	defer c.ws.Close()
	go c.receiveWS()
	for {
		select {
		case ackChan := <-c.closeChan:
			err := c.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err == nil {
				select {
				case <-c.doneChan:
				case <-time.After(time.Second):
				}
			}
			c.connected = false // !!!: we're not c.Lock()'d, but we wrote to it anyway...
			ackChan <- struct{}{}
		}
	}
}

// receiveWS receives messages via Websocket
func (c *Client) receiveWS() {
	for {
		t, m, err := c.ws.ReadMessage()
		if err != nil {
			log.Println(err)
			c.Close()
			return // ???: is this the best way?
		}
		log.Println(t, m)
	}
}
