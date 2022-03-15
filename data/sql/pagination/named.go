package pagination

import (
	"github.com/simiancreative/simiango/service"
)

func (p *Page) NamedSelect(items interface{}, arg interface{}) (*service.ContentResponse, error) {
	total, err := getNamedCount(p, arg)
	if err != nil {
		return nil, err
	}

	err = getNamed(p, items, arg)
	if err != nil {
		return nil, err
	}

	meta := service.ContentResponseMeta{Size: p.PageSize, Page: p.Page, Total: total}
	resp := service.ToContentResponse(items, meta)

	return &resp, nil
}

func getNamedCount(p *Page, arg interface{}) (int, error) {
	var total int
	query := p.buildCountQuery()

	nstmt, err := p.Cx.PrepareNamed(query)
	if err != nil {
		return 0, err
	}

	err = nstmt.Get(&total, arg)

	return total, err
}

func getNamed(p *Page, items interface{}, arg interface{}) error {
	query := p.buildQuery()

	nstmt, err := p.Cx.PrepareNamed(query)
	if err != nil {
		return err
	}

	err = nstmt.Select(&items, arg)

	return err
}
