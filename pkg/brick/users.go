package brick

type User struct {
	ID         string
	Username   string
	Location   string
	BrickCount int     `json:"brickCount"`
	Pieces     []Piece `json:"collection"`
}

func UserCollections() ([]User, error) {
	summaryURL := Domain + "/api/users"
	detailsURL := Domain + "/api/user/by-id"

	users := struct {
		Users []User
	}{}
	err := GetFromJSON(&users, summaryURL)
	if err != nil {
		return nil, err
	}

	for n, user := range users.Users {
		user := user // no longer necessary in Go 1.22
		url := detailsURL + "/" + user.ID

		err = GetFromJSON(&user, url)
		if err != nil {
			return nil, err
		}

		users.Users[n] = user
	}

	return users.Users, nil
}
