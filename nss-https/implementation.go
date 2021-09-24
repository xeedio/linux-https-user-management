package main

import (
	. "github.com/protosam/go-libnss"
	. "github.com/protosam/go-libnss/structs"
)

// Placeholder main() stub is neccessary for compile.
func main() {}

func init() {
	// We set our implementation to "TestImpl", so that go-libnss will use the methods we create
	SetImpl(TestImpl{})
}

// We're creating a struct that implements LIBNSS stub methods.
type TestImpl struct {
	LIBNSS
}

// PasswdByName() returns a single entry by name.
func (self TestImpl) PasswdByName(name string) (Status, Passwd) {
	for _, entry := range dbtest_passwd {
		if entry.Username == name {
			return StatusSuccess, entry
		}
	}
	return StatusNotfound, Passwd{}
}

// PasswdByUid() returns a single entry by uid.
func (self TestImpl) PasswdByUid(uid uint) (Status, Passwd) {
	for _, entry := range dbtest_passwd {
		if entry.UID == uid {
			return StatusSuccess, entry
		}
	}
	return StatusNotfound, Passwd{}
}
