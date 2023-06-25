package cslb

import (
	"net"
	"sync/atomic"
	"unsafe"
)

type rrDNSService struct {
	ipv4      bool
	ipv6      bool
	hostnames []string
	nodes     unsafe.Pointer // Pointer to []*net.IPAddr
}

// NewRRDNSService is for Round-robin DNS load balancing solution.
// Usually multiple A or AAAA records are associated with single hostname.
// Node type: *net.IPAddr
//
// For example:
//   Hostname www.a.com
//     |- A 1.2.3.4
//     |- A 2.3.4.5
//     |- A 3.4.5.6
//     ...
// Everytime a client wants to send a request, one of these A records will be chosen to establish connection.
func NewRRDNSService(hostnames []string, ipv4 bool, ipv6 bool) *rrDNSService {
	return &rrDNSService{
		ipv4:      ipv4,
		ipv6:      ipv6,
		hostnames: hostnames,
		nodes:     nil,
	}
}

func (s *rrDNSService) Nodes() []Node {
	nodes := (*[]*net.IPAddr)(atomic.LoadPointer(&s.nodes))
	result := make([]Node, 0, len(*nodes))
	for _, n := range *nodes {
		result = append(result, n)
	}
	return result
}

func (s *rrDNSService) NodeFailedCallbackFunc() func(node Node) {
	return nil
}

func (s *rrDNSService) Refresh() {
	ips := make([]*net.IPAddr, 0, len(s.hostnames))
	for _, h := range s.hostnames {
		if results, err := net.LookupIP(h); err == nil {
			for _, ip := range results {
				switch {
				case ip.To4() != nil && s.ipv4:
					fallthrough
				case ip.To16() != nil && s.ipv6:
					ips = append(ips, &net.IPAddr{IP: ip})
				}
			}
		}
	}
	atomic.StorePointer(&s.nodes, (unsafe.Pointer)(&ips))
}
