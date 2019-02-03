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
	files, err := collectData(src)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Println(f)
		if out, err := f.ExecutePlfit(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println(out)
		}
	}
}
