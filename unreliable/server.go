package unreliable

import (
	"log"
	"net"
)

// TODO: This was a first thought on how to deal with client "registration"
//
// I need something where I can check if the messenger is a valid client (has connected to TCP *FIRST*)
// and also something where I can be informed when the client disconnects (should no longer receive broadcasts).
//
// Just the IP probably isn't good enough to tie TCP + UDP connections - may need a special "join" message
// to do that. But this may be good enough for the simple case of playing with my friends.
type ClientInformant interface {
	IsRegisteredClient(*net.UDPAddr) bool
}

func ListenAndServeUDP(network, address string) error {
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP(network, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf(`Listening for %s on "%s"`, conn.LocalAddr().Network(), conn.LocalAddr().String())

	srv := server{
		conn:    conn,
		clients: make(map[cmpUDPAddr]*net.UDPAddr),
	}
	return srv.Serve()
}

type server struct {
	conn    *net.UDPConn
	clients map[cmpUDPAddr]*net.UDPAddr
}

func (srv *server) Serve() error {
	buf := make([]byte, 2048)

	for {
		n, tmpAddr, err := srv.conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}
		msg := buf[:n]
		addr := addrFromUDP(tmpAddr)
		// TODO: Only care about clients that are connected to the ClientInformant
		// if srv.ci.IsRegisteredClient(tmpAddr){
		srv.clients[addr] = tmpAddr
		srv.broadcast(addr, msg)
	}
}

func (srv *server) broadcast(player cmpUDPAddr, msg []byte) {
	for cmpAddr, sendPlr := range srv.clients {
		if cmpAddr != player {
			srv.conn.WriteToUDP(msg, sendPlr)
		}
	}
}

// cmpUDPAddr is a UDPAddr that can be compared and converted back to UDPAddr.
// This is used specifically because "UDPAddr.IP" is a []byte, which isn't comparable.
type cmpUDPAddr struct {
	IP   string // Convert to net.IP by casting to []byte
	Port int
	Zone string
}

func (a *cmpUDPAddr) ToUDPAddr() net.UDPAddr {
	return net.UDPAddr{
		IP:   []byte(a.IP),
		Port: a.Port,
		Zone: a.Zone,
	}
}

func addrFromUDP(udpAddr *net.UDPAddr) cmpUDPAddr {
	return cmpUDPAddr{
		IP:   string(udpAddr.IP),
		Port: udpAddr.Port,
		Zone: udpAddr.Zone,
	}
}
