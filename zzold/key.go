package gabby

import (
	"sync"
	"time"

	"github.com/dedis/crypto/abstract"
	uuid "github.com/satori/go.uuid"
)

// Key used for session keys and whatever
// Either gabber or conversation will be non-nil, not both.  If neither, it's no longer a valid key.
type Key struct {
	id           uuid.UUID // ???: where does this come from?
	gabber       *Gabber
	conversation *Conversation
	key          abstract.Point
	expiration   time.Time
	sync.RWMutex
}

// NewKey returns a key
func NewKey(id uuid.UUID, key abstract.Point) *Key {
	k := &Key{
		id:  id,
		key: key,
	}
	return k
}

// Destroy destroys the key by wiping the key, the reference to the Gabber, and immediately expires
func (k *Key) Destroy() {
	k.Lock()
	defer k.Unlock()
	k.expiration = time.Now()
	k.key = nil
	k.gabber = nil
	k.conversation = nil
}

// Expiration returns the expiration
func (k *Key) Expiration() time.Time {
	k.RLock()
	defer k.RUnlock()
	return k.expiration
}

// Expired returns true if the key is expired and should no longer be used/trusted
func (k *Key) Expired() bool {
	k.RLock()
	defer k.RUnlock()
	return k.expiration.Before(time.Now())
}
