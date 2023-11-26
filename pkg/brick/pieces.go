package brick

import (
	"encoding/json"
	"net/http"
)

// Piece represents a LEGO piece.
type Piece struct {
	ID       string `json:"pieceId"`
	Variants []PieceVariant
}

// Color represents a LEGO piece color.
type Color string

// PieceVariant represents a LEGO piece variant.
type PieceVariant struct {
	Color Color
	Count int
}

// Domain is the domain of the LEGO API.
const Domain = "https://d16m5wbro86fg2.cloudfront.net"

// GetFromJSON make a GET request to the given url and decodes the response
// into the given struct.
func GetFromJSON[T any](t *T, url string) error {
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
