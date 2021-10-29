package humcommon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const userID = 1337
const groupID = 100

type User struct {
	Admin    bool   `json:"is_staff"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func (u User) GetPasswdLine() []byte {
	return []byte(fmt.Sprintf("%s:x:%d:%d::/home/%s:/bin/bash\n", u.Username, userID, groupID, u.Username))
}

type TokenUser struct {
	Token string `json:"token"`
	User  User   `json:"user"`
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

	return &tokenUser, nil
}
