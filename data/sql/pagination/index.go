package pagination

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/simiancreative/simiango/data/sql"
	m "github.com/simiancreative/simiango/mocks/data/mysql"
	"github.com/simiancreative/simiango/service"
)

type Page struct {
	Cx         sql.ConnX
	Attributes string
	From       string
	Join       string
	Where      string
	Order      string
	Page       int
	PageSize   int
}

func (p *Page) Select(items interface{}, params ...interface{}) (*service.ContentResponse, error) {
	countQuery := p.buildCountQuery()
	query := p.buildQuery()

	var total int

	if err := p.Cx.Get(&total, countQuery, params...); err != nil {
		return nil, err
	}

	if err := p.Cx.Select(items, query, params...); err != nil {
		return nil, err
	}

	meta := service.ContentResponseMeta{Size: p.PageSize, Page: p.Page, Total: total}
	resp := service.ToContentResponse(items, meta)

	return &resp, nil
}

func (p *Page) buildCountQuery() string {
	return buildQuery(`
	SELECT COUNT(*)
	FROM {{ .From }}
	{{if .Join -}} {{.Join}} {{- end}}
	{{if .Where -}} WHERE {{.Where}} {{- end}}
	`, p)
}

func (p *Page) buildQuery() string {
	return buildQuery(`
	SELECT {{.Attributes}}
	FROM {{.From}}
	{{if .Join -}} {{.Join}} {{- end}}
	{{if .Where -}} WHERE {{.Where}} {{- end}}
	{{if .Order -}} ORDER BY {{.Order}} {{- end}}
	LIMIT {{ .Limit }}
	OFFSET {{ .Offset }}
	`, p)
}

func (p *Page) Limit() int {
	return p.PageSize
}

func (p *Page) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func buildQuery(queryTpl string, data interface{}) string {
	var tpl bytes.Buffer

	t := template.New("query")
	t.Parse(queryTpl)

	if err := t.Execute(&tpl, data); err != nil {
		fmt.Printf("Build failed %v", err.Error())
		return ""
	}

	return tpl.String()
}

func (p *Page) SetupTest(
	response interface{},
	responseHandler func(interface{}, string, ...interface{}) error,
	params ...interface{},
) {
	var total int
	cx := &m.ConnX{}

	count := params
	count = append([]interface{}{p.buildCountQuery()}, count...)
	count = append([]interface{}{&total}, count...)

	cx.On("Get", count...).
		Return(
			func(v interface{}, s string, s2 ...interface{}) error {
				r, _ := v.(*int)
				*r = 10000
				return nil
			},
		)

	query := params
	query = append([]interface{}{p.buildQuery()}, query...)
	query = append([]interface{}{response}, query...)

	cx.On("Select", query...).Return(responseHandler)

	p.Cx = cx
}
