package main

import (
	"log"

	"github.com/XiaoMi/go-fds/cslb"
)

type DB interface{}

// Let's pretend DB is the client type in some kind of database SDK

type Client struct {
	id string
	db *DB
}

func NewClient(addr string) *Client {
	return &Client{
		id: addr,
		db: nil, // Connect to your database node here
	}
}

func (client Client) String() string {
	return client.id
}

func (client Client) DB() *DB {
	return client.db
}

func main() {
	clients := []cslb.Node{
		NewClient("192.168.1.100"),
		NewClient("10.0.0.100"),
	}

	lb := cslb.NewLoadBalancer(
		cslb.NewStaticService(clients),
		cslb.NewRoundRobinStrategy(),
	)

	for i := 0; i < 10; i++ {
		if node, err := lb.Next(); err != nil {
			log.Println(err)
		} else {
			log.Println("using node " + node.(*Client).String())
			// db := node.(*Client).DB()
			// Query something with db
		}
	}
}
