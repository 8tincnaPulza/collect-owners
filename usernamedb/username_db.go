package usernamedb

import (
	"io/ioutil"
	"strings"
)

type UsernameDB struct {
	userMap       map[string]string
	Loader        DBLoader
	AddUnresolved bool
}

func (u *UsernameDB) Load() (err error) {
	u.userMap, err = u.Loader.Load()
	return
}

func (u *UsernameDB) ToUsername(email string) (string, bool) {
	username, ok := u.userMap[email]
	return username, ok
}

func (u *UsernameDB) ToUsernames(emails []string) (usernames []string) {
	for _, email := range emails {
		if username, ok := u.ToUsername(email); ok {
			usernames = append(usernames, "@"+username)
		} else if u.AddUnresolved {
			usernames = append(usernames, email)
		}
	}
	return
}

type DBLoader interface {
	Load() (map[string]string, error)
}

type ContributorsFileDBLoader struct {
	Filename string
}

func (c *ContributorsFileDBLoader) Load() (contributors map[string]string, err error) {
	bytes, err := ioutil.ReadFile(c.Filename)
	if err != nil {
		return
	}
	fileContent := strings.Split(string(bytes), "\n")
	contributors = make(map[string]string)

	for _, line := range fileContent {
		line := strings.Trim(line, " \n")

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, " ")
		contributors[parts[0]] = parts[1]
	}

	return
}

type SimpleDBLoader struct {
	users map[string]string
}

func (s *SimpleDBLoader) AddEntry(email string, username string) {
	if s.users == nil {
		s.users = make(map[string]string)
	}

	s.users[email] = username
}

func (s *SimpleDBLoader) Load() (map[string]string, error) {
	return s.users, nil
}
