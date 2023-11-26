package main

import (
	"fmt"
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
	sets, err := service.AvailableSets("brickfan35")
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(now)

	for _, set := range sets {
		println(set)
	}

	fmt.Printf("Took %s\n", elapsed)
}
