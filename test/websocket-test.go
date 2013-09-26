package main

import (
    "github.com/fjukstad/gowebsocket"
    "code.google.com/p/go.net/websocket"
    "flag"
    "net"
    "log"
    "fmt"
    "time"
)

type Message struct {
    T string `json:"type"`
    Content string `json:"content`
}

type Client struct {
    Connection *websocket.Conn
    Listener net.Listener
}


func graphConfig(addr, path string) *websocket.Config {
    config, _ := websocket.NewConfig(fmt.Sprintf("ws://%s%s",addr,path), "http://localhost")
    return config
}


func main() {
    
    // cmd line flags
    var ip = flag.String("ip", "127.0.0.1", "ip to run on")
    var port = flag.String("port", ":3999" ,"port to run on")
    flag.Parse() 
    
    server := gowebsocket.New(*ip, *port)
    server.Start() 


    time.Sleep(1 * time.Second) 

    var err error
    
    address := "ws://"+*ip+*port
    origin := "http://"+*ip

    ws, err := websocket.Dial(address, "ws", origin) 
    
    if err != nil {
        log.Panic("ws dial:", err)
    }

    msg := []byte("Hello, websocket man\n")

    log.Print("Trying to send: ", string(msg))
    if _, err := ws.Write(msg); err != nil {
        log.Panic("Message", msg," could not be sent", err)
    }

    log.Print("sent") 

    var incoming = make([]byte, 512)
    var n int
    if n, err = ws.Read(incoming); err != nil {
            log.Fatal(err)
    }
    fmt.Printf("Received: %s.\n", incoming[:n])



    
}
