package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func FastSearch(out io.Writer) {
	file, err := os.Open("./data/users.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	seenBrowsers := make(map[string]struct{}, 100)
	names := make([]string, 0, 100)
	emails := make([]string, 0, 100)
	nums := make([]int, 0, 100)

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		user := make(map[string]any)
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}

		hasAndroid := false
		hasMSIE := false
		browsers, ok := user["browsers"].([]any)
		if !ok {
			continue
		}

		for _, rawBrowser := range browsers {
			browser, ok := rawBrowser.(string)
			if !ok {
				continue
			}

			isAndroid := strings.Contains(browser, "Android")
			if isAndroid {
				hasAndroid = true
			}

			isMSIE := strings.Contains(browser, "MSIE")
			if isMSIE {
				hasMSIE = true
			}

			if _, has := seenBrowsers[browser]; !has && (isAndroid || isMSIE) {
				seenBrowsers[browser] = struct{}{}
			}
		}

		if !hasAndroid || !hasMSIE {
			continue
		}

		name, ok := user["name"].(string)
		if !ok {
			continue
		}
		names = append(names, name)

		email, ok := user["email"].(string)
		if !ok {
			continue
		}
		formattedEmail := strings.Replace(email, "@", " [at] ", 1)
		emails = append(emails, formattedEmail)

		nums = append(nums, i)
	}

	fmt.Fprintln(out, "found users:")
	for i, name := range names {
		email := emails[i]
		num := nums[i]
		foundUser := fmt.Sprintf("[%d] %s <%s>", num, name, email)
		fmt.Fprintln(out, foundUser)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
