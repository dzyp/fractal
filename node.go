package fractal

import (
	"sort"

	"github.com/satori/go.uuid"
)

type NodeMap struct {
	Layers map[int][]byte
}

func (n *NodeMap) copy() *NodeMap {
	cp := make(map[int][]byte)
	for layer, key := range n.Layers {
		cp[layer] = key
	}

	return &NodeMap{
		Layers: cp,
	}
}

func newNodeMap() *NodeMap {
	return &NodeMap{
		Layers: make(map[int][]byte, 8),
	}
}

type Node struct {
	Key       []byte
	Items     Items
	MaxLength int
}

func (n *Node) canAdd(number int) bool {
	return len(n.Items)+number <= n.MaxLength
}

func (n *Node) copy() *Node {
	cp := make(Items, len(n.Items))
	copy(cp, n.Items)
	return &Node{
		Key:       uuid.NewV4().Bytes(),
		MaxLength: n.MaxLength,
		Items:     cp,
	}
}

func (n *Node) merge(t *tree, items Items) {
	if len(items) == 0 {
		n.Items = items
		return
	}

	n.Items = append(n.Items, items...)
	isw := &itemSortWrapper{
		c:     t.c.Comparator,
		items: n.Items,
	}
	sort.Sort(isw)
}

func (n *Node) pull() Items {
	oldItems := n.Items
	n.Items = make(Items, 0)
	return oldItems
}

func (n *Node) reset() {
	for i := 0; i < len(n.Items); i++ {
		n.Items[i] = nil
	}

	n.Items = n.Items[:0]
}

func loadNode(t *tree, key []byte) (*Node, error) {
	_, err := t.c.Persister.Get(key)
	if err != nil {
		return nil, err
	}

	n := &Node{}
	return n, nil
}

func newNode(maxLength int) *Node {
	return &Node{
		Key:       uuid.NewV4().Bytes(),
		MaxLength: maxLength,
	}
}
