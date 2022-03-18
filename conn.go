package h2conn

import (
	"context"
	"io"
	"net"
	"sync"
)

// Conn is client/server symmetric connection.
// It implements the io.Reader/io.Writer/io.Closer to read/write or close the connection to the other side.
// It also has a Send/Recv function to use channels to communicate with the other side.
type Conn struct {
	ca *conAddr
	r  io.Reader
	wc io.WriteCloser

	cancel context.CancelFunc

	wLock sync.Mutex
	rLock sync.Mutex
}

func newConn(ctx context.Context, ca *conAddr, r io.Reader, wc io.WriteCloser) (*Conn, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Conn{
		ca:     ca,
		r:      r,
		wc:     wc,
		cancel: cancel,
	}, ctx
}

// Write writes data to the connection
func (c *Conn) Write(data []byte) (int, error) {
	c.wLock.Lock()
	defer c.wLock.Unlock()
	return c.wc.Write(data)
}

// Read reads data from the connection
func (c *Conn) Read(data []byte) (int, error) {
	c.rLock.Lock()
	defer c.rLock.Unlock()
	return c.r.Read(data)
}

// LocalAddr is used to get the local address of the
// underlying connection.
func (s *Conn) LocalAddr() net.Addr {
	return s.ca.localAddr
}

// RemoteAddr is used to get the address of remote end
// of the underlying connection
func (s *Conn) RemoteAddr() net.Addr {
	return s.ca.remoteAddr
}

// Close closes the connection
func (c *Conn) Close() error {
	c.cancel()
	return c.wc.Close()
}
