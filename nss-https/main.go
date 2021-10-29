package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	. "github.com/protosam/go-libnss"
	. "github.com/protosam/go-libnss/structs"

	humcommon "github.com/xeedio/linux-https-user-management"
)

var validGroupNames []string
var groupsById map[uint]string
var groupsByName map[string]uint

// We're creating a struct that implements HTTPSRemoteUserImpl stub methods.
type HTTPSRemoteUserImpl struct {
	LIBNSS
}

// PasswdByName() returns a single entry by name.
func (self HTTPSRemoteUserImpl) PasswdByName(name string) (Status, Passwd) {
	if humcommon.ConfigError {
		humcommon.Log().Info("Exit early due to config error")
		return StatusNotfound, Passwd{}
	}

	humcommon.Log().Infof("PasswordByName: %s", name)
	user := &humcommon.User{}
	if err := user.ReadUserFile(); err != nil {
		humcommon.Log().Infof("Can't get user info: %v", err)
		return StatusNotfound, Passwd{}
	}

	entry := Passwd{
		Username: user.Username,
		Password: "x",
		UID:      user.UID,
		GID:      humcommon.GroupID, // users
		Dir:      fmt.Sprintf("/home/%s", user.Username),
		Shell:    "/bin/bash",
	}
	return StatusSuccess, entry
}

// PasswdByUid() returns a single entry by uid.
func (self HTTPSRemoteUserImpl) PasswdByUid(uid uint) (Status, Passwd) {
	if humcommon.ConfigError {
		humcommon.Log().Info("Exit early due to config error")
		return StatusNotfound, Passwd{}
	}

	humcommon.Log().Infof("PasswordByUid: %d", uid)
	user := &humcommon.User{}
	if err := user.ReadUserFile(); err != nil {
		humcommon.Log().Infof("Can't get user info: %v", err)
		return StatusNotfound, Passwd{}
	}

	if uid != user.UID {
		return StatusNotfound, Passwd{}
	}

	entry := Passwd{
		Username: user.Username,
		Password: "x",
		UID:      user.UID,
		GID:      100, // users
		Dir:      fmt.Sprintf("/home/%s", user.Username),
		Shell:    "/bin/bash",
	}
	return StatusSuccess, entry
}

func (self HTTPSRemoteUserImpl) ShadowByName(name string) (Status, Shadow) {
	if humcommon.ConfigError {
		humcommon.Log().Info("Exit early due to config error")
		return StatusNotfound, Shadow{}
	}
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

func (self HTTPSRemoteUserImpl) GroupAll() (Status, []Group) {
	if humcommon.ConfigError {
		humcommon.Log().Info("Exit early due to config error")
		return StatusNotfound, []Group{}
	}

	user := &humcommon.User{}
	if err := user.ReadUserFile(); err != nil {
		humcommon.Log().Infof("Can't get user info: %v", err)
		return StatusNotfound, []Group{}
	}

	groupList := make([]Group, 0)
	for groupName, groupId := range groupsByName {
		groupList = append(groupList, Group{
			Groupname: groupName,
			Password:  "x",
			GID:       groupId,
			Members:   []string{user.Username},
		},
		)
	}

	return StatusSuccess, groupList
}

func (self HTTPSRemoteUserImpl) GroupByName(name string) (Status, Group) {
	if humcommon.ConfigError {
		humcommon.Log().Info("Exit early due to config error")
		return StatusNotfound, Group{}
	}

	user := &humcommon.User{}
	if err := user.ReadUserFile(); err != nil {
		humcommon.Log().Infof("Can't get user info: %v", err)
		return StatusNotfound, Group{}
	}

	if !isValidGroup(name) {
		return StatusNotfound, Group{}
	}

	return StatusSuccess, Group{
		Groupname: name,
		Password:  "x",
		GID:       groupsByName[name],
		Members:   []string{user.Username},
	}
}

func (self HTTPSRemoteUserImpl) GroupByGid(gid uint) (Status, Group) {
	if humcommon.ConfigError {
		humcommon.Log().Info("Exit early due to config error")
		return StatusNotfound, Group{}
	}

	user := &humcommon.User{}
	if err := user.ReadUserFile(); err != nil {
		humcommon.Log().Infof("Can't get user info: %v", err)
		return StatusNotfound, Group{}
	}

	if _, ok := groupsById[gid]; !ok {
		return StatusNotfound, Group{}
	}

	return StatusSuccess, Group{
		Groupname: groupsById[gid],
		Password:  "x",
		GID:       gid,
		Members:   []string{user.Username},
	}
}

func parseEtcGroup() error {
	groupFile, err := os.Open("/etc/group")
	if err != nil {
		humcommon.Log().Warnf("Error opening groups: %v", err)
		return err
	}
	defer groupFile.Close()

	r := csv.NewReader(groupFile)
	r.Comma = ':'
	r.Comment = '#'

	for {
		record, err := r.Read()
		if record == nil || err != nil {
			break
		}
		groupName := record[0]
		if isValidGroup(groupName) {
			groupId, _ := strconv.Atoi(record[2])
			uGroupId := uint(groupId)
			groupsById[uGroupId] = groupName
			groupsByName[groupName] = uGroupId
		}
	}

	return nil
}

func isValidGroup(groupName string) bool {
	for _, testGroup := range validGroupNames {
		if testGroup == groupName {
			return true
		}
	}

	return false
}

func init() {
	validGroupNames = append(validGroupNames, "adm", "cdrom", "sudo", "plugdev", "lpadmin")

	groupsById = make(map[uint]string)
	groupsByName = make(map[string]uint)

	if err := parseEtcGroup(); err != nil {
		humcommon.Log().Warnf("Unable to parse etc group: %v", err)
	}

	// We set our implementation to "HTTPSRemoteUserImpl", so that go-libnss will use the methods we create
	SetImpl(HTTPSRemoteUserImpl{})
}

// Placeholder main() stub is neccessary for compile.
func main() {}
