package user

import "github.com/simon-siggaard/lego/pkg/brick"

// Store is a LEGO user store.
type Store struct{}

// Summary returns a LEGO user summary.
func (s Store) Summary(username string) (brick.User, error) {
	user := brick.User{}
	userSummaryURL := brick.Domain + "/api/user/by-username/" + username
	err := brick.GetFromJSON(&user, userSummaryURL)
	if err != nil {
		return brick.User{}, err
	}

	return user, nil
}

// Details returns a LEGO user details.
func (s Store) Details(id string) (brick.User, error) {
	user := brick.User{}
	userDetailsURL := brick.Domain + "/api/user/by-id/" + id
	err := brick.GetFromJSON(&user, userDetailsURL)
	if err != nil {
		return brick.User{}, err
	}

	return user, nil
}

// All returns details on all LEGO users.
func (s Store) All() ([]brick.User, error) {
	summaryURL := brick.Domain + "/api/users"
	detailsURL := brick.Domain + "/api/user/by-id"

	users := struct {
		Users []brick.User
	}{}
	err := brick.GetFromJSON(&users, summaryURL)
	if err != nil {
		return nil, err
	}

	for n, user := range users.Users {
		user := user // no longer necessary in Go 1.22
		url := detailsURL + "/" + user.ID

		err = brick.GetFromJSON(&user, url)
		if err != nil {
			return nil, err
		}

		users.Users[n] = user
	}

	return users.Users, nil
}
