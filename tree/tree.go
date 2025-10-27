package main

import (
	"fmt"
	"io"
	"os"
)

const (
	cross = "├"
	hline = "───"
	vline = "│"
	end   = "└"
	space = "   "
)

func dirTree(writer io.Writer, dirPath string, isPrintFiles bool) error {
	return dirTreeRecursive(writer, dirPath, isPrintFiles, []bool{})
}

func dirTreeRecursive(writer io.Writer, dirPath string, isPrintFiles bool, seps []bool) error {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	var filteredDirEntries []os.DirEntry
	if isPrintFiles {
		filteredDirEntries = dirEntries
	} else {
		for _, dirEntry := range dirEntries {
			if dirEntry.IsDir() {
				filteredDirEntries = append(filteredDirEntries, dirEntry)
			}
		}
	}

	for _, dirEntry := range filteredDirEntries {
		info, err := dirEntry.Info()
		if err != nil {
			return err
		}

		isLast := dirEntry == filteredDirEntries[len(filteredDirEntries)-1]

		if err := printDirEntry(writer, info.Name(), info.Size(), seps, info.IsDir(), isLast); err != nil {
			return err
		}

		nextDirPath := dirPath + string(os.PathSeparator) + info.Name()

		nextSeps := make([]bool, len(seps)+1)
		copy(nextSeps, seps)
		nextSeps[len(seps)] = !isLast

		if info.IsDir() {
			if err := dirTreeRecursive(writer, nextDirPath, isPrintFiles, nextSeps); err != nil {
				return err
			}
		}
	}

	return nil
}

func printDirEntry(writer io.Writer, name string, bSize int64, seps []bool, isDir bool, isLast bool) error {
	formattedDirEntry := formatDirEntry(name, bSize, seps, isDir, isLast)
	_, err := fmt.Fprintf(writer, "%s\n", formattedDirEntry)
	return err
}

func formatDirEntry(name string, bSize int64, seps []bool, isDir bool, isLast bool) string {
	decoration := cross
	if isLast {
		decoration = end
	}

	decoration += hline + name

	if !isDir {
		if bSize == 0 {
			decoration += " (empty)"
		} else {
			decoration += fmt.Sprintf(" (%db)", bSize)
		}
	}

	for deep := len(seps) - 1; deep >= 0; deep-- {
		if hasSep := seps[deep]; hasSep {
			decoration = vline + space + decoration
		} else {
			decoration = space + " " + decoration
		}
	}

	return decoration
}
