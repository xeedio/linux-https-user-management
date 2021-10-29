package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	nss "github.com/protosam/go-libnss"
)

func TestPasswdByName(t *testing.T) {
	status, _ := HTTPSRemoteUserImpl{}.PasswdByName("user0")
	assert.Equal(t, nss.Status(0), status, "Expected success")
}

func TestPasswdByUidBad(t *testing.T) {
	status, _ := HTTPSRemoteUserImpl{}.PasswdByUid(2001)
	assert.Equal(t, nss.Status(0), status, "Expected success")
}

func TestPasswdByUidGood(t *testing.T) {
	status, _ := HTTPSRemoteUserImpl{}.PasswdByUid(0234)
	assert.Equal(t, nss.Status(0), status, "Expected NotFound")
}

func TestShadowByName(t *testing.T) {
	status, _ := HTTPSRemoteUserImpl{}.ShadowByName("user0")
	assert.Equal(t, nss.Status(0), status, "Expected success")
}

func TestParseEtcGroup(t *testing.T) {
	if err := parseEtcGroup(); err != nil {
		t.Errorf("Error parsing etc group: %v", err)
	}
	t.Logf("GroupsById: %+v", groupsById)
	t.Logf("GroupsByName: %+v", groupsByName)
}
