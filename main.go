package main /* import "ivo.qa/plfit-tablegen" */

import (
	"fmt"
	"os"
	"ivo.qa/plfit-tablegen/fittools"
)

// This program executes `plfit` for all files in the specified folders,
// collecting data from the filenames (frame, type, period, and struct number)
// and xmin and alpha, found in plfit's output.
// Then it outputs the collected data as a formatted table to stdin.
// In order to run the program, make sure you have plfit preinstalled.

const usageInfo = `Usage: plfit-tablegen <source path> [--silent]
When silent is on, only the table will be output.
Source: https://github.com/ivaivalous/plfit-tablegen`

func main() {
	if len(os.Args) == 1 {
		fmt.Println(usageInfo)
		os.Exit(1)
	}
	src := os.Args[1]
	silent := len(os.Args) >= 3 && os.Args[2] == "--silent"
	files, err := fittools.CollectData(src)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}

	for _, f := range files {
		if !silent {
			fmt.Println("Processing", f.Filename)
		}
		if err := f.ExecutePlfit(); err != nil {
			fmt.Println(err)
		}
	}

	fittools.AsTable(os.Stdout, files)
}
