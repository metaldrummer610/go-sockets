package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

const host = "localhost:9000"

func main() {
	ctx := context.Background()
	go server(ctx)

	client1()
	client2()
}

func client1() {
	addr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		panic(err)
	}

	header := make([]byte, 4)
	conn.Read(header)

	size := binary.BigEndian.Uint32(header)
	body := make([]byte, size)

	conn.Read(body)

	p := &person{}
	json.Unmarshal(body, p)

	fmt.Printf("Got %v from server\n", p)

	conn.Close()
}

func client2() {
	addr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP("tcp4", nil, addr)
	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(conn)

	p := &person{}
	json.Unmarshal(body[4:], p)

	fmt.Printf("Got %v from server2\n", p)

	conn.Close()
}

func server(ctx context.Context) {
	addr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			c, err := conn.AcceptTCP()
			if err != nil {
				panic(err)
			}

			p := &person{
				Name: "Robbie",
				Age:  28,
			}

			data, err := json.Marshal(p)

			binary.Write(c, binary.BigEndian, uint32(len(data)))
			c.Write(data)

			c.Close()
		}
	}
}
