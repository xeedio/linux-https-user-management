package main

import (
	"github.com/donpark/pam"
	humcommon "github.com/xeedio/linux-https-user-management"
)

const etcPasswd = "/etc/passwd"

type PAMHttps struct {
	// your pam vars
}

var ph PAMHttps

func (ph *PAMHttps) Authenticate(hdl pam.Handle, args pam.Args) pam.Value {
	if humcommon.ConfigError {
		humcommon.Log().Info("Authenticate exiting early due to config error")
		return pam.AuthError
	}
	user, err := hdl.GetUser()
	if err != nil {
		return pam.AuthError
	}
	humcommon.Log().Debugf("Got request for user: %v", user)

	userPassword, err := hdl.GetItem(pam.AuthToken)
	if err != nil {
		humcommon.Log().Warnf("Error getting PAM passwd for user: %v", err)
		return pam.AuthError
	}

	humcommon.Log().Debugf("Got password for user %v", user)

	if userPassword == "" {
		humcommon.Log().Info("User password was empty!")
		replies, err := hdl.Conversation(pam.Message{Msg: "Password: ", Style: pam.MessageEchoOff})
		if err != nil {
			humcommon.Log().Warnf("Error getting PAM passwd conversation for user: %v!", err)
			return pam.AuthError
		}
		if len(replies) > 0 {
			userPassword = replies[0]
		}
	}

	if err := hdl.SetItem(pam.AuthToken, userPassword); err != nil {
		humcommon.Log().Warnf("Error setting PAM passwd for user: %v!", err)
		return pam.AuthError
	}

	tokenUser, err := humcommon.Authenticate(user, userPassword)
	if err != nil {
		humcommon.Log().Warnf("Auth error: %v", err)
		return pam.AuthInfoUnavailable
	}

	if tokenUser.Token != "" {
		humcommon.Log().Debugf("Token: %s, User: %+v", tokenUser.Token, tokenUser.User)
		if err := appendLineToFile(tokenUser.User.GetPasswdLine(), etcPasswd); err != nil {
			humcommon.Log().Warnf("Error appending passwd file: %v", err)
			return pam.AuthInfoUnavailable
		}
		if err := tokenUser.WriteTokenFile(); err != nil {
			humcommon.Log().Warnf("Error writing token file: %v", err)
			return pam.AuthInfoUnavailable
		}
		if err := tokenUser.User.WriteUserFile(); err != nil {
			humcommon.Log().Warnf("Error writing user file: %v", err)
			return pam.AuthInfoUnavailable
		}
		return pam.Success
	}

	return pam.PermissionDenied
}

func (ph *PAMHttps) SetCredential(hdl pam.Handle, args pam.Args) pam.Value {
	if humcommon.ConfigError {
		humcommon.Log().Info("SetCredential exiting early due to config error")
		return pam.AuthError
	}
	return pam.Success
}

func init() {
	if err := humcommon.InitTLS(); err != nil {
		humcommon.Log().Warningf("Error init tls: %v", err)
		humcommon.ConfigError = true
	}
	pam.RegisterAuthHandler(&ph)
}

func main() {}
