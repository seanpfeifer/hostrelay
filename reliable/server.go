package reliable

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/seanpfeifer/hostrelay/scan"
)

var ErrServerClosed = errors.New("tcp: Server closed")

func ListenAndServeTCP(network, address string) error {
	ln, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	log.Printf(`Listening for %s on "%s"`, ln.Addr().Network(), ln.Addr().String())

	srv := Server{
		doneChan:   make(chan struct{}),
		activeConn: make(map[*playerConn]struct{}),
	}

	return srv.Serve(ln)
}

type Server struct {
	doneChan   chan struct{}
	activeConn map[*playerConn]struct{}
	connMutex  sync.Mutex
}

func (srv *Server) trackConn(c *playerConn, add bool) {
	srv.connMutex.Lock()
	defer srv.connMutex.Unlock()

	if add {
		srv.activeConn[c] = struct{}{}
	} else {
		delete(srv.activeConn, c)
	}
}

func (srv *Server) Close() {
	// Close the server so we don't accept new connections
	select {
	case <-srv.doneChan:
		// Already closed. Don't close again
	default:
		close(srv.doneChan)
	}

	// Cycle through all active connections and close them, then delete from the map
	for c := range srv.activeConn {
		c.Close()
		delete(srv.activeConn, c)
	}
}

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()

	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		conn, err := l.Accept()
		// Handle any errors we get, delaying if there's a temporary error
		if err != nil {
			select {
			case <-srv.doneChan:
				return ErrServerClosed
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				// Retry after the delay
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		// No errors, let's reset our delay then spawn a goroutine to handle the connection
		tempDelay = 0
		player := srv.newConn(conn)
		go srv.handlePlayer(player)
	}
}

func (srv *Server) newConn(conn net.Conn) *playerConn {
	c := &playerConn{
		doneChan: make(chan struct{}),
		outgoing: make(chan []byte),
		conn:     conn,
		srv:      srv,
	}

	srv.trackConn(c, true)

	return c
}

func (srv *Server) broadcast(player *playerConn, msg []byte) {
	// We need to lock our mutex, or we can't cycle through all connections
	srv.connMutex.Lock()
	defer srv.connMutex.Unlock()

	for conn := range srv.activeConn {
		if conn != player {
			// Only send to other players
			conn.send(msg)
		}
	}
}

func (srv *Server) handlePlayer(player *playerConn) {
	defer player.Close()

	// Listen on another goroutine.
	// We can't just do a blocking read in our select, even in "default", as it will hold all broadcasts up for quiet clients.
	go player.listen()

	for {
		select {
		case out := <-player.outgoing:
			_, err := player.conn.Write(out)
			if err != nil {
				log.Println(err)
				return
			}
		case <-player.doneChan:
			return
		}
	}
}

type playerConn struct {
	doneChan chan struct{}
	outgoing chan []byte
	conn     net.Conn
	srv      *Server
}

func (c *playerConn) Read(p []byte) (int, error) {
	return c.conn.Read(p)
}

func (c *playerConn) Close() {
	c.srv.trackConn(c, false)
	c.conn.Close()
}

func (c *playerConn) send(msg []byte) {
	c.outgoing <- msg
}

func (c *playerConn) listen() {
	scan.ListenAndDispatch(c, c.onMessageReceived)

	// If the scanner's done, we're done listening to this connection
	c.doneChan <- struct{}{}
}

func (c *playerConn) onMessageReceived(prefix [scan.PrefixLengthBytes]byte, data []byte) {
	// TODO: Profile and optimize this slice allocation by pre-allocating per connection if necessary.
	out := make([]byte, 0, len(data)+scan.PrefixLengthBytes)
	out = append(out, prefix[:]...)
	out = append(out, data...)
	c.srv.broadcast(c, out)
}
