package mx

type HealthState int

const (
	Unknown HealthState = iota
	Critical
	Warning
	Healthy
)

type MRID string

type ManagedResource interface {
	ID() MRID
	Health() HealthState
	Readyz() bool
	Livez() bool
	Alarms() []Alarm
	MXChildren() []ManagedResource
	MXAttach(mr ManagedResource)
}

type BaseManagedResource struct {
	Id       MRID
	Children []ManagedResource
}

func NewBaseManagedResource(id string) ManagedResource {
	return &BaseManagedResource{
		Id: MRID(id),
	}
}

func (mr *BaseManagedResource) ID() MRID {
	return mr.Id
}

func (mr *BaseManagedResource) Health() HealthState {
	return Unknown
}

func (mr *BaseManagedResource) Readyz() bool {
	return false
}

func (mr *BaseManagedResource) Livez() bool {
	return false
}

func (mr *BaseManagedResource) Alarms() []Alarm {
	return nil
}

func (mr *BaseManagedResource) MXChildren() []ManagedResource {
	return mr.Children
}

func (mr *BaseManagedResource) MXAttach(resource ManagedResource) {
	for _, r := range mr.Children {
		if r == resource {
			// resource is already attached
			return
		}
	}
	mr.Children = append(mr.Children, resource)
}
