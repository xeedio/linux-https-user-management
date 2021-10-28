package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	nss "github.com/protosam/go-libnss"
)

func TestPasswdByName(t *testing.T) {
	status, _ := HTTPSRemoteUserImpl{}.PasswdByName("user1")
	assert.Equal(t, nss.Status(1), status, "Expected success")
}

func TestPasswdByUidBad(t *testing.T) {
	status, _ := HTTPSRemoteUserImpl{}.PasswdByUid(httpsRemoteUserID)
	assert.Equal(t, nss.Status(1), status, "Expected success")
}

func TestPasswdByUidGood(t *testing.T) {
	status, _ := HTTPSRemoteUserImpl{}.PasswdByUid(1234)
	assert.Equal(t, nss.Status(0), status, "Expected NotFound")
}

func TestShadowByName(t *testing.T) {
	status, _ := HTTPSRemoteUserImpl{}.ShadowByName("user1")
	assert.Equal(t, nss.Status(1), status, "Expected success")
}
