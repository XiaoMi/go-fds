# cslb

Client-Side Load Balancer

**This Project is in early developing state**

## Feature

- [ ] Multiple client-side load balancing solutions support
  - [x] [Round-Robin DNS](https://en.wikipedia.org/wiki/Round-robin_DNS)
  - [ ] [SRV DNS](https://en.wikipedia.org/wiki/SRV_record)
  - [x] Static Node List
- [ ] Multiple distributing strategies
  - [x] Round-Robin
  - [ ] Weighted Round-Robin
  - [x] Hashed
- [x] Exile unhealthy node
- [x] Node list TTL 

## Usage

Example:

```go
package main

import (
	"log"
	
	"github.com/RangerCD/cslb"
)

func main() {
	lb := cslb.NewLoadBalancer(
		cslb.NewRRDNSService([]string{"example.com"}, true, true), 
		cslb.NewRoundRobinStrategy(),
	)

	log.Println(lb.Next()) // IP 1
	log.Println(lb.Next()) // IP 2
}
```
