package gabby

// Packet is a something
type Packet struct {
	From             *Gabber
	To               *Gabber
	PayloadEncrypted PayloadEncrypted
}
