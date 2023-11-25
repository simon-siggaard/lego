package main

import (
	"fmt"
	"time"

	"github.com/simon-siggaard/lego/pkg/brick"
)

func main() {
	now := time.Now()
	megabuilder := brick.User{}
	err := brick.GetFromJSON(
		&megabuilder,
		brick.Domain+"/api/user/by-id/d174c807-8880-4f49-866b-6e1ec6527ccf",
	)
	if err != nil {
		panic(err)
	}

	type userCountAndQuantity struct {
		userCount   int
		minQuantity int
	}
	fiftyPercentSet := make(map[string]userCountAndQuantity)
	for _, piece := range megabuilder.Pieces {
		for _, variant := range piece.Variants {
			fiftyPercentSet[piece.ID+":"+string(variant.Color)] = userCountAndQuantity{
				userCount:   0,
				minQuantity: variant.Count,
			}
		}
	}

	users, err := brick.UserCollections()
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
	elapsed := time.Since(now)

	fmt.Println("len(users)", len(users))
	total := 0
	for piece, cnq := range fiftyPercentSet {
		if cnq.userCount < len(users)/2+1 {
			continue
		}
		total += cnq.minQuantity

		fmt.Println(piece, cnq.userCount, cnq.minQuantity)
	}
	fmt.Println("Total", total)
	fmt.Println("Took", elapsed)
}
