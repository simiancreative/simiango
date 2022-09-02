package main

import (
	"fmt"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/simiancreative/simiango/data/mysql"
	"github.com/simiancreative/simiango/data/sql"
	"github.com/simiancreative/simiango/service"

	"github.com/doug-martin/goqu/v9"
)

type Item struct {
}
type Items []Item

type Kind int

const (
	AND Kind = iota
	OR
)

var c = CRUD{
	Table: "Devices",
	Filterable: Filterable{
		"type": func(value interface{}) (Kind, interface{}) {
			return AND, Filters{"type": value}
		},
	},
}

func main() {
	c.Initialize(mysql.Cx, "mysql")

	items := &Items{}

	result, err := c.Page(items, Filters{}, Orders{}, 0, 0)
	fmt.Println(result)
	fmt.Println(err)

	item := &Item{}
	err = c.One(item, Filters{"type": "gateway"}, Orders{"type": "asc"}, 23)
	fmt.Println(item)
	fmt.Println(err)
}

type FilterableFunc func(interface{}) (Kind, interface{})
type Filterable map[string]FilterableFunc
type AugmentListQuery *func(*goqu.SelectDataset) *goqu.SelectDataset

type CRUD struct {
	cx               sql.ConnX
	dialect          goqu.DialectWrapper
	MultiTenant      bool
	Paginate         bool
	Table            string
	Filterable       Filterable
	AugmentListQuery AugmentListQuery
	AugmentOneQuery  AugmentListQuery
}

type Filters map[string]interface{}
type Orders map[string]interface{}
type Size int
type Page int

func (c *CRUD) Initialize(cx sql.ConnX, dialect string) {
	c.cx = cx
	c.dialect = goqu.Dialect(dialect)
}

func (c *CRUD) One(items interface{}, filters Filters, orders Orders, id interface{}) error {
	ds := c.query(filters, orders)
	ds = ds.Where(goqu.Ex{"id": id})

	query, params, _ := ds.ToSQL()
	if err := c.cx.Get(items, query, params...); err != nil {
		return err
	}

	return nil
}

func (c *CRUD) Page(items interface{}, filters Filters, orders Orders, page int, size int) (*service.ContentResponse, error) {
	ds := c.query(filters, orders)
	count, content := c.pageQueries(ds, page, size)

	var total int

	countQuery, countParams, _ := count.ToSQL()
	if err := c.cx.Get(&total, countQuery, countParams...); err != nil {
		return nil, err
	}

	contentQuery, contentParams, _ := content.ToSQL()
	if err := c.cx.Select(items, contentQuery, contentParams...); err != nil {
		return nil, err
	}

	meta := service.ContentResponseMeta{Size: size, Page: page, Total: total}
	resp := service.ToContentResponse(items, meta)

	return &resp, nil
}

func (c *CRUD) query(filters Filters, orders Orders) *goqu.SelectDataset {
	ds := c.dialect.From(c.Table)

	for key, value := range filters {
		filter, ok := c.Filterable[key]
		if !ok {
			continue
		}

		kind, expressions := filter(value)
		if kind == AND {
			ds = ds.Where(goqu.Ex(expressions.(Filters)))
		}

		if kind == OR {
			ds = ds.Where(goqu.Ex(expressions.(Filters)))
		}
	}

	for key, value := range orders {
		ordr := goqu.C(key)

		if value == "asc" {
			ds = ds.Order(ordr.Asc())
		}

		if value == "dsc" {
			ds = ds.Order(ordr.Desc())
		}
	}

	if c.AugmentListQuery != nil {
		cb := *c.AugmentListQuery
		ds = cb(ds)
	}

	return ds
}

func (c *CRUD) pageQueries(ds *goqu.SelectDataset, size int, page int) (*goqu.SelectDataset, *goqu.SelectDataset) {
	countQuery := ds.Select(goqu.COUNT("*"))
	ds = ds.Limit(uint(size)).Offset(uint((page - 1) * size))
	return countQuery, ds
}
