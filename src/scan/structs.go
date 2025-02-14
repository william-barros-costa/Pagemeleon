package scan

import "bytes"

type Trailer struct {
	Root    Object
	Encrypt Object
	Info    Object
	Ids     [][]byte
}

type Object struct {
	Name       string
	Id         uint
	Generation uint
}

func (t Trailer) Equal(other Trailer) bool {
	if t.Root != other.Root {
		return false
	}
	if t.Encrypt != other.Encrypt {
		return false
	}
	if t.Info != other.Info {
		return false
	}
	for i := range t.Ids {
		if !bytes.Equal(t.Ids[i], other.Ids[i]) {
			return false
		}
	}
	return true
}

func (o Object) Equal(other Object) bool {
	if o.Name != other.Name {
		return false
	}
	if o.Id != other.Id {
		return false
	}
	if o.Generation != other.Generation {
		return false
	}
	return true
}
