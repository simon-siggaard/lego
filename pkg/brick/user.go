package brick

type User struct {
	ID         string
	Username   string
	Location   string
	BrickCount int     `json:"brickCount"`
	Pieces     []Piece `json:"collection"`
}
