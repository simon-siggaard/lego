package brick

import "fmt"

type SetPiecePart struct {
	Id       string `json:"designId"`
	Material int
	PartType string `json:"partType"`
}

type SetPiece struct {
	Part     SetPiecePart
	Quantity int
}

type set struct {
	Id         string
	Name       string
	SetNumber  string     `json:"setNumber"`
	SetPieces  []SetPiece `json:"pieces"`
	TotalCount int        `json:"totalPieces"`
}

func (s set) AsSet() Set {
	variants := make(map[string][]PieceVariant)
	for _, setPiece := range s.SetPieces {
		if len(variants[setPiece.Part.Id]) == 0 {
			variants[setPiece.Part.Id] = []PieceVariant{
				{
					Color: Color(fmt.Sprint(setPiece.Part.Material)),
					Count: setPiece.Quantity,
				},
			}
		} else {
			variants[setPiece.Part.Id] = append(variants[setPiece.Part.Id], PieceVariant{
				Color: Color(fmt.Sprint(setPiece.Part.Material)),
				Count: setPiece.Quantity,
			})
		}
	}

	pieces := make([]Piece, len(variants))
	n := 0
	for pieceId, variants := range variants {
		pieces[n] = Piece{
			Id:       pieceId,
			Variants: variants,
		}
		n += 1
	}

	return Set{
		Id:         s.Id,
		Name:       s.Name,
		SetNumber:  s.SetNumber,
		Pieces:     pieces,
		TotalCount: s.TotalCount,
	}
}

type Set struct {
	Id         string
	Name       string
	SetNumber  string
	Pieces     []Piece
	TotalCount int `json:"totalPieces"`
}

func setSummaries() ([]Set, error) {
	summaryUrl := domain + "/api/sets"

	s := struct {
		Sets []set
	}{}

	err := getFromJson(&s, summaryUrl)
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
	detailsUrl := domain + "/api/set/by-id/" + id
	set := set{}
	err := getFromJson(&set, detailsUrl)
	if err != nil {
		return Set{}, err
	}

	return set.AsSet(), nil
}

func AvailableSets(userId string) ([]string, error) {
	setsSummaries, err := setSummaries()
	if err != nil {
		return nil, err
	}

	user := User{}
	userDetailsUrl := domain + "/api/user/by-id/" + userId
	err = getFromJson(&user, userDetailsUrl)
	if err != nil {
		return nil, err
	}

	userPieces := make(map[string]int)
	for _, piece := range user.Pieces {
		for _, variant := range piece.Variants {
			userPieces[piece.Id+string(variant.Color)] += variant.Count
		}
	}

	availableSets := []string{}
	for _, ss := range setsSummaries {
		details, err := setDetails(ss.Id)
		if err != nil {
			continue
		}

		setPieces := make(map[string]int)
		for _, piece := range details.Pieces {
			for _, variant := range piece.Variants {
				setPieces[piece.Id+string(variant.Color)] += variant.Count
			}
		}

		available := true
		for pieceId, setPieceCount := range setPieces {
			if userPieceCount, ok := userPieces[pieceId]; !ok || userPieceCount < setPieceCount {
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
