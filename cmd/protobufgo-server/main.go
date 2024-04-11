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
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(500)
		fmt.Fprintln(w, msg)
		return
	}

	book := &addressbook.AddressBook{}
	if errProto := proto.Unmarshal(reqBody, book); errProto != nil {
		msg := fmt.Sprintf("protobuf unmarshal addressbook error: %v", errProto)
		log.Print(msg)
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(500)
		fmt.Fprintln(w, msg)
		return
	}

	log.Printf("protobuf addressbook: %v", book)

	// 1/3. response headers

	w.Header().Add("content-type", "application/octet-stream")

	// 2/3. response status

	w.WriteHeader(200)

	// 3/3. response body

	fmt.Fprintln(w, book)
}
