package main

import (
	"fmt"
	"log"

	"github.com/simon-siggaard/lego/pkg/brick"
)

func main() {
	users, err := brick.UserCollections()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", users)
}
