package cache

import "encoding/json"

type Marshaler struct {
	D interface{}
}

func (m *Marshaler) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m.D)
}

type Unmarshaler struct {
	D interface{}
}

func (m *Unmarshaler) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &m.D)
}
