package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/donpark/pam"
)

func TestAuthenticate(t *testing.T) {
	args := pam.Args{}
	hdl := pam.Handle{}
	val := mp.Authenticate(hdl, args)
	assert.Equal(t, pam.AuthError, val, "Expected AuthError")
}
