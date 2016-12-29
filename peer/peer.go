package peer

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/dustywilson/gabby/dns"
	"github.com/dustywilson/gabby/domain"
	"github.com/dustywilson/gabby/encryption"
)

// Peer is a Gabby peer
type Peer struct {
	id     string
	domain *domain.Domain
	key    *encryption.Key
}

// ErrNameInvalid is an error returned when the provided peer name is invalid
var ErrNameInvalid = errors.New("invalid peer name format")

// String returns the string representation of the Peer
func (p *Peer) String() string {
	return p.id + "@" + p.domain.String()
}

// Parse returns a Peer from the string representation of a peer
func Parse(name string) (*Peer, error) {
	parts := strings.Split(name, "@")
	if len(parts) != 2 {
		return nil, ErrNameInvalid
	}

	id := parts[0]
	if strings.Contains(id, ".") {
		return nil, ErrNameInvalid
	}

	d, err := domain.New(parts[1])
	if err != nil {
		return nil, err
	}

	return &Peer{
		id:     id,
		domain: d,
	}, nil
}

// Key returns the Peer's Key
func (p *Peer) Key() *encryption.Key {
	return p.key
}

// FetchKey retrieves the peer's public key.  This returns the key to the caller; more importantly it updates the key in the Peer struct.
func (p *Peer) FetchKey() (*encryption.Key, error) {
	txts, err := dns.LookupTXT(strings.Join([]string{p.id, "_keys._gabby", p.domain.String()}, ".")) // "$id._keys._gabby.$domain" TXT record
	if err != nil {
		return nil, err
	}
	for _, txt := range txts { // note: it's best if there's only one key listed; remembering old keys for key rolling is a problem for the key's owner, not the client
		keyB64 := txt.Value
		keyBytes, err := base64.StdEncoding.DecodeString(keyB64)
		if err != nil {
			continue
		}
		p.key = encryption.NewPublicKey(keyBytes, txt.TTL)
		return p.key, nil // we only want a single key
	}
	return nil, errors.New("no public key found for this peer")
}
