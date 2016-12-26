package gabby

import (
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Conversation is a conversation, including the session key and the participants
type Conversation struct {
	id         uuid.UUID
	key        *Key
	members    map[uuid.UUID]*Gabber
	expiration time.Time
	sync.RWMutex
}

// NewConversation returns a conversation for the provided members, key, expiration
func NewConversation(id uuid.UUID, key *Key, members map[uuid.UUID]*Gabber, expiration time.Time) *Conversation {
	c := &Conversation{
		id:         id,
		key:        key,
		members:    members,
		expiration: expiration,
	}
	for _, member := range c.members {
		member.addConversation(c)
	}
	return c
}

// Destroy destroys the conversation and references to it in the member Gabbers
func (c *Conversation) Destroy() {
	c.Lock()
	defer c.Unlock()
	c.expiration = time.Now()
	c.key = nil
	for _, v := range c.members {
		v.removeConversation(c)
	}
}

// Key returns the session key for this conversation
// I'm hoping we can not provide this and do it under the covers instead.
func (c *Conversation) Key() *Key {
	c.RLock()
	defer c.RUnlock()
	if c.key.Expired() {
		return nil
	}
	return c.key
}

// Members returns the members of this conversation
func (c *Conversation) Members() map[uuid.UUID]*Gabber {
	c.RLock()
	defer c.RUnlock()
	m := make(map[uuid.UUID]*Gabber, len(c.members))
	for i := range c.members {
		m[i] = c.members[i]
	}
	return m
}

// Expiration returns the expiration
func (c *Conversation) Expiration() time.Time {
	c.RLock()
	defer c.RUnlock()
	return c.expiration
}

// Expired returns true if the conversation is expired and should no longer be used/trusted
func (c *Conversation) Expired() bool {
	c.RLock()
	defer c.RUnlock()
	return c.expiration.Before(time.Now())
}
