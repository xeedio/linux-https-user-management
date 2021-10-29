package main

import (
	"bytes"
	"encoding/json"
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
	humcommon.Log().Infof("Got request for user: %v", user)
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
		humcommon.Log().Infof("Token: %s, User: %+v", tokenUser.Token, tokenUser.User)
		if err := appendLineToFile(tokenUser.User.GetPasswdLine(), etcPasswd); err != nil {
			humcommon.Log().Warnf("Error appending passwd file: %v", err)
			return pam.AuthInfoUnavailable
		}
		if err := writeTokenFile(tokenUser.Token); err != nil {
			humcommon.Log().Warnf("Error writing token file: %v", err)
			return pam.AuthInfoUnavailable
		}
		if err := writeUserFile(tokenUser.User); err != nil {
			humcommon.Log().Warnf("Error writing user file: %v", err)
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

func writeUserFile(user humcommon.User) error {
	data, _ := json.MarshalIndent(user, "", " ")
	if _, err := os.Stat(humcommon.AppConfig.UserFile); os.IsNotExist(err) {
		return ioutil.WriteFile(humcommon.AppConfig.UserFile, data, 0644)
	}
	return nil
}

func fileContains(line []byte, filePath string) (bool, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		humcommon.Log().Warnf("Error reading file %s: %v", filePath, err)
		return false, err
	}
	return bytes.Contains(data, line), nil
}

func appendLineToFile(line []byte, filePath string) error {
	present, err := fileContains(line, filePath)
	if err != nil {
		humcommon.Log().Warnf("Error from fileContains: %v", err)
		return err
	}

	if present {
		humcommon.Log().Debug("Line was present in file!")
		return nil
	} else {
		humcommon.Log().Debug("Line was NOT present in file!")
	}

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		humcommon.Log().Warnf("Error opening %s: %v", filePath, err)
		return err
	}
	if _, err := f.Write(line); err != nil {
		f.Close() // ignore error; Write error takes precedence
		humcommon.Log().Warnf("Error writing line %s to file: %v", line, err)
		return err
	}
	if err := f.Close(); err != nil {
		humcommon.Log().Warnf("Error closing file: %v", err)
		return err
	}

	return nil
}

func (mp *mypam) SetCredential(hdl pam.Handle, args pam.Args) pam.Value {
	if humcommon.ConfigError {
		humcommon.Log().Info("Exit early due to config error")
		return pam.AuthError
	}
	humcommon.Log().Debugf("SetCredential args: %v", args)
	return pam.Success
}

var mp mypam

func init() {
	pam.RegisterAuthHandler(&mp)
}

func main() {}
