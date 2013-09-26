package main

import (
    "github.com/fjukstad/gowebsocket"
    "log"
    "flag"
)

func main() {
    
    // cmd line flags
    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":3999" ,"port to run on")
    flag.Parse() 
    
    gowebsocket.Start(*ip, *port)

    log.Print("Started websocket man") 
}
