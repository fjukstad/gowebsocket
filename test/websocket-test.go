package main

import (
    "github.com/fjukstad/gowebsocket"
    "flag"
)

type Client struct {

}


func main() {
    
    // cmd line flags
    var ip = flag.String("ip", "localhost", "ip to run on")
    var port = flag.String("port", ":3999" ,"port to run on")
    flag.Parse() 
    
    server := gowebsocket.New(*ip, *port)
    server.Start() 




    
}
