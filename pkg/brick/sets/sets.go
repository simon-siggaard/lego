package sets

import (
	"fmt"

	"github.com/simon-siggaard/lego/pkg/brick"
)

type setPiecePart struct {
	ID       string `json:"designId"`
	Material int
	PartType string `json:"partType"`
}

type setPiece struct {
	Part     setPiecePart
	Quantity int
}

type Set struct {
	ID         string
	Name       string
	SetNumber  string
	Pieces     []brick.Piece
	TotalCount int `json:"totalPieces"`
}

type set struct {
	ID         string
	Name       string
	SetNumber  string     `json:"setNumber"`
	SetPieces  []setPiece `json:"pieces"`
	TotalCount int        `json:"totalPieces"`
}

func (s set) AsSet() Set {
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

	return Set{
		ID:         s.ID,
		Name:       s.Name,
		SetNumber:  s.SetNumber,
		Pieces:     pieces,
		TotalCount: s.TotalCount,
	}
}

func setSummaries() ([]Set, error) {
	summaryURL := brick.Domain + "/api/sets"

	s := struct {
		Sets []set
	}{}

	err := brick.GetFromJSON(&s, summaryURL)
	if err != nil {
		return nil, err
	}

	sets := make([]Set, len(s.Sets))
	for n, set := range s.Sets {
		sets[n] = set.AsSet()
	}

	return sets, nil
}

func setDetails(id string) (Set, error) {
	detailsURL := brick.Domain + "/api/set/by-id/" + id
	set := set{}
	err := brick.GetFromJSON(&set, detailsURL)
	if err != nil {
		return Set{}, err
	}

	return set.AsSet(), nil
}

func AvailableSets(userID string) ([]string, error) {
	setsSummaries, err := setSummaries()
	if err != nil {
		return nil, err
	}

	user := brick.User{}
	userDetailsURL := brick.Domain + "/api/user/by-id/" + userID
	err = brick.GetFromJSON(&user, userDetailsURL)
	if err != nil {
		return nil, err
	}

	userPieces := make(map[string]int)
	for _, piece := range user.Pieces {
		for _, variant := range piece.Variants {
			userPieces[piece.ID+string(variant.Color)] += variant.Count
		}
	}

	availableSets := []string{}
	for _, ss := range setsSummaries {
		details, err := setDetails(ss.ID)
		if err != nil {
			continue
		}

		setPieces := make(map[string]int)
		for _, piece := range details.Pieces {
			for _, variant := range piece.Variants {
				setPieces[piece.ID+string(variant.Color)] += variant.Count
			}
		}

		available := true
		for pieceID, setPieceCount := range setPieces {
			if userPieceCount, ok := userPieces[pieceID]; !ok || userPieceCount < setPieceCount {
				available = false
				break
			}
		}

		if available {
			availableSets = append(availableSets, details.Name)
		}
	}

	return availableSets, nil
}
