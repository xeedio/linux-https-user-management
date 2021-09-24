package main

import (
	. "github.com/protosam/go-libnss/structs"
)

// Test database objects.
var dbtest_passwd []Passwd
var dbtest_group []Group
var dbtest_shadow []Shadow

func init() {
	// Populates the passwd test db.
	dbtest_passwd = append(dbtest_passwd,
		Passwd{
			Username: "testguy1",
			Password: "x",
			UID:      1500,
			GID:      1500,
			Gecos:    "Test user 1",
			Dir:      "/home/testguy1",
			Shell:    "/bin/bash",
		},
		Passwd{
			Username: "testguy2",
			Password: "x",
			UID:      1501,
			GID:      1501,
			Gecos:    "Test user 2",
			Dir:      "/home/testguy2",
			Shell:    "/bin/bash",
		},
	)

	// Populates the group test db.
	dbtest_group = append(dbtest_group,
		Group{
			Groupname: "testguy1",
			Password:  "x",
			GID:       1500,
			Members:   []string{"testguy1"},
		},
		Group{
			Groupname: "testguy2",
			Password:  "x",
			GID:       1501,
			Members:   []string{"testguy2"},
		},
		Group{
			Groupname: "testguyz",
			Password:  "x",
			GID:       1499,
			Members:   []string{"testguy1", "testguy2"},
		},
	)

	// Populates the shadow test db.
	dbtest_shadow = append(dbtest_shadow,
		Shadow{
			Username:        "testguy1",
			Password:        "$6$yZcX.DOY$7bgsJhILMYl3DfMZsYUwoObbVt5Sj9FuujuhVn05Vg9hk.2AXLNy6o1DcPNq0SIyaRZ5YBZer2rYaycuh3qtg1", // Password is "password"
			LastChange:      17920,
			MinChange:       0,
			MaxChange:       99999,
			PasswordWarn:    7,
			InactiveLockout: -1,
			ExpirationDate:  -1,
			Reserved:        -1,
		},
		Shadow{
			Username:        "testguy2",
			Password:        "$6$yZcX.DOY$7bgsJhILMYl3DfMZsYUwoObbVt5Sj9FuujuhVn05Vg9hk.2AXLNy6o1DcPNq0SIyaRZ5YBZer2rYaycuh3qtg1", // Password is "password"
			LastChange:      17920,
			MinChange:       0,
			MaxChange:       99999,
			PasswordWarn:    7,
			InactiveLockout: 0,
			ExpirationDate:  0,
			Reserved:        -1,
		},
	)
}
