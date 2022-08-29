package pagination

var total int

type sampleItem interface{}

type Item struct {
	ID string `db:"id"`
}

type Items []Item
