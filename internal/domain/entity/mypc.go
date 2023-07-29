package entity

import "encoding/json"

type MyPc struct {
	ModePowerOff string
	TimePowerOff string
}

func (m *MyPc) String() string {
	json, _ := json.Marshal(m)
	return string(json)
}
