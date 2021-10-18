package schemas

import (
	"encoding/json"
	"fmt"
)

type Payload struct {
	Data string `json:"data"`
}

func (s *Payload) String() string {
	return fmt.Sprintf("Data: %s\n", s.Data)
}

func (s *Payload) JSON() ([]byte, error) {
	return json.Marshal(s)
}
