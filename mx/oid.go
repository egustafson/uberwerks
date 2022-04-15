package mx

// Object Identifier (oid)

type Oid string

const RootOID = ""

func (oid Oid) Child(name string) string {
	return oid + "." + name
}

func (oid Oid) Parent() Oid {
	// TODO
}
