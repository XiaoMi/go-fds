package cslb

import (
	"errors"
	"sync/atomic"
	"unsafe"
)

type HashFunc func(interface{}) (uint64, error)

type hashedStrategy struct {
	hashFunc HashFunc
	nodes    unsafe.Pointer // pointer to []node.Node
}

func NewHashedStrategy(hashFunc HashFunc) *hashedStrategy {
	return &hashedStrategy{
		hashFunc: hashFunc,
		nodes:    nil,
	}
}

func (s *hashedStrategy) SetNodes(nodes []Node) {
	atomic.StorePointer(&s.nodes, unsafe.Pointer(&nodes))
}

func (s *hashedStrategy) Next() (Node, error) {
	return nil, errors.New("hashed strategy needs an input, use NextFor instead")
}

func (s *hashedStrategy) NextFor(key interface{}) (Node, error) {
	nodes := (*[]Node)(atomic.LoadPointer(&s.nodes))
	if len(*nodes) > 0 {
		hash, err := s.hashFunc(key)
		if err != nil {
			return nil, err
		}
		index := hash % uint64(len(*nodes))
		return (*nodes)[index], nil
	} else {
		return nil, errors.New("empty node list")
	}
}
