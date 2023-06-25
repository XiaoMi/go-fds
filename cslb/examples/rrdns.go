package main

import (
	"log"
	"net"

	"github.com/XiaoMi/go-fds/cslb"
)

func main() {
	lb := cslb.NewLoadBalancer(
		cslb.NewRRDNSService([]string{"example.com"}, true, true),
		cslb.NewRoundRobinStrategy(),
	)

	for i := 0; i < 10; i++ {
		if node, err := lb.Next(); err != nil {
			log.Println(err)
		} else {
			ip := node.(*net.IPAddr)
			log.Println("using IP " + ip.String())
			// Do something with ip
		}
	}
}
