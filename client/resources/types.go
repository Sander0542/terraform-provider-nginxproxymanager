package resources

import (
	"encoding/json"
)

type boolAsInt bool

type Meta map[string]json.RawMessage

func (b *boolAsInt) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*b = boolAsInt(i != 0)
	return nil
}

func (b boolAsInt) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal(1)
	}
	return json.Marshal(0)
}

func (b boolAsInt) Bool() bool {
	return bool(b)
}

func (m *Meta) Map() map[string]string {
	res := map[string]string{}
	for k, v := range *m {
		res[k] = string(v)
	}
	return res
}
