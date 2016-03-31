package fractal

type context struct {
	nodes map[string]*Node
}

func (c *context) getNode(key []byte) *Node {
	return c.nodes[string(key)]
}

func (c *context) addNode(n *Node) {
	c.nodes[string(n.Key)] = n
}

func newContext() *context {
	return &context{
		nodes: make(map[string]*Node, 10),
	}
}
