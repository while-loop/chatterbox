package main

import (
	"flag"
	"log"
	"net/http"
)

//go:generate go get github.com/elazarl/go-bindata-assetfs/... github.com/jteeuwen/go-bindata/...
//go:generate go-bindata-assetfs www/...

var (
	laddr = flag.String("laddr", ":8080", "local bind address")
)

func main() {

	h := NewHub()
	handler := &http.ServeMux{}
	handler.Handle("/", http.FileServer(assetFS()))
	handler.Handle("/chat", h)
	go h.Start()

	server := http.Server{
		Handler: handler,
		Addr:    *laddr,
	}

	log.Println("Serving http on:", *laddr)
	server.ListenAndServe()
}
