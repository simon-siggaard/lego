package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/simon-siggaard/lego/pkg/brick"
	"github.com/simon-siggaard/lego/pkg/brick/set"
	"github.com/simon-siggaard/lego/pkg/brick/user"
)

func main() {
	service := brick.NewService(
		user.NewCachedStore(user.Store{}),
		set.NewCachedStore(set.Store{}),
	)
	now := time.Now()
	pieces, err := service.FiftyPercent("megabuilder99")
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(now)

	total := 0
	for _, piece := range pieces {
		fmt.Println(piece)
		countString := strings.Split(piece, ":")[2]
		count, err := strconv.Atoi(countString)
		if err != nil {
			panic(err)
		}
		total += count
	}
	fmt.Println("total:", total)

	fmt.Printf("Took %s\n", elapsed)
}
