package peer

// Conn is a Gabby peer connection.  This represents a connection with someone on the Gabby network that we want to communicate with.
// Conn is an io.ReadWriteCloser.  You have to open a Conn before using it.
type Conn struct {
	//
}

// NewConn is a new PeerConn
func NewConn(name string) *Conn {
	return &Conn{}
}

// Write is an io.Writer that sends data to our peer.
func (p *Conn) Write(b []byte) (n int, err error) {
	return 999, nil
}
