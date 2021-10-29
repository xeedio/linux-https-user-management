package humcommon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const baseUserUID = 2000
const GroupID = 100 // users

type User struct {
	Admin    bool   `json:"is_staff"`
	Email    string `json:"email"`
	ID       uint   `json:"id"`
	UID      uint   `json:"uid"`
	Username string `json:"username"`
}

func (u User) GetPasswdLine() []byte {
	return []byte(fmt.Sprintf("%s:x:%d:%d::/home/%s:/bin/bash\n", u.Username, u.UID, GroupID, u.Username))
}

type TokenUser struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func (t TokenUser) WriteTokenFile() error {
	if _, err := os.Stat(AppConfig.TokenFile); os.IsNotExist(err) {
		return ioutil.WriteFile(AppConfig.TokenFile, []byte(t.Token), 0644)
	}
	return nil
}

func (u User) WriteUserFile() error {
	data, _ := json.MarshalIndent(u, "", "  ")
	if _, err := os.Stat(AppConfig.UserFile); os.IsNotExist(err) {
		return ioutil.WriteFile(AppConfig.UserFile, data, 0644)
	}
	return nil
}

func (u *User) ReadUserFile() error {
	data, err := ioutil.ReadFile(AppConfig.UserFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, u); err != nil {
		return err
	}

	return nil
}

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

func postAuthData(b []byte) (*http.Response, error) {
	buf := bytes.NewBuffer(b)
	client := GetHTTPClient()
	return client.Post(AppConfig.URL, "application/json", buf)
}

func postAuthDataWithRetries(b []byte) (*http.Response, error) {
	var err error
	var resp *http.Response

	for _, backoff := range backoffSchedule {
		resp, err = postAuthData(b)
		if err == nil {
			break
		}

		logger.Infof("Request error: %+v", err)
		logger.Infof("Retrying in %v", backoff)
		time.Sleep(backoff)
	}

	// All retries failed
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Authenticate(user, password string) (*TokenUser, error) {
	logger.Debugf("Making request to %s", AppConfig.URL)

	authStruct := struct {
		User     string `json:"username"`
		Password string `json:"password"`
	}{
		user,
		password,
	}

	b, err := json.Marshal(authStruct)
	if err != nil {
		return nil, err
	}

	resp, err := postAuthDataWithRetries(b)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Debugf("Response: code=%d(%s), length=%d, content-type=%s", resp.StatusCode, resp.Status, resp.ContentLength, resp.Header.Get("Content-Type"))

	if resp.StatusCode != 200 {
		b, _ = ioutil.ReadAll(resp.Body)
		logger.Infof("Invalid response body: %s", b)
		return nil, errors.New(resp.Status)
	}

	tokenUser := TokenUser{}
	err = json.NewDecoder(resp.Body).Decode(&tokenUser)
	if err != nil {
		return nil, err
	}
	tokenUser.User.UID = tokenUser.User.ID + baseUserUID

	return &tokenUser, nil
}
