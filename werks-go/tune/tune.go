package tune

// Tunables

import (
	"errors"
	"time"

	"gopkg.in/yaml.v3"
)

type Tunable struct {
	Name    string
	Value   any
	Default any
}

var (
	tunables map[string]Tunable = make(map[string]Tunable)
)

func Init(yml string) error {
	keys := make(map[string]any)
	err := yaml.Unmarshal([]byte(yml), &keys)
	if err != nil {
		return err
	}
	for k, v := range keys {
		tunables[k] = Tunable{
			Name:    k,
			Value:   v,
			Default: v,
		}
	}
	return nil
}

func String(name string, def string) (*string, error) {
	var (
		t  Tunable
		ok bool
	)
	if t, ok = tunables[name]; !ok {
		t = Tunable{
			Name:    name,
			Value:   def,
			Default: def,
		}
		tunables[name] = t
	}
	if v, ok := t.Value.(string); ok {
		return &v, nil
	}
	return nil, errors.New("tunable is not type string")
}

func GetString(name string) (v string, ok bool) {
	if t, ok := tunables[name]; ok {
		v, ok = t.Value.(string)
	}
	return
}

func GetInt(name string) int                { return 0 }
func GetFloat(name string) float64          { return 0.0 }
func GetDuration(name string) time.Duration { return 0 }
