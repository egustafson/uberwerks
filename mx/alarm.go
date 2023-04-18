package mx

// Alarm abstraction

type Alarm interface {
	Alarming() bool
}

type Alarmer interface {
	Alarming() bool
	ActiveAlarms() []Alarm
	AllAlarms() []Alarm
}

type BasicAlarm struct {
	OID Oid
}

type AlarmerMixin struct {
	AlarmList []Alarm
}

// static check: *AlarmerMixin isA Alarm
var _ Alarmer = (*AlarmerMixin)(nil)

func (a *AlarmerMixin) Alarming() bool {
	for _, alarm := range a.AlarmList {
		if alarm.Alarming() {
			return true
		}
	}
	return false
}

func (a *AlarmerMixin) ActiveAlarms() []Alarm {
	alarming := make([]Alarm, 0)
	for _, alarm := range a.AlarmList {
		if alarm.Alarming() {
			alarming = append(alarming, alarm)
		}
	}
	return alarming
}

func (a *AlarmerMixin) AllAlarms() []Alarm {
	return a.AlarmList
}
