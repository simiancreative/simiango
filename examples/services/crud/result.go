package crud

import (
	"fmt"

	"github.com/simiancreative/simiango/service"
)

func result(req service.Req) (interface{}, error) {
	_, exists := req.Params.Get("as_descendants")
	if !exists {
		return nil, fmt.Errorf("as_descendants param is required")
	}

	items := &Products{}
	return c.PageFromReq(items, req)

	// Read, Create, Update, Delete

	// item := &Product{}
	// err = c.One(item, 6)
	// result = append(result, *item)

	// params := &ProductProperties{Name: "Roller Blades"}
	// err = c.Create(params, item)
	// result = append(result, *item)

	// params = &ProductProperties{Name: "Cars"}
	// err = c.Update(item.ID, params, item)
	// result = append(result, *item)

	// err = c.Delete(item.ID)
}
