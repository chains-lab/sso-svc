package bodies

import "encoding/json"

type InternalEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
