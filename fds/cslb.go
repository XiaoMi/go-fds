package fds

import (
	"context"
	"net"
	"sync"

	"github.com/RangerCD/cslb"
	"github.com/RangerCD/cslb/node"
	"github.com/RangerCD/cslb/service"
	"github.com/RangerCD/cslb/strategy"
)

const (
	DNSCacheTTLSecond = 600
)

type cslbDialer struct {
	net.Dialer
	maxNodeCount uint
	lbs          sync.Map // string => *cslb.LoadBalancer
}

func (d *cslbDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	if host, port, err := net.SplitHostPort(address); err == nil {
		val, ok := d.lbs.Load(host)
		var lb *cslb.LoadBalancer
		if !ok {
			maxNodeCount := d.maxNodeCount
			if maxNodeCount <= 0 {
				maxNodeCount = node.NodeCountUnlimited
			}
			lb = cslb.NewLoadBalancer(
				service.NewRRDNSService([]string{host}, true, true),
				strategy.NewRoundRobinStrategy(),
				cslb.LoadBalancerOption{
					MaxNodeCount:        int(maxNodeCount),
					TTL:                 DNSCacheTTLSecond,
					MinHealthyNodeRatio: cslb.HealthyNodeAny,
				})
			d.lbs.Store(host, lb)
		} else {
			lb = val.(*cslb.LoadBalancer)
		}
		if addr, err := lb.Next(); err == nil {
			return d.Dialer.DialContext(ctx, network, net.JoinHostPort(addr.String(), port))
		}
	}
	// Fall back to default behavior if any error
	return d.Dialer.DialContext(ctx, network, address)
}