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

// Structs populated on init
var validGroupNames []string
var groupsById map[uint]string
var groupsByName map[string]uint

// Singleton user
var user *humcommon.User

// We're creating a struct that implements HTTPSRemoteUserImpl stub methods.
type HTTPSRemoteUserImpl struct {
	LIBNSS
}

func loadUser() error {
	if user != nil {
		return nil
	}

	tmpUser := &humcommon.User{}
	if err := tmpUser.ReadUserFile(); err != nil {
		return err
	}

	// Set global user to temp user
	user = tmpUser

	return nil
}

// PasswdByName() returns a single entry by name.
func (self HTTPSRemoteUserImpl) PasswdByName(name string) (Status, Passwd) {
	if humcommon.ConfigError {
		humcommon.Log().Debug("Exit early due to config error")
		return StatusNotfound, Passwd{}
	}

	humcommon.Log().Debugf("PasswordByName Start: %s", name)
	if err := loadUser(); err != nil {
		humcommon.Log().Debugf("PasswdByName Can't get user info: %v", err)
		return StatusNotfound, Passwd{}
	}

	if name != user.Username {
		humcommon.Log().Debugf("PasswdByName Wrong user: %s", name)
		return StatusNotfound, Passwd{}
	}

	humcommon.Log().Debugf("PasswordByName Success: %s", name)
	return StatusSuccess, Passwd{
		Username: user.Username,
		Password: "x",
		UID:      user.UID,
		GID:      humcommon.GroupID, // users
		Dir:      fmt.Sprintf("/home/%s", user.Username),
		Shell:    "/bin/bash",
	}
}

// PasswdByUid() returns a single entry by uid.
func (self HTTPSRemoteUserImpl) PasswdByUid(uid uint) (Status, Passwd) {
	if humcommon.ConfigError {
		humcommon.Log().Debug("PasswdByUid Exit early due to config error")
		return StatusNotfound, Passwd{}
	}

	humcommon.Log().Debugf("PasswdByUid Start: %d", uid)
	if err := loadUser(); err != nil {
		humcommon.Log().Debugf("PasswdByUid Can't get user info: %v", err)
		return StatusNotfound, Passwd{}
	}

	if uid != user.UID {
		humcommon.Log().Debugf("PasswdByUid Wrong uid: %d", uid)
		return StatusNotfound, Passwd{}
	}

	humcommon.Log().Debugf("PasswordByUid Success: %d", uid)
	return StatusSuccess, Passwd{
		Username: user.Username,
		Password: "x",
		UID:      user.UID,
		GID:      100, // users
		Dir:      fmt.Sprintf("/home/%s", user.Username),
		Shell:    "/bin/bash",
	}
}

func (self HTTPSRemoteUserImpl) ShadowByName(name string) (Status, Shadow) {
	if humcommon.ConfigError {
		humcommon.Log().Debug("Exit early due to config error")
		return StatusNotfound, Shadow{}
	}
	humcommon.Log().Debugf("ShadowByName Start: %s", name)

	if err := loadUser(); err != nil {
		humcommon.Log().Debugf("ShadowByName Can't get user info: %v", err)
		return StatusNotfound, Shadow{}
	}

	if name != user.Username {
		humcommon.Log().Debugf("ShadowByName Wrong user: %s", name)
		return StatusNotfound, Shadow{}
	}

	humcommon.Log().Debugf("ShadowByName Success: %s", name)
	return StatusSuccess, Shadow{
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
}

func (self HTTPSRemoteUserImpl) GroupAll() (Status, []Group) {
	if humcommon.ConfigError {
		humcommon.Log().Debug("Exit early due to config error")
		return StatusNotfound, []Group{}
	}

	if err := loadUser(); err != nil {
		humcommon.Log().Debugf("Can't get user info: %v", err)
		return StatusNotfound, []Group{}
	}

	humcommon.Log().Debugf("GroupAll Start")

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

	humcommon.Log().Debugf("GroupAll Success")

	return StatusSuccess, groupList
}

func (self HTTPSRemoteUserImpl) GroupByName(name string) (Status, Group) {
	if humcommon.ConfigError {
		humcommon.Log().Debug("Exit early due to config error")
		return StatusNotfound, Group{}
	}

	humcommon.Log().Debugf("GroupByName Start: %s", name)

	if err := loadUser(); err != nil {
		humcommon.Log().Debugf("GroupByName Can't get user info: %v", err)
		return StatusNotfound, Group{}
	}

	if !isValidGroup(name) {
		humcommon.Log().Debugf("GroupByName Invalid: %s", name)
		return StatusNotfound, Group{}
	}

	humcommon.Log().Debugf("GroupByName Success: %s", name)

	return StatusSuccess, Group{
		Groupname: name,
		Password:  "x",
		GID:       groupsByName[name],
		Members:   []string{user.Username},
	}
}

func (self HTTPSRemoteUserImpl) GroupByGid(gid uint) (Status, Group) {
	if humcommon.ConfigError {
		humcommon.Log().Debug("Exit early due to config error")
		return StatusNotfound, Group{}
	}

	if err := loadUser(); err != nil {
		humcommon.Log().Debugf("Can't get user info: %v", err)
		return StatusNotfound, Group{}
	}

	humcommon.Log().Debugf("GroupByGid: %d", gid)

	if _, ok := groupsById[gid]; !ok {
		humcommon.Log().Debugf("GroupByGid Invalid: %d", gid)
		return StatusNotfound, Group{}
	}

	humcommon.Log().Debugf("GroupByGid Success: %d", gid)

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
	validGroupNames = append(validGroupNames, "adm", "cdrom", "sudo", "plugdev", "lpadmin", "users")

	groupsById = make(map[uint]string)
	groupsByName = make(map[string]uint)

	if err := parseEtcGroup(); err != nil {
		humcommon.Log().Warnf("Unable to parse etc group: %v", err)
	}

	if err := loadUser(); err != nil {
		humcommon.Log().Debugf("Init can't yet load user json: %v", err)
	}

	// We set our implementation to "HTTPSRemoteUserImpl", so that go-libnss will use the methods we create
	SetImpl(HTTPSRemoteUserImpl{})
}

// Placeholder main() stub is neccessary for compile.
func main() {}
