package main

import (
	"fmt"
	"hash/crc32"
	"log"
	"net"

	"github.com/XiaoMi/go-fds/cslb"
)

func main() {
	lb := cslb.NewLoadBalancer(
		cslb.NewStaticService([]cslb.Node{
			net.IPv4(192, 168, 0, 1),
			net.IPv4(192, 168, 0, 2),
			net.IPv4(192, 168, 0, 3),
		}),
		cslb.NewHashedStrategy(func(key interface{}) (uint64, error) {
			return uint64(crc32.ChecksumIEEE([]byte(key.(string)))), nil
		}),
	)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("request_%d", i)
		if node, err := lb.NextFor(key); err != nil {
			log.Println(err)
		} else {
			ip := node.(net.IP)
			log.Println("using IP " + ip.String())
			// Do something with ip
		}
	}
}
