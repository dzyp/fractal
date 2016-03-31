package fractal

import (
	"log"

	"github.com/satori/go.uuid"
)

type tree struct {
	c       Config
	Key     []byte   `msg:"k"`
	NodeMap *NodeMap `msg:"nm"`
	ctx     *context
}

func (t *tree) AsMutable() Writer {
	return &tree{
		c:       t.c,
		Key:     uuid.NewV4().Bytes(),
		NodeMap: t.NodeMap.copy(),
		ctx:     newContext(),
	}
}

func (t *tree) Commit() error {
	return nil
}

func (t *tree) save(item *Item) error {
	var err error
	items := Items{item}
	i := uint(0)
	for {
		layer := 1 << i
		log.Printf(`LAYER: %+v, i: %+v`, layer, i)
		key, ok := t.NodeMap.Layers[layer]
		if !ok {
			n := newNode(layer)
			t.ctx.addNode(n)
			t.NodeMap.Layers[layer] = n.Key
			n.merge(t, items)
			break
		}

		n := t.ctx.getNode(key)
		if n == nil {
			n, err = loadNode(t, key)
			if err != nil {
				return err
			}
			n = n.copy()
			t.ctx.addNode(n)
			t.NodeMap.Layers[layer] = n.Key
		}

		log.Printf(`N: %+v`, n)

		log.Printf(`LEN ITEMS: %+v`, len(items)+len(n.Items))
		if !n.canAdd(len(items)) {
			println(`THIS HAPPENED`)
			previousItems := n.pull()
			previousItems = append(previousItems, items...)
			items = previousItems
		} else {
			n.merge(t, items)
			break
		}
		i++
	}

	return nil
}

func (t *tree) Save(items ...*Item) error {
	for _, item := range items {
		err := t.save(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func New(config Config) Reader {
	return &tree{
		c:       config,
		Key:     uuid.NewV4().Bytes(),
		NodeMap: newNodeMap(),
	}
}
