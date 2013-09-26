package main

import (
    "github.com/fjukstad/gowebsocket"
    "flag"
    "log"
    "time"
    "strconv"
)

// Usage example
func main() {
    
    // cmd line flags
    var ip = flag.String("ip", "127.0.0.1", "ip to run on")
    var port = flag.String("port", ":3999" ,"port to run on")
    flag.Parse() 
    
    // Start Server
    server := gowebsocket.New(*ip, *port)
    server.Start() 


    // Start a single client which sends and receives messages 
    c := gowebsocket.NewClient(*ip, *port) 
    i := 0
    for {
        time.Sleep(1 * time.Second) 
        
        msg := "Hello, Websocket. This is message "+strconv.Itoa(i)

        log.Printf("Sending: %s \n", msg)

        c.Send(msg) 
        recv := c.Receive() 

        log.Printf("Received: %s\n", recv)
        i++
    }


    
}
