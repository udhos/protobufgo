// Package main implements the tool.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/udhos/protobufgo/proto/addressbook"
	"google.golang.org/protobuf/proto"
)

func main() {

	//
	// create address book
	//
	book := &addressbook.AddressBook{}
	p1 := &addressbook.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*addressbook.Person_PhoneNumber{
			{Number: "555-4321", Type: addressbook.PhoneType_PHONE_TYPE_HOME},
		},
	}
	book.People = append(book.People, p1)

	//
	// protobuf marshal
	//
	out, errProto := proto.Marshal(book)
	if errProto != nil {
		log.Fatalf("protobuf marshal addressbook error: %v", errProto)
	}

	reqBodyReader := bytes.NewBuffer(out)

	req, errReq := http.NewRequestWithContext(context.TODO(), "POST", "http://localhost:8080/service", reqBodyReader)
	if errReq != nil {
		log.Fatalf("request error: %v", errReq)
	}

	req.Header.Add("accept", "text/plain")
	req.Header.Set("content-type", "application/octet-stream")

	httpClient := http.DefaultClient

	resp, errDo := httpClient.Do(req)
	if errDo != nil {
		log.Fatalf("send error: %v", errDo)
	}

	defer resp.Body.Close()

	log.Printf("sent protobuf addressbook: %v", book)

	respBody, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		log.Fatalf("response body error: %v", errBody)
	}

	log.Printf("response status:%d content-type:%v body:\n",
		resp.StatusCode, resp.Header.Values("content-type"))

	str := string(respBody)

	fmt.Print(str)
}
