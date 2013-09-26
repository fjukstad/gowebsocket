
// Implementation inspired by http://gary.beagledreams.com/page/go-websocket-chat.html

package gowebsocket

import (
    "code.google.com/p/go.net/websocket"
    "log"
    "net/http"
)

type hub struct {
    
    // registered cnnections
    connections map[*connection] bool

    // inbound messages from connections
    broadcast chan string

    // register requests from connections
    register chan *connection

    // unregister requests from connections
    unregister chan *connection
    
}

var h = hub {
    connections:    make(map[*connection] bool),
    broadcast:      make(chan string),
    register:       make(chan *connection),
    unregister:     make(chan *connection),
}


func (h *hub) run(){
    for {
        select {
            // Register new connection
            case c:= <-h.register:
                h.connections[c] = true

            // Unregister connection
            case c:= <-h.unregister:
                delete(h.connections, c)
                close(c.send)

            // Broadcast message to all connections. If send buffer
            // is full unregister and close websocket connection 
            case m:= <-h.broadcast:
                for c:= range h.connections {
                    log.Println("sending message to", c.conn.RemoteAddr())
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

func (c *connection) reader() {
    for {
        var message [1000]byte
        n, err := c.conn.Read(message[:])
        if err != nil{
            break
        }
        h.broadcast <- string(message[:n])
    }
    c.conn.Close()
}


func (c *connection) writer() {
    for message := range c.send{
        err := websocket.Message.Send(c.conn, message)
        if err != nil {
            break
        }
    }
}

func connHandler(conn *websocket.Conn){
    c := &connection{send: make(chan string, 256), conn: conn}
    h.register <- c
    defer func() {
        h.unregister <- c
    }() 
    
    go c.writer()
    c.reader()
}


func Start(ip, port string) {
    address := ip+port
    go h.run() 
    http.Handle("/", websocket.Handler(connHandler))
    if err := http.ListenAndServe(address, nil); err != nil{
        log.Panic("gowebsocket could not start: ", err)
    }
}

/*
func main() {
    
    // cmd line flags
    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":3999" ,"port to run on")
    flag.Parse() 

    address := *ip + *port


    log.Println("Websocket broadcaster started on", address)

    go h.run()
    http.Handle("/", websocket.Handler(connHandler))

    if err := http.ListenAndServe(address, nil); err != nil {
        log.Panic("ListenAndServe:", err)
    }
}
*/

