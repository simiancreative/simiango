package natsjs

import (
	"encoding/json"

	"github.com/simiancreative/simiango/errors"
)

func UnmarshalEvent(dest interface{}, str string) (err error) {
	err = json.Unmarshal([]byte(str), &dest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal event")
	}

	return nil
}
