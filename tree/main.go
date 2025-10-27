package main

import (
	"os"
)

const cmdError = "usage: go run . dir [-f]"

func main() {
	argc := len(os.Args)

	if argc == 1 || argc > 3 {
		panic(cmdError)
	}

	dirPath := os.Args[1]
	isShowFiles := false

	if argc == 3 {
		isShowFilesFlag := os.Args[2]
		if isShowFilesFlag != "-f" {
			panic(cmdError)
		}
		isShowFiles = true
	}

	if err := dirTree(os.Stdout, dirPath, isShowFiles); err != nil {
		panic(err)
	}
}
