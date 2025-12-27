package config

import (
	"encoding/json"
	"fmt"
)

type IntMap map[string]int

func (i *IntMap) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*i = map[string]int{}
		return nil
	}

	var jsonMap map[string]int
	if err := json.Unmarshal(text, &jsonMap); err == nil {
		*i = jsonMap
		return nil
	}

	return fmt.Errorf("invalid IntMap format: %s", string(text))
}
