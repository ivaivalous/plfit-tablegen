package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Usage: plfit-tablegen <source path>")
		os.Exit(1)
	}
	src := os.Args[1]
	paths, err := collectData(src)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(paths)
}
