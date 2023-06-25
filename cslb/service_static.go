package cslb

import (
	"sync/atomic"
	"unsafe"
)

type staticService struct {
	staticNodes []Node
	nodes       unsafe.Pointer // Pointer to []node.Node
}

// NewStaticService represents a simple static list of Node
// Node type: node.Node
func NewStaticService(nodes []Node) *staticService {
	return &staticService{
		staticNodes: nodes,
		nodes:       nil,
	}
}

func (s *staticService) Nodes() []Node {
	nodes := (*[]Node)(atomic.LoadPointer(&s.nodes))
	result := make([]Node, 0, len(*nodes))
	result = append(result, *nodes...)
	return result
}

func (s *staticService) Refresh() {
	atomic.StorePointer(&s.nodes, (unsafe.Pointer)(&s.staticNodes))
}

func (s *staticService) NodeFailedCallbackFunc() func(node Node) {
	return nil
}
