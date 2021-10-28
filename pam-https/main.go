package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/donpark/pam"
	humcommon "github.com/xeedio/linux-https-user-management"
)

const etcPasswd = "/etc/passwd"

type mypam struct {
	// your pam vars
}

func (mp *mypam) Authenticate(hdl pam.Handle, args pam.Args) pam.Value {
	if humcommon.ConfigError {
		humcommon.Log().Info("Exit early due to config error")
		return pam.AuthError
	}
	user, err := hdl.GetUser()
	if err != nil {
		return pam.AuthError
	}
	humcommon.Log().Debug("AUTH", "Got request for user:", user)

	userPassword, err := hdl.GetItem(pam.AuthToken)
	if err != nil {
		humcommon.Log().Fatal("Error getting PAM passwd for user", err)
		return pam.AuthError
	}

	if userPassword == "" {
		humcommon.Log().Info("USER-PASSWORD", "User password was empty!")
		replies, err := hdl.Conversation(pam.Message{Msg: "Password: ", Style: pam.MessageEchoOff})
		if err != nil {
			humcommon.Log().Fatal("Error getting PAM passwd conversation for user!", err)
			return pam.AuthError
		}
		if len(replies) > 0 {
			userPassword = replies[0]
		}
	}

	if err := hdl.SetItem(pam.AuthToken, userPassword); err != nil {
		humcommon.Log().Fatal("Error setting PAM passwd for user!", err)
		return pam.AuthError
	}

	tokenUser, err := humcommon.Authenticate(user, userPassword)
	if err != nil {
		humcommon.Log().Fatal("GET-AUTH", err)
		return pam.AuthInfoUnavailable
	}

	if tokenUser.Token != "" {
		humcommon.Log().Info("DEBUG-AUTH", fmt.Sprintf("Token: %s, User: %+v", tokenUser.Token, tokenUser.User))
		if err := appendLineToFile(tokenUser.User.GetPasswdLine(), etcPasswd); err != nil {
			humcommon.Log().Fatal("PASSWD-USER", err)
			return pam.AuthInfoUnavailable
		}
		if err := writeTokenFile(tokenUser.Token); err != nil {
			humcommon.Log().Fatal("WRITE-TOKEN", err)
			return pam.AuthInfoUnavailable
		}
		return pam.Success
	}

	return pam.PermissionDenied
}

func writeTokenFile(token string) error {
	if _, err := os.Stat(humcommon.AppConfig.TokenFile); os.IsNotExist(err) {
		return ioutil.WriteFile(humcommon.AppConfig.TokenFile, []byte(token), 0644)
	}
	return nil
}

func fileContains(line []byte, filePath string) (bool, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		humcommon.Log().Fatal("FILE-CONTAINS", err)
		return false, err
	}
	return bytes.Contains(data, line), nil
}

func appendLineToFile(line []byte, filePath string) error {
	present, err := fileContains(line, filePath)
	if err != nil {
		humcommon.Log().Fatal("PASSWD-CONTAINS", err)
		return err
	}

	if present {
		humcommon.Log().Info("APPEND-LINE", "Line was present in file!")
		return nil
	} else {
		humcommon.Log().Debug("APPEND-LINE", "Line was NOT present in file!")
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		humcommon.Log().Fatal("PASSWD-OPEN", err)
		return err
	}
	if _, err := f.Write(line); err != nil {
		f.Close() // ignore error; Write error takes precedence
		humcommon.Log().Fatal("PASSWD-WRITE", err)
		return err
	}
	if err := f.Close(); err != nil {
		humcommon.Log().Fatal("PASSWD-CLOSE", err)
		return err
	}

	return nil
}

func (mp *mypam) SetCredential(hdl pam.Handle, args pam.Args) pam.Value {
	if humcommon.ConfigError {
		humcommon.Log().Info("Exit early due to config error")
		return pam.AuthError
	}
	fmt.Println("SetCredential:", args)
	return pam.Success
}

var mp mypam

func init() {
	pam.RegisterAuthHandler(&mp)
}

func main() {}
