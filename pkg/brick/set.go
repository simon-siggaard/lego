package brick

// Set represents a LEGO set.
type Set struct {
	ID         string
	Name       string
	SetNumber  string
	Pieces     []Piece
	TotalCount int `json:"totalPieces"`
}
