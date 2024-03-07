package mx

type Alarm interface {
	Alarming() bool
}

type BaseAlarm struct {
	Owner      MRID
	IsAlarming bool
}

func NewAlarm(owner MRID) Alarm {
	return &BaseAlarm{
		Owner:      MRID(owner),
		IsAlarming: false,
	}
}

func (al BaseAlarm) Alarming() bool {
	return al.IsAlarming
}
