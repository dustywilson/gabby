package encryption

// Encryption is a poorly named thing.
type Encryption struct {
	key *Key
}

// New returns a thing.
func New(key *Key) *Encryption {
	return &Encryption{
		key: key,
	}
}
