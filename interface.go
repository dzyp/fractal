package fractal

// Comparator should instruct the tree how to compare two items.
// A negative number should be returned if item1 is less than
// item2, a zero should be returned if they are equal, and a
// positive number should be returned if item1 is greater than
// item2.
type Comparator func(item1, item2 *Item) int

type Item struct {
	Key     []byte
	Payload []byte
}

type Items []*Item

type itemSortWrapper struct {
	c     Comparator
	items Items
}

func (isw *itemSortWrapper) Len() int {
	return len(isw.items)
}

func (isw *itemSortWrapper) Swap(i, j int) {
	isw.items[i], isw.items[j] = isw.items[j], isw.items[i]
}

func (isw *itemSortWrapper) Less(i, j int) bool {
	return isw.c(isw.items[i], isw.items[j]) < 0
}

type common interface {
}

type Writer interface {
	Commit() error
	Save(items ...*Item) error
}

type Reader interface {
	AsMutable() Writer
}

type Persister interface {
	Add(key, payload []byte) error
	Get(key []byte) ([]byte, error)
}
