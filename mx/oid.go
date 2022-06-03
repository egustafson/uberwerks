package mx

import (
	"fmt"
	"strings"
)

// Object Identifier (oid)

type Oid string

const (
	RootOID   Oid = ""
	Separator     = "."
)

func (oid Oid) Child(name string) string {
	return fmt.Sprintf("%s.%s", oid, name)
}

func (oid Oid) Parent() Oid {
	parts := strings.Split(string(oid), Separator)
	plen := len(parts)
	if plen < 2 {
		return RootOID
	}
	return Oid(strings.Join(parts[:plen-1], Separator))
}
