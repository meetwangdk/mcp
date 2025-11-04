package sid

import (
	"github.com/sony/sonyflake"
	"time"
)

type Sid struct {
	sf *sonyflake.Sonyflake
}

func NewSid() *Sid {
	settings := sonyflake.Settings{
		StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), // 合法的过去时间
		MachineID: func() (uint16, error) { return 1, nil },
	}
	sf := sonyflake.NewSonyflake(settings)
	if sf == nil {
		panic("sonyflake not created")
	}
	return &Sid{sf}
}
func (s Sid) GenString() (string, error) {
	id, err := s.sf.NextID()
	if err != nil {
		return "", err
	}
	return IntToBase62(int(id)), nil
}
func (s Sid) GenUint64() (uint64, error) {
	return s.sf.NextID()
}
