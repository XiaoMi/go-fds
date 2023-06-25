package cslb

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCSLB(t *testing.T) {
	srv := NewRRDNSService(
		[]string{
			"example.com",
		}, true, true,
	)
	stg := NewRoundRobinStrategy()
	lb := NewLoadBalancer(
		srv,
		stg,
	)

	nodes := srv.Nodes()

	for i := 0; i < 10; i++ {
		next, err := lb.Next()
		log.Println(next)
		assert.Contains(t, nodes, next)
		assert.Nil(t, err)
	}
}

func Test100RCSLB(t *testing.T) {
	var counter uint64 = 0
	srv := NewRRDNSService(
		[]string{
			"example.com",
		}, true, true,
	)
	stg := NewRoundRobinStrategy()
	lb := NewLoadBalancer(
		srv,
		stg,
	)

	ctx, cancel := context.WithCancel(context.Background())
	// 100 concurrent read
	for i := 0; i < 100; i++ {
		go func() {
			done := ctx.Done()
			for {
				select {
				case <-done:
					return
				default:
				}
				_, err := lb.Next()
				assert.Nil(t, err)
				atomic.AddUint64(&counter, 1)
			}
		}()
	}

	time.Sleep(time.Second * 1)
	cancel()

	log.Println("Next() called", counter, "times")
}

func Test100RCSLBRandomFail(t *testing.T) {
	var counter uint64 = 0
	var failedCounter uint64 = 0
	srv := NewRRDNSService(
		[]string{
			"example.com",
		}, true, true,
	)
	stg := NewRoundRobinStrategy()
	lb := NewLoadBalancer(
		srv,
		stg,
	)
	randSlice := make([]bool, 0, 1000)
	for i := 0; i < 1000; i++ {
		randSlice = append(randSlice, rand.Intn(10) < 1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// 100 concurrent read & 10% random fail
	for i := 0; i < 100; i++ {
		go func() {
			done := ctx.Done()
			for {
				select {
				case <-done:
					return
				default:
				}
				n, err := lb.Next()
				assert.Nil(t, err)
				if randSlice[atomic.AddUint64(&counter, 1)%1000] {
					lb.NodeFailed(n)
					atomic.AddUint64(&failedCounter, 1)
				}
			}
		}()
	}

	time.Sleep(time.Second * 1)
	cancel()

	log.Println("Next() called", counter, "times")
	log.Println("NodeFailed() called", failedCounter, "times")
}

func TestFailedRatio(t *testing.T) {
	nodes := []Node{
		&net.TCPAddr{
			IP:   net.IPv4(1, 2, 3, 4),
			Port: 1234,
		},
		&net.TCPAddr{
			IP:   net.IPv4(2, 3, 4, 5),
			Port: 2345,
		},
	}
	srv := NewStaticService(nodes)
	stg := NewRoundRobinStrategy()
	lb := NewLoadBalancer(
		srv,
		stg,
		LoadBalancerOption{
			MaxNodeCount:       NodeCountUnlimited,
			TTL:                TTLUnlimited,
			MaxNodeFailedRatio: 0.2,
			MinSampleSize:      10,
		},
	)
	failOrder := []bool{
		// 0-9
		false, false, false, false, false, false, true, false, false, true,
		// 10-19
		true, false, false, false, false, false, false, false, false, false,
	}
	exiled := []bool{
		// 0-9
		false, false, false, false, false, false, false, false, false, false,
		// 10-19
		false, true, true, true, true, true, true, true, true, true,
	}

	// test refresh for 2 rounds
	for k := 0; k < 2; k++ {
		for i := 0; i < 20; i++ {
			actual := make([]string, 0, 2)
			for j := 0; j < len(nodes); j++ {
				node, err := lb.Next()
				assert.Nil(t, err)
				actual = append(actual, node.String())
			}
			if exiled[i] {
				assert.NotContains(t, actual, nodes[1].String())
			} else {
				assert.Contains(t, actual, nodes[1].String())
			}
			if failOrder[i] {
				lb.NodeFailed(nodes[1])
			}
		}
		<-lb.refresh()
	}
}

func BenchmarkCSLB(b *testing.B) {
	srv := NewRRDNSService(
		[]string{
			"example.com",
		}, true, true,
	)
	stg := NewRoundRobinStrategy()
	lb := NewLoadBalancer(
		srv,
		stg,
	)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		lb.Next()
	}
	b.StopTimer()
}
