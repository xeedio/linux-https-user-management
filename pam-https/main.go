package main

import (
	"fmt"

	"github.com/donpark/pam"
	humcommon "github.com/xeedio/linux-https-user-management"
)

type mypam struct {
	// your pam vars
}

func (mp *mypam) Authenticate(hdl pam.Handle, args pam.Args) pam.Value {
	humcommon.SetLogPrefix("PAM HTTPS")
	user, err := hdl.GetUser()
	if err != nil {
		return pam.AuthError
	}
	humcommon.LogDebug("AUTH", "Got request for user: ", user)

	userPassword, err := hdl.GetItem(pam.AuthToken)
	if err != nil {
		humcommon.LogFatal("Error getting PAM passwd for user", err)
		return pam.AuthError
	}

	if userPassword == "" {
		humcommon.LogInfo("USER-PASSWORD", "User password was empty!")
		replies, err := hdl.Conversation(pam.Message{Msg: "Password: ", Style: pam.MessageEchoOff})
		if err != nil {
			humcommon.LogFatal("Error getting PAM passwd conversation for user!", err)
			return pam.AuthError
		}
		if len(replies) > 0 {
			userPassword = replies[0]
		}
	}

	if err := hdl.SetItem(pam.AuthToken, userPassword); err != nil {
		humcommon.LogFatal("Error setting PAM passwd for user!", err)
		return pam.AuthError
	}

	tokenUser, err := humcommon.Authenticate(user, userPassword)
	if err != nil {
		humcommon.LogFatal("GET-AUTH", err)
		return pam.AuthInfoUnavailable
	}

	if tokenUser.Token != "" {
		humcommon.LogInfo("DEBUG-AUTH", fmt.Sprintf("Token: %s, User: %+v", tokenUser.Token, tokenUser.User))
		return pam.Success
	}

	return pam.PermissionDenied
}

func (mp *mypam) SetCredential(hdl pam.Handle, args pam.Args) pam.Value {
	fmt.Println("SetCredential:", args)
	return pam.Success
}

var mp mypam

func init() {
	pam.RegisterAuthHandler(&mp)
}

func main() {}
