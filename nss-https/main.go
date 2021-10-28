package main

import (
	"fmt"

	. "github.com/protosam/go-libnss"
	. "github.com/protosam/go-libnss/structs"

	humcommon "github.com/xeedio/linux-https-user-management"
)

const httpsRemoteUserID = 1337

var defaultHttpsRemoteUser Passwd

// We're creating a struct that implements LIBNSS stub methods.
type HTTPSRemoteUserImpl struct {
	LIBNSS
}

// PasswdByName() returns a single entry by name.
func (self HTTPSRemoteUserImpl) PasswdByName(name string) (Status, Passwd) {
	humcommon.Log().Infof("PasswordByName: %s", name)
	entry := Passwd{
		Username: name,
		Password: "x",
		UID:      httpsRemoteUserID,
		GID:      100, // users
		Gecos:    "HTTPS Remote User",
		Dir:      fmt.Sprintf("/home/%s", name),
		Shell:    "/bin/bash",
	}
	return StatusSuccess, entry
}

// PasswdByUid() returns a single entry by uid.
func (self HTTPSRemoteUserImpl) PasswdByUid(uid uint) (Status, Passwd) {
	humcommon.Log().Infof("PasswordByUid: %d", uid)
	if uid == httpsRemoteUserID {
		return StatusSuccess, defaultHttpsRemoteUser
	}
	return StatusNotfound, Passwd{}
}

func (self HTTPSRemoteUserImpl) ShadowByName(name string) (Status, Shadow) {
	humcommon.Log().Infof("ShadowByName: %s", name)
	entry := Shadow{
		Username:        name,
		Password:        "!",
		LastChange:      18000,
		MinChange:       0,
		MaxChange:       99999,
		PasswordWarn:    7,
		InactiveLockout: -1,
		ExpirationDate:  -1,
		Reserved:        -1,
	}
	return StatusSuccess, entry
}

func init() {
	defaultHttpsRemoteUser = Passwd{
		Username: "remoteuser",
		Password: "x",
		UID:      httpsRemoteUserID,
		GID:      100, // users
		Gecos:    "HTTPS Remote User",
		Dir:      "/home/remoteuser",
		Shell:    "/bin/bash",
	}

	// We set our implementation to "HTTPSRemoteUserImpl", so that go-libnss will use the methods we create
	SetImpl(HTTPSRemoteUserImpl{})
}

// Placeholder main() stub is neccessary for compile.
func main() {}
