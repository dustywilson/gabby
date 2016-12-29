package gabber

import (
	"github.com/dustywilson/gabby/domain"
	"github.com/dustywilson/gabby/encryption"
)

// Gabber is a Gabby client.  A Gabby server happens to also be a Gabby client, too.
type Gabber struct {
	domain *domain.Domain
	enc    *encryption.Encryption
}

// New returns a new Gabber.
// Domain is required.
// Key can be nil; one will be generated if nil.
func New(d *domain.Domain, k *encryption.Key) *Gabber {
	e := encryption.New(k)
	g := &Gabber{
		domain: d,
		enc:    e,
	}
	return g
}
