package entity

import "encoding/json"

type MyPc struct {
	ModePowerOff string `json:"modepoweroff"`
	TimePowerOff string `json:"timepoweroff"`
}

func (m *MyPc) String() string {
	json, _ := json.Marshal(m)
	return string(json)
}
