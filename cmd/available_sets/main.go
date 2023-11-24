package main

import "github.com/simon-siggaard/lego/pkg/brick"

func main() {
	sets, _ := brick.AvailableSets("d174c807-8880-4f49-866b-6e1ec6527ccf")

	for _, set := range sets {
		println(set)
	}
}
