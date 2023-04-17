package combinators

import "encoding/json"

func (m Model) SelectFromJSON(dst interface{}, query string, params interface{}) error {
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return m.Cx.Select(dst, query, string(b))
}
