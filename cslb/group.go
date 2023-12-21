package cslb

import (
	"math/bits"
	"sync"
	"sync/atomic"
)

const (
	NodeCountUnlimited = (1 << bits.UintSize) / 2 - 1
)

type Group struct {
	m             sync.Map // string => *Node
	originalCount int64
	currentCount  int64

	maxNodeCount int
}

func NewGroup(maxNodeCount int) *Group {
	return &Group{
		m:             sync.Map{},
		originalCount: 0,
		currentCount:  0,
		maxNodeCount:  maxNodeCount,
	}
}

func (g *Group) Set(nodes []Node) {
	nodesCheck := make(map[string]struct{}, len(nodes))
	for _, node := range nodes {
		key := node.String()
		nodesCheck[key] = struct{}{}
		g.m.Store(key, node)
	}

	g.m.Range(func(key, value interface{}) bool {
		if _, ok := nodesCheck[key.(string)]; !ok {
			g.m.Delete(key)
		}
		return true
	})

	atomic.StoreInt64(&g.originalCount, int64(len(nodes)))
	atomic.StoreInt64(&g.currentCount, int64(len(nodes)))
}

func (g *Group) Get() []Node {
	result := make([]Node, 0)
	g.m.Range(func(key, value interface{}) bool {
		result = append(result, value.(Node))
		if len(result) >= g.maxNodeCount {
			return false
		}
		return true
	})
	return result
}

func (g *Group) GetNode(key string) Node {
	if val, loaded := g.m.Load(key); loaded {
		return (val).(Node)
	}
	return nil
}

func (g *Group) GetOriginalCount() int64 {
	return atomic.LoadInt64(&g.originalCount)
}

func (g *Group) GetCurrentCount() int64 {
	return atomic.LoadInt64(&g.currentCount)
}

func (g *Group) Exile(node Node) bool {
	_, loaded := g.m.LoadAndDelete(node.String())
	if loaded {
		atomic.AddInt64(&g.currentCount, -1)
	}
	return loaded
}
