package natsjs

import (
	"encoding/base64"
	"encoding/json"

	"github.com/simiancreative/simiango/errors"
)

func ArgoEventFromString(str string) (ArgoEvent, error) {
	var e ArgoEvent
	err := json.Unmarshal([]byte(str), &e)
	return e, err
}

func ArgoUnmarshalEvent(dest interface{}, str string) (err error) {
	newEvent, err := ArgoEventFromString(str)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal event")
	}

	str, err = newEvent.JSONBody()
	if err != nil {
		return errors.Wrap(err, "failed to extract JSON body")
	}

	err = json.Unmarshal([]byte(str), &dest)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal event")
	}

	return nil
}

// EventContext represents the context of an event
type ArgoEventContext struct {
	ID              string `json:"id"`
	Source          string `json:"source"`
	SpecVersion     string `json:"specversion"`
	Type            string `json:"type"`
	DataContentType string `json:"datacontenttype"`
	Subject         string `json:"subject"`
	Time            string `json:"time"`
}

// Event represents an event received from an argo event source
type ArgoEvent struct {
	Context ArgoEventContext `json:"context"`
	Data    string           `json:"data"`
}

type ArgoEventData struct {
	Subject string      `json:"subject"`
	Body    interface{} `json:"body"`
}

func (e ArgoEvent) JSONBody() (string, error) {
	var data ArgoEventData

	bytes, err := base64.StdEncoding.DecodeString(e.Data)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return "", err
	}

	bytes, err = json.Marshal(data.Body)

	return string(bytes), err
}
