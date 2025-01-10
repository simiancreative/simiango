package nats

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/simiancreative/simiango/messaging/natsjs"
)

func BuildArgoEvent(workflowName, action, msg string) string {
	subject := fmt.Sprintf("%s-%s", workflowName, action)

	data := natsjs.ArgoEventData{
		Subject: subject,
	}

	json.Unmarshal([]byte(msg), &data.Body)
	dataStr, _ := json.Marshal(data)

	event := natsjs.ArgoEvent{
		Context: natsjs.ArgoEventContext{
			ID:              fmt.Sprintf("%x", md5.Sum([]byte(subject))),
			Source:          "nats-event-source",
			SpecVersion:     "1.0",
			Type:            "nats",
			DataContentType: "application/json",
			Subject:         subject,
			Time:            time.Now().Format(time.RFC3339),
		},
		Data: base64.StdEncoding.EncodeToString(dataStr),
	}

	jsonData, _ := json.Marshal(event)

	return string(jsonData)
}
