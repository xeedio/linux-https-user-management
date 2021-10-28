package humcommon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	buf := bytes.NewBuffer(b)

	client := GetHTTPClient()
	resp, err := client.Post(AppConfig.URL, "application/json", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Debugf("Response: code=%d(%s),length=%d,content-type=%s", resp.StatusCode, resp.Status, resp.ContentLength, resp.Header.Get("Content-Type"))

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
