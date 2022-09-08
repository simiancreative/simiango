package crud

import (
	sqlMock "github.com/simiancreative/simiango/mocks/data/sql"
	"github.com/simiancreative/simiango/service"
)

func (m *Model) Page(items interface{}, filters Filters, order Order, page int, size int) (*service.ContentResponse, error) {
	ds := m.query(filters, order)
	ds = m.handleAugmentList(ds)

	count, content := m.pageQueries(ds, page, size)

	var total int

	countQuery, countParams, _ := count.ToSQL()
	if err := m.cx.Get(&total, countQuery, countParams...); err != nil {
		return nil, err
	}

	contentQuery, contentParams, _ := content.ToSQL()
	if err := m.cx.Select(items, contentQuery, contentParams...); err != nil {
		return nil, err
	}

	meta := service.ContentResponseMeta{Size: size, Page: page, Total: total}
	resp := service.ToContentResponse(items, meta)

	return &resp, nil
}

func (m *Model) SetupPageTest(
	response interface{},
	filters Filters, order Order, page int, size int,
	responseHandler func(interface{}, string, ...interface{}) error,
) {
	var total int
	cx := &sqlMock.ConnX{}

	m.Initialize("mysql", cx)

	ds := m.query(filters, order)
	ds = m.handleAugmentList(ds)

	count, content := m.pageQueries(ds, page, size)

	countQuery, countParams, _ := count.ToSQL()
	countParams = append([]interface{}{&total, countQuery}, countParams...)

	cx.On("Get", countParams...).
		Return(
			func(v interface{}, s string, s2 ...interface{}) error {
				r, _ := v.(*int)
				*r = 10000
				return nil
			},
		)

	contentQuery, contentParams, _ := content.ToSQL()
	contentParams = append([]interface{}{response, contentQuery}, contentParams...)

	cx.On("Select", contentParams...).Return(responseHandler)

	m.cx = cx
}
