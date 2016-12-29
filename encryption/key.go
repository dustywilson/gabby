package encryption

import (
	"fmt"
	"time"
)

// Key is a key
type Key struct {
	privateKey []byte // !!!: quite sure that []byte isn't going to be the right way to store this.
	publicKey  []byte // !!!: maybe it'd be better to use one of the key built-ins?
	expiration time.Time
	refresh    time.Time
}

// NewKey returns a new Key
func NewKey(keyBytes []byte) *Key {
	k := &Key{
		privateKey: keyBytes,
	}
	if len(k.privateKey) == 0 {
		k.privateKey = []byte{0, 1, 2, 3, 4} // !!!: actually generate a new key instead of this garbage
	}
	return k
}

// NewPublicKey returns a Key for the public key provided
func NewPublicKey(keyBytes []byte, ttl time.Duration) *Key {
	expiration := time.Now().Add(ttl)   // expiration = now + TTL
	refresh := expiration.Add(-ttl / 3) // refresh = 66% of max age
	pubkey := &Key{
		publicKey:  keyBytes,
		expiration: expiration,
		refresh:    refresh,
	}
	return pubkey
}

// Expired returns true if the key is expired and therefore is not usable
// When the key is expired, refresh the Peer via (peer.Peer).GetKey() to get an updated key.
func (k *Key) Expired() bool {
	fmt.Println(k.expiration.String())
	return time.Now().After(k.expiration)
}

// ShouldRefresh returns true if the key is beyond the refresh value
// You can refresh the Peer via (peer.Peer).GetKey() to get an updated key.
func (k *Key) ShouldRefresh() bool {
	fmt.Println(k.refresh.String())
	return time.Now().After(k.refresh)
}

// PublicKey returns the Key's public key value
func (k *Key) PublicKey() []byte {
	return k.publicKey
}
