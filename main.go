package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var countFlag bool
	var duplicateFlag bool
	var uniqueFlag bool
	var ignoreCaseFlag bool
	var skipFields int
	var skipChars int
	var inputFile string
	var outputFile string

	args := os.Args[1:]
	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == "-c" {
			countFlag = true
			i++
		} else if arg == "-d" {
			duplicateFlag = true
			i++
		} else if arg == "-u" {
			uniqueFlag = true
			i++
		} else if arg == "-i" {
			ignoreCaseFlag = true
			i++
		} else if arg == "-f" {
			i++
			if i < len(args) {
				skipFields, _ = strconv.Atoi(args[i])
				i++
			}
		} else if arg == "-s" {
			i++
			if i < len(args) {
				skipChars, _ = strconv.Atoi(args[i])
				i++
			}
		} else {
			if inputFile == "" {
				inputFile = arg
			} else {
				outputFile = arg
			}
			i++
		}
	}

	flagCount := 0
	if countFlag {
		flagCount++
	}
	if duplicateFlag {
		flagCount++
	}
	if uniqueFlag {
		flagCount++
	}

	if flagCount > 1 {
		fmt.Println("Error: can't use -c, -d, -u together")
		fmt.Println("Usage: uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
		return
	}

	var lines []string
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	}

	result := processLines(lines, countFlag, duplicateFlag, uniqueFlag, ignoreCaseFlag, skipFields, skipChars)
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		for _, line := range result {
			fmt.Fprintln(file, line)
		}
	} else {
		for _, line := range result {
			fmt.Println(line)
		}
	}
}

func processLines(lines []string, countFlag, duplicateFlag, uniqueFlag,
	ignoreCaseFlag bool, skipFields, skipChars int) []string {
	if len(lines) == 0 {
		return []string{}
	}
	var result []string
	currentLine := lines[0]
	currentKey := makeCompareKey(lines[0], ignoreCaseFlag, skipFields, skipChars)
	count := 1
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		key := makeCompareKey(line, ignoreCaseFlag, skipFields, skipChars)
		if key == currentKey {
			count++
		} else {
			if shouldPrint(count, duplicateFlag, uniqueFlag) {
				if countFlag {
					result = append(result, fmt.Sprintf("%d %s", count, currentLine))
				} else {
					result = append(result, currentLine)
				}
			}
			currentLine = line
			currentKey = key
			count = 1
		}
	}

	if shouldPrint(count, duplicateFlag, uniqueFlag) {
		if countFlag {
			result = append(result, fmt.Sprintf("%d %s", count, currentLine))
		} else {
			result = append(result, currentLine)
		}
	}

	return result
}

func makeCompareKey(line string, ignoreCase bool, skipFields, skipChars int) string {
	key := line

	if skipFields > 0 {
		parts := strings.Fields(key)
		if skipFields < len(parts) {
			key = strings.Join(parts[skipFields:], " ")
		} else {
			key = ""
		}
	}

	if skipChars > 0 {
		if len(key) > skipChars {
			key = key[skipChars:]
		} else {
			key = ""
		}
	}

	if ignoreCase {
		key = strings.ToLower(key)
	}

	return key
}

func shouldPrint(count int, duplicateFlag, uniqueFlag bool) bool {
	if duplicateFlag {
		return count > 1
	}
	if uniqueFlag {
		return count == 1
	}
	return true
}
