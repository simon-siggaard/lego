package brick

import (
	"fmt"
)

type SetService interface {
	Summaries() ([]Set, error)
	Details(username string) (Set, error)
}

type UserService interface {
	Summary(username string) (User, error)
	Details(username string) (User, error)
	All() ([]User, error)
}

type Service struct {
	userStore UserService
	setStore  SetService
}

func NewService(userStore UserService, setStore SetService) *Service {
	return &Service{
		userStore: userStore,
		setStore:  setStore,
	}
}

func (s *Service) AvailableSets(username string) ([]string, error) {
	user, err := s.userStore.Summary(username)
	if err != nil {
		return nil, err
	}

	user, err = s.userStore.Details(user.ID)
	if err != nil {
		return nil, err
	}

	userPieces := make(map[string]int)
	for _, piece := range user.Pieces {
		for _, variant := range piece.Variants {
			userPieces[piece.ID+":"+string(variant.Color)] += variant.Count
		}
	}

	setsSummaries, err := s.setStore.Summaries()
	if err != nil {
		return nil, err
	}

	availableSets := []string{}
	for _, ss := range setsSummaries {
		details, err := s.setStore.Details(ss.ID)
		if err != nil {
			return nil, err
		}

		setPieces := make(map[string]int)
		for _, piece := range details.Pieces {
			for _, variant := range piece.Variants {
				setPieces[piece.ID+":"+string(variant.Color)] += variant.Count
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

func (s *Service) FiftyPercent(username string) ([]string, error) {
	user, err := s.userStore.Summary(username)
	if err != nil {
		return nil, err
	}

	user, err = s.userStore.Details(user.ID)
	if err != nil {
		return nil, err
	}

	type userCountAndQuantity struct {
		userCount   int
		minQuantity int
	}
	fiftyPercentSet := make(map[string]userCountAndQuantity)
	for _, piece := range user.Pieces {
		for _, variant := range piece.Variants {
			fiftyPercentSet[piece.ID+":"+string(variant.Color)] = userCountAndQuantity{
				userCount:   0,
				minQuantity: variant.Count,
			}
		}
	}

	users, err := s.userStore.All()
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		for _, piece := range user.Pieces {
			for _, variant := range piece.Variants {
				if _, ok := fiftyPercentSet[piece.ID+":"+string(variant.Color)]; !ok {
					continue
				}
				cnq := fiftyPercentSet[piece.ID+":"+string(variant.Color)]

				cnq.userCount++
				if cnq.minQuantity == 0 || variant.Count < cnq.minQuantity {
					cnq.minQuantity = variant.Count
				}
				fiftyPercentSet[piece.ID+":"+string(variant.Color)] = cnq
			}
		}
	}

	pieces := []string{}
	for piece, cnq := range fiftyPercentSet {
		if cnq.userCount < len(users)/2+1 {
			continue
		}
		pieces = append(pieces, fmt.Sprintf("%s:%d", piece, cnq.minQuantity))
	}

	return pieces, nil
}
