// Package jsondb holds the common, (client and server) objects for working with
// JSON-DB.
package jsondb

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	IDKey JID = "_id"
)

// JSONObj represents a JSON Object
type JSONObj map[JID]any

// JID is a JSON Map Key.
type JID string

func (jo JSONObj) ID() JID {
	if id, ok := jo[IDKey]; ok {
		if id, ok := id.(string); ok {
			return JID(id)
		} else {
			return JID(fmt.Sprintf("%v", id))
		}
	}
	return ""
}

// Identify extracts or assigns the objects ID and returns it.
func Identify(jo JSONObj) JID {
	jid, ok := jo[IDKey].(string)
	if !ok {
		jid = uuid.NewString()
		jo[IDKey] = jid
	}
	return JID(jid)
}
