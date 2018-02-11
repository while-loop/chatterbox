package main

import (
	"flag"
	"log"
	"net/http"
)

//go:generate go get github.com/elazarl/go-bindata-assetfs/... github.com/jteeuwen/go-bindata/...
//go:generate go-bindata-assetfs www/...

var (
	// create a flag `-laddr` to bind a given local bind address
	laddr = flag.String("laddr", ":8080", "local bind address")
)

func main() {
	flag.Parse() // parse flags given by the cli

	// create and start a new hub manager
	h := NewHub()
	go h.Start()


	handler := &http.ServeMux{}

	// create an http handler for the Chatterbox website created and handled by bindata
	handler.Handle("/", http.FileServer(assetFS()))
	handler.Handle("/chat", FromHTTP(h)) // handle websocket connections at ws://host/chat

	server := http.Server{
		Handler: handler,
		Addr:    *laddr,
	}

	log.Println("Serving http on:", *laddr)
	server.ListenAndServe() // blocking call to create a server socket and accept new connections
}
