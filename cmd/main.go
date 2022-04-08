package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ratnadeepb/goproxy/internal/incoming"
	"github.com/ratnadeepb/goproxy/internal/loadbalance"
	"github.com/ratnadeepb/goproxy/internal/outgoing"
	irequest "github.com/ratnadeepb/goproxy/internal/request"
)

func main() {
	inProxy := incoming.NewProxy("google.com")

	send := make(chan *irequest.Request, 10)
	outProxy := outgoing.NewRequestHandler(send)
	outMux := mux.NewRouter()
	outMux.PathPrefix("/").HandlerFunc(outProxy.HandleOutgoing)

	inMux := mux.NewRouter()
	inMux.PathPrefix("/").HandlerFunc(inProxy.Handle)

	lb := loadbalance.NewLBThread(send)
	defer lb.Close()
	go func() {
		lb.Run()
	}()

	go func() {
		log.Fatal(http.ListenAndServe(":8080", outMux))
	}()
	log.Fatal(http.ListenAndServe(":8081", inMux))
}
