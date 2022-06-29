package tune

// Tunables

import (
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

func GetString(name string) string          { return "" }
func GetInt(name string) int                { return 0 }
func GetFloat(name string) float64          { return 0.0 }
func GetDuration(name string) time.Duration { return 0 }
