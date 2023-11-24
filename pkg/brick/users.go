package brick

import (
	"encoding/json"
	"net/http"
)

const domain = "https://d16m5wbro86fg2.cloudfront.net"

type User struct {
	Id         string
	Username   string
	Location   string
	BrickCount int     `json:"brickCount"`
	Pieces     []Piece `json:"collection"`
}

// getFromJson make a GET request to the given url and decodes the response
// into the given struct.
func getFromJson[T any](t *T, url string) error {
	result, err := http.Get(url)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	decoder := json.NewDecoder(result.Body)
	err = decoder.Decode(t)
	if err != nil {
		return err
	}

	return nil
}

func UserCollections() ([]User, error) {
	summaryUrl := domain + "/api/users"
	detailsUrl := domain + "/api/user/by-id"

	users := struct {
		Users []User
	}{}
	err := getFromJson(&users, summaryUrl)
	if err != nil {
		return nil, err
	}

	for n, user := range users.Users {
		user := user // no longer necessary in Go 1.22
		url := detailsUrl + "/" + user.Id

		err = getFromJson(&user, url)
		if err != nil {
			return nil, err
		}

		users.Users[n] = user
	}

	return users.Users, nil
}
