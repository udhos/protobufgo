// Package main implements the tool.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/udhos/protobufgo/proto/addressbook"
	"google.golang.org/protobuf/proto"
)

func main() {

	addr := ":8080"
	route := "POST /service"

	mux := http.NewServeMux()

	register(mux, route, handler)

	server := &http.Server{Addr: addr, Handler: mux}

	log.Printf("server listening on %s", addr)
	err := server.ListenAndServe()
	log.Printf("server listening on %s: exited: %v", addr, err)
}

func register(mux *http.ServeMux, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	log.Printf("registering: %s", pattern)
	mux.HandleFunc(pattern, handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Values("accept")
	log.Printf("request from %s: accept:%v", r.RemoteAddr, accept)

	reqBody, errBody := io.ReadAll(r.Body)
	if errBody != nil {
		msg := fmt.Sprintf("body error: %v", errBody)
		log.Print(msg)
		http.Error(w, msg, 500)
		return
	}

	book := &addressbook.AddressBook{}
	if errProto := proto.Unmarshal(reqBody, book); errProto != nil {
		msg := fmt.Sprintf("protobuf unmarshal addressbook error: %v", errProto)
		log.Print(msg)
		http.Error(w, msg, 500)
		return
	}

	log.Printf("received protobuf addressbook: %v", book)

	http.Error(w, fmt.Sprint(book), 200)
}
