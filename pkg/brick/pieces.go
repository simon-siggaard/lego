package brick

type Piece struct {
	ID       string `json:"pieceId"`
	Variants []PieceVariant
}

type Color string

type PieceVariant struct {
	Color Color
	Count int
}
