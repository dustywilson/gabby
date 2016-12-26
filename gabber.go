package gabby

import (
	"sync"

	uuid "github.com/satori/go.uuid"
)

// Gabber is a representation of another Gabby client
type Gabber struct {
	id            uuid.UUID
	keys          map[uuid.UUID]*Key
	conversations map[uuid.UUID]*Conversation
	sync.RWMutex
}

// NewGabber returns a Gabber
func NewGabber(id uuid.UUID, key *Key) *Gabber {
	g := &Gabber{
		id:            id,
		keys:          make(map[uuid.UUID]*Key),
		conversations: make(map[uuid.UUID]*Conversation),
	}
	if key != nil {
		g.keys[key.id] = key
	}
	return g
}

// ID is the Gabber's ID (as far as we know, anyway)
func (g *Gabber) ID() uuid.UUID {
	g.RLock()
	defer g.RUnlock()
	return g.id
}

// Keys returns all of the Gabber's keys that we use in our private one-to-one conversation
// There are potentially multiple keys for the purpose of key rolling
// This strips out expired keys
// I'm hoping we can not provide this and do it under the covers instead.
func (g *Gabber) Keys() map[uuid.UUID]*Key {
	g.RLock()
	defer g.RUnlock()
	k := make(map[uuid.UUID]*Key, len(g.keys))
	for i, v := range g.keys {
		if v.Expired() {
			continue
		}
		k[i] = v
	}
	return k
}

// Conversations returns all the conversations that we have running with this Gabber
// This strips out expired conversations
// Keep in mind that a single remote process/person/whatever could have many Gabber IDs.
// We specifically do not know this and can never know this unless they tell us about their other identities, and even if they did, we probably shouldn't believe them.
func (g *Gabber) Conversations() map[uuid.UUID]*Conversation {
	g.RLock()
	defer g.RUnlock()
	c := make(map[uuid.UUID]*Conversation, len(g.conversations))
	for i, v := range g.conversations {
		if v.Expired() {
			continue
		}
		c[i] = v
	}
	return c
}

// addConversation stuffs a conversation into a Gabber's conversation list
// This will not add an expired conversation
// This does NOT add the conversation to the Gabber
func (g *Gabber) addConversation(c *Conversation) {
	g.Lock()
	defer g.Unlock()
	if c.Expired() {
		return
	}
	g.conversations[c.id] = c
}

// removeConversation removes a conversation from a Gabber's conversation list
// This does NOT remove the conversation from the Gabber
func (g *Gabber) removeConversation(c *Conversation) {
	g.Lock()
	defer g.Unlock()
	delete(g.conversations, c.id)
}
