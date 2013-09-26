package gowebsocket

import (
    "log"
    "code.google.com/p/go.net/websocket"
    "net"
)    

type Client struct {
    Conn *websocket.Conn
    Listener net.Listener
}

func NewClient(ip, port string) (c *Client) {
    c = new(Client) 
    address := "ws://"+ip+port
    origin := "http://"+ip
    
    var err error
    c.Conn, err = websocket.Dial(address, "ws", origin) 
    
    if err != nil {
        log.Panic("ws dial:", err)
    }

    return c
}

func (c *Client) Send (message string) {
    msg := []byte(message)

    if _, err := c.Conn.Write(msg); err != nil {
        log.Panic("Message", msg," could not be sent", err)
    }
}

func (c *Client) Receive() string {
    
    var incoming = make([]byte, 512)
    var n int
    var err error 

    if n, err = c.Conn.Read(incoming); err != nil {
            log.Fatal(err)
    }
    return string(incoming[:n])

}

