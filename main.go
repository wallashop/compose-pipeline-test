package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Test struct {
	buf bytes.Buffer
}

func (t *Test) Write(p []byte) (n int, err error) {
	return t.buf.Write(p)
}
func main() {

	//func (b *Buffer) Write(p []byte) (n int, err error) {
	//(p []byte) (n int, err error)
	host := os.Getenv("host")
	if host == "" {
		host = "localhost"
	}
	var t Test
	//
	conn := struct {
		Host string
	}{host}
	err := template.Must(template.New("name").
		Parse("mongodb://{{.Host}}:27017")).Execute(&t, conn)
	check(err)

	fmt.Println("url:" + t.buf.String())

	client, err :=
		mongo.NewClient(options.Client().ApplyURI(
			t.buf.String()))

	check(err)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	check(err)

	err = client.Ping(ctx, readpref.Primary())

	check(err)
	fmt.Println("connection established")
}
