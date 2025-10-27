package main

import (
	"os"
	"strings"
	"testing"
)

type formatDirEntryArgs struct {
	name   string
	bSize  int64
	seps   []bool
	isDir  bool
	isLast bool
}

func TestDirAndNotFormatDirEntry(T *testing.T) {
	cases := []formatDirEntryArgs{
		{
			name: "file1",
		},
		{
			name:  "file2",
			isDir: true,
		},
	}
	expectedResults := []string{
		cross + hline + "file1 (empty)",
		cross + hline + "file2",
	}

	for i := range cases {
		result := formatDirEntry(cases[i].name, cases[i].bSize, cases[i].seps, cases[i].isDir, cases[i].isLast)
		if result != expectedResults[i] {
			T.Errorf("actual: '%s'\nexpected: '%s'\n", result, expectedResults[i])
		}
	}
}

func TestEmtyAndNotFormatDirEntry(T *testing.T) {
	cases := []formatDirEntryArgs{
		{
			name: "file1",
		},
		{
			name:  "file2",
			bSize: 100,
		},
	}
	expectedResults := []string{
		cross + hline + "file1 (empty)",
		cross + hline + "file2 (100b)",
	}

	for i := range cases {
		result := formatDirEntry(cases[i].name, cases[i].bSize, cases[i].seps, cases[i].isDir, cases[i].isLast)
		if result != expectedResults[i] {
			T.Errorf("actual: '%s'\nexpected: '%s'\n", result, expectedResults[i])
		}
	}
}

func TestDeep0FormatDirEntry(T *testing.T) {
	cases := []formatDirEntryArgs{
		{
			name: "file1",
			seps: []bool{},
		},
		{
			name:   "file2",
			seps:   []bool{},
			isLast: true,
		},
	}
	expectedResults := []string{
		cross + hline + "file1 (empty)",
		end + hline + "file2 (empty)",
	}

	for i := range cases {
		result := formatDirEntry(cases[i].name, cases[i].bSize, cases[i].seps, cases[i].isDir, cases[i].isLast)
		if result != expectedResults[i] {
			T.Errorf("actual: '%s'\nexpected: '%s'\n", result, expectedResults[i])
		}
	}
}

func TestDeep1FormatDirEntry(T *testing.T) {
	cases := []formatDirEntryArgs{
		{
			name: "file1",
			seps: []bool{false},
		},
		{
			name:   "file1",
			seps:   []bool{false},
			isLast: true,
		},
		{
			name: "file1",
			seps: []bool{true},
		},
		{
			name:   "file1",
			seps:   []bool{true},
			isLast: true,
		},
	}
	expectedResults := []string{
		space + " " + cross + hline + "file1 (empty)",
		space + " " + end + hline + "file1 (empty)",
		vline + space + cross + hline + "file1 (empty)",
		vline + space + end + hline + "file1 (empty)",
	}

	for i := range cases {
		result := formatDirEntry(cases[i].name, cases[i].bSize, cases[i].seps, cases[i].isDir, cases[i].isLast)
		if result != expectedResults[i] {
			T.Errorf("actual: '%s'\nexpected: '%s'\n", result, expectedResults[i])
		}
	}
}

func TestDeep2FormatDirEntry(T *testing.T) {
	cases := []formatDirEntryArgs{
		{
			name: "file1",
			seps: []bool{false, false},
		},
		{
			name: "file1",
			seps: []bool{false, true},
		},
		{
			name: "file1",
			seps: []bool{true, false},
		},
		{
			name: "file1",
			seps: []bool{true, true},
		},
	}
	expectedResults := []string{
		space + " " + space + " " + cross + hline + "file1 (empty)",
		space + " " + vline + space + cross + hline + "file1 (empty)",
		vline + space + space + " " + cross + hline + "file1 (empty)",
		vline + space + vline + space + cross + hline + "file1 (empty)",
	}

	for i := range cases {
		result := formatDirEntry(cases[i].name, cases[i].bSize, cases[i].seps, cases[i].isDir, cases[i].isLast)
		if result != expectedResults[i] {
			T.Errorf("actual: '%s'\nexpected: '%s'\n", result, expectedResults[i])
		}
	}
}

// По хорошему тесты должны генерировать структуру каталога, а не сделанную заранее
func TestOnlyDirsDirPrint(T *testing.T) {
	expected := `
└───dir2
    ├───dir3
    │   ├───dir 6
    │   │   └───dir7
    │   │       └───dir8
    │   └───dir5
    └───dir4
`
	expected = strings.TrimSpace(expected)

	outputFileName := "tmp"
	outputFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	if err := dirTree(outputFile, "test-dir1", false); err != nil {
		panic(err)
	}

	if err := outputFile.Close(); err != nil {
		panic(err)
	}

	data, err := os.ReadFile(outputFileName)
	if err != nil {
		panic(err)
	}
	result := strings.TrimSpace(string(data))

	if result != expected {
		T.Errorf("actual: '%s'\nexpected: '%s'\n", result, expected)
	}

	os.Remove(outputFileName)
}

func TestAllFilesDirPrint(T *testing.T) {
	expected := `
├───a (empty)
├───b (20b)
├───dir2
│   ├───dir3
│   │   ├───c (empty)
│   │   ├───d (empty)
│   │   ├───dir 6
│   │   │   └───dir7
│   │   │       └───dir8
│   │   │           └───xxx (empty)
│   │   ├───dir5
│   │   ├───file5 (5b)
│   │   ├───file6 (empty)
│   │   └───file7 (27b)
│   ├───dir4
│   │   └───file4 (empty)
│   └───flie3 (empty)
├───file1 (empty)
└───file2 (empty)
`
	expected = strings.TrimSpace(expected)

	outputFileName := "tmp"
	outputFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	if err := dirTree(outputFile, "test-dir1", true); err != nil {
		panic(err)
	}

	if err := outputFile.Close(); err != nil {
		panic(err)
	}

	data, err := os.ReadFile(outputFileName)
	if err != nil {
		panic(err)
	}
	result := strings.TrimSpace(string(data))

	if result != expected {
		T.Errorf("actual: '%s'\nexpected: '%s'\n", result, expected)
	}

	os.Remove(outputFileName)
}
