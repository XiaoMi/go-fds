package cslb

import (
	"math"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	NodeFailedKey = "node-failed."
	RefreshKey    = "refresh"
)

type LoadBalancer struct {
	service  Service
	strategy Strategy
	option   LoadBalancerOption

	sf       *singleflight.Group
	nodes    *Group
	ttlTimer *time.Timer

	metrics *Metrics
}

func NewLoadBalancer(service Service, strategy Strategy, option ...LoadBalancerOption) *LoadBalancer {
	opt := DefaultLoadBalancerOption
	if len(option) > 0 {
		opt = option[0]
	}

	lb := &LoadBalancer{
		service:  service,
		strategy: strategy,
		option:   opt,
		sf:       new(singleflight.Group),
		nodes:    NewGroup(opt.MaxNodeCount),
		ttlTimer: nil,
		metrics:  NewMetrics(opt.MaxNodeFailedRatio, opt.MinSampleSize),
	}
	<-lb.refresh()

	if lb.option.TTL != TTLUnlimited {
		lb.ttlTimer = time.NewTimer(lb.option.TTL)
	}

	return lb
}

func (lb *LoadBalancer) next(nextFunc func() (Node, error)) (Node, error) {
	next, err := nextFunc()
	if err != nil {
		// Refresh and retry
		<-lb.refresh()
		next, err = nextFunc()
	}

	// Check TTL
	if lb.ttlTimer != nil {
		select {
		case <-lb.ttlTimer.C:
			// Background refresh
			lb.refresh()
		default:
		}
	}

	if lb.metrics != nil {
		lb.metrics.NodeInc(next)
	}

	return next, err
}

func (lb *LoadBalancer) Next() (Node, error) {
	return lb.next(lb.strategy.Next)
}

func (lb *LoadBalancer) NextFor(input interface{}) (Node, error) {
	return lb.next(func() (Node, error) {
		return lb.strategy.NextFor(input)
	})
}

func (lb *LoadBalancer) NodeFailed(node Node) {
	if lb.metrics == nil {
		return
	}

	lb.metrics.NodeFailedInc(node)

	if ratio, err := lb.metrics.GetNodeFailedRatio(node); err == nil && ratio > lb.option.MaxNodeFailedRatio {
		lb.sf.Do(NodeFailedKey+node.String(), func() (interface{}, error) {
			lb.metrics.ResetNode(node)
			lb.nodes.Exile(node)
			if fn := lb.service.NodeFailedCallbackFunc(); fn != nil {
				go fn(node)
			}
			nodes := lb.nodes.Get()
			if len(nodes) <= 0 ||
				math.Round(float64(lb.nodes.GetOriginalCount())*lb.option.MinHealthyNodeRatio) > float64(lb.nodes.GetCurrentCount()) {
				<-lb.refresh()
			} else {
				lb.strategy.SetNodes(nodes)
			}
			return nil, nil
		})
	}
}

func (lb *LoadBalancer) refresh() <-chan singleflight.Result {
	return lb.sf.DoChan(RefreshKey, func() (interface{}, error) {
		lb.service.Refresh()

		if lb.ttlTimer != nil {
			select {
			case <-lb.ttlTimer.C:
			default:
			}
			lb.ttlTimer.Reset(lb.option.TTL)
		}

		if lb.metrics != nil {
			lb.metrics.ResetAllNodes()
		}
		lb.nodes.Set(lb.service.Nodes())
		lb.strategy.SetNodes(lb.nodes.Get())
		return nil, nil
	})
}
