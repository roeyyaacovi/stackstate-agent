package topology

import (
	"encoding/json"
	"fmt"
)


type Data = map[string]interface{}

// Component is a representation of a topology component
type Component struct {
	ExternalID string                 `json:"externalId"`
	Type       Type                   `json:"type"`
	Data       Data `json:"data"`
}

// JSONString returns a JSON string of the Component
func (c Component) JSONString() string {
	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("{\"error\": \"%s\"}", err.Error())
	}
	return string(b)
}
