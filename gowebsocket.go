// Implementation inspired by http://gary.beagledreams.com/page/go-websocket-chat.html

package gowebsocket

import (
<<<<<<< HEAD
    "code.google.com/p/go.net/websocket"
    "log"
    "net/http"
    "time"
=======
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
>>>>>>> 7e197148b88cab4129b0e362cb863c252e242f46
)

type hub struct {

	// registered cnnections
	connections map[*connection]bool

	// inbound messages from connections
	broadcast chan string

	// register requests from connections
	register chan *connection

	// unregister requests from connections
	unregister chan *connection
}

var h = hub{
	connections: make(map[*connection]bool),
	broadcast:   make(chan string),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
}

func (h *hub) Run() {
	for {
		select {

		// Register new connection
		case c := <-h.register:
			h.connections[c] = true

		// Unregister connection
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)

		// Broadcast message to all connections. If send buffer
		// is full unregister and close websocket connection
		case m := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- m:
				default:
					delete(h.connections, c)
					close(c.send)
					go c.conn.Close()
				}
			}

		}
	}
}

type connection struct {
	// connection
	conn *websocket.Conn

	// buffered channel of outbund messages
	send chan string
}

func (c *connection) reader(h *hub) {
	for {
		var message [1000]byte
		n, err := c.conn.Read(message[:])
		if err != nil {
			break
		}
		h.broadcast <- string(message[:n])
	}
	c.conn.Close()
}

func (c *connection) writer(h *hub) {
	for message := range c.send {
		err := websocket.Message.Send(c.conn, message)
		if err != nil {
			break
		}
	}
}

func (s *WSServer) connHandler(conn *websocket.Conn) {
	s.Conn = &connection{send: make(chan string, 256), conn: conn}
	s.Hub.register <- s.Conn
	defer func() {
		s.Hub.unregister <- s.Conn
	}()

	go s.Conn.writer(s.Hub)
	s.Conn.reader(s.Hub)
}

type WSServer struct {
	Hub    *hub
	Server *http.Server
	Conn   *connection
}

func New(ip, port string) (s *WSServer) {
	s = new(WSServer)

	s.Hub = &hub{
		connections: make(map[*connection]bool),
		broadcast:   make(chan string),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}

	s.Server = &http.Server{
		Addr:    ip + port,
		Handler: websocket.Handler(s.connHandler),
	}

	return s
}

func (s *WSServer) Start() {
	go s.Hub.Run()

	go func() {
		http.Handle("/ws", s.Server.Handler)
		err := s.Server.ListenAndServe()
		if err != nil {
			log.Panic("Websocket server could not start")
		}
	}()
	log.Print("Websocket server started successfully. Go have fun! ")
}
