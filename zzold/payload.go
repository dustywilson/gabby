package gabby

// PayloadEncrypted is an encrypted payload plus its metadata
type PayloadEncrypted struct {
	Signature    Signature
	Conversation Conversation
	Cypher       Cypher
	Payload      []byte
}
