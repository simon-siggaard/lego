package set

import (
	"fmt"

	"github.com/simon-siggaard/lego/pkg/brick"
)

// Store is a LEGO set store.
type Store struct{}

type setPiecePart struct {
	ID       string `json:"designId"`
	Material int
	PartType string `json:"partType"`
}

type setPiece struct {
	Part     setPiecePart
	Quantity int
}

type rawSet struct {
	ID         string
	Name       string
	SetNumber  string     `json:"setNumber"`
	SetPieces  []setPiece `json:"pieces"`
	TotalCount int        `json:"totalPieces"`
}

func (s rawSet) asSet() brick.Set {
	variants := make(map[string][]brick.PieceVariant)
	for _, setPiece := range s.SetPieces {
		if len(variants[setPiece.Part.ID]) == 0 {
			variants[setPiece.Part.ID] = []brick.PieceVariant{
				{
					Color: brick.Color(fmt.Sprint(setPiece.Part.Material)),
					Count: setPiece.Quantity,
				},
			}
		} else {
			variants[setPiece.Part.ID] = append(variants[setPiece.Part.ID], brick.PieceVariant{
				Color: brick.Color(fmt.Sprint(setPiece.Part.Material)),
				Count: setPiece.Quantity,
			})
		}
	}

	pieces := make([]brick.Piece, len(variants))
	n := 0
	for pieceID, variants := range variants {
		pieces[n] = brick.Piece{
			ID:       pieceID,
			Variants: variants,
		}
		n++
	}

	return brick.Set{
		ID:         s.ID,
		Name:       s.Name,
		SetNumber:  s.SetNumber,
		Pieces:     pieces,
		TotalCount: s.TotalCount,
	}
}

// Summaries returns a list of LEGO set summaries.
func (s Store) Summaries() ([]brick.Set, error) {
	summaryURL := brick.Domain + "/api/sets"

	sts := struct {
		Sets []rawSet
	}{}

	err := brick.GetFromJSON(&sts, summaryURL)
	if err != nil {
		return nil, err
	}

	sets := make([]brick.Set, len(sts.Sets))
	for n, set := range sts.Sets {
		sets[n] = set.asSet()
	}

	return sets, nil
}

// Details returns the details of a LEGO set.
func (s Store) Details(id string) (brick.Set, error) {
	detailsURL := brick.Domain + "/api/set/by-id/" + id
	set := rawSet{}
	err := brick.GetFromJSON(&set, detailsURL)
	if err != nil {
		return brick.Set{}, err
	}

	return set.asSet(), nil
}
