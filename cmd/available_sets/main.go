package main

import (
	"fmt"
	"time"

	"github.com/simon-siggaard/lego/pkg/brick/sets"
)

func main() {
	now := time.Now()
	sets, _ := sets.AvailableSets("d174c807-8880-4f49-866b-6e1ec6527ccf")
	elapsed := time.Since(now)

	for _, set := range sets {
		println(set)
	}

	fmt.Printf("Took %s\n", elapsed)
}
