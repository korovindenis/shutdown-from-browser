package memory

import "github.com/korovindenis/shutdown-from-browser/v2/internal/domain/entity"

type Storage struct {
	HostInfo entity.MyPc
}

func New() (*Storage, error) {
	return &Storage{}, nil
}

func (s *Storage) SetPoff(pc entity.MyPc) error {
	s.HostInfo.ModePowerOff = pc.ModePowerOff
	s.HostInfo.TimePowerOff = pc.TimePowerOff
	return nil
}
func (s *Storage) GetTimePoff() (string, error) {
	return s.HostInfo.TimePowerOff, nil
}
func (s *Storage) GetModePoff() (string, error) {
	return s.HostInfo.ModePowerOff, nil
}
