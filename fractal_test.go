package fractal

import (
	"bytes"
	"encoding/binary"
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type persister struct {
	items map[string][]byte
}

func (p *persister) Add(key, payload []byte) error {
	p.items[string(key)] = payload
	return nil
}

func (p *persister) Get(key []byte) ([]byte, error) {
	return p.items[string(key)], nil
}

func newPersister() *persister {
	return &persister{
		items: make(map[string][]byte, 10),
	}
}

var comparator = func(item1, item2 *Item) int {
	return bytes.Compare(item1.Key, item2.Key)
}

func newItem(value int64) *Item {
	b := make([]byte, 8)
	binary.PutVarint(b, value)
	return &Item{
		Key: b,
	}
}

func defaultConfig() Config {
	return Config{
		Persister:  newPersister(),
		Comparator: comparator,
	}
}

func (t *tree) pprint() {
	println(`PPRINTING`)
	for layer, key := range t.NodeMap.Layers {
		log.Printf(`LAYER: %+v`, layer)
		n := t.ctx.getNode(key)
		log.Printf(
			`N.MAXLENGTH: %+v, N.KEY: %+v, N.LENGTH: %+v`, n.MaxLength, n.Key, len(n.Items),
		)
		ints := make([]int64, 0, len(n.Items))
		for _, item := range n.Items {
			value, _ := binary.Varint(item.Key)
			ints = append(ints, value)
		}
		log.Printf(`ITEMS: %+v`, ints)
	}
}

func (t *tree) verify(tb testing.TB) {
	for layer, key := range t.NodeMap.Layers {
		n := t.ctx.getNode(key)
		assert.Equal(tb, layer, n.MaxLength)
		if len(n.Items) != 0 && len(n.Items) != n.MaxLength {
			tb.Errorf(`wrong number of items for n; layer: %+v, len: %+v`, layer, len(n.Items))
		}

		ints := make([]int, 0, len(n.Items))
		for _, item := range n.Items {
			value, _ := binary.Varint(item.Key)
			ints = append(ints, int(value))
		}

		if !assert.True(tb, sort.IntsAreSorted(ints)) {
			tb.Logf(`ints: %+v`, ints)
		}
	}
}

func TestSimpleAdd(t *testing.T) {
	r := New(defaultConfig())
	w := r.AsMutable()

	err := w.Save(newItem(4))
	require.NoError(t, err)

	err = w.Save(newItem(2))
	require.NoError(t, err)
	w.(*tree).verify(t)

	err = w.Save(newItem(1))
	require.NoError(t, err)

	err = w.Save(newItem(7))
	require.NoError(t, err)
	w.(*tree).verify(t)
}
