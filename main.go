package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	uniq "github.com/Caesarqq/dz4/uniqq"
)

func main() {
	opts, inputFile, outputFile, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		fmt.Println("Usage: uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]")
		return
	}

	var reader io.Reader
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Println("Error opening input file:", err)
			return
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	lines, err := readLines(reader)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	result := uniq.ProcessLines(lines, opts)
	var writer io.Writer
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		defer file.Close()
		writer = file
	} else {
		writer = os.Stdout
	}

	writeLines(writer, result)
}

func parseArgs(args []string) (uniq.Options, string, string, error) {
	var opts uniq.Options
	var inputFile, outputFile string
	i := 0

	for i < len(args) {
		arg := args[i]
		switch arg {
		case "-c":
			opts.Count = true
		case "-d":
			opts.Duplicate = true
		case "-u":
			opts.Unique = true
		case "-i":
			opts.IgnoreCase = true
		case "-f":
			if i+1 >= len(args) {
				return opts, "", "", fmt.Errorf("missing value for -f")
			}
			val, err := strconv.Atoi(args[i+1])
			if err != nil {
				return opts, "", "", fmt.Errorf("invalid -f value: %v", args[i+1])
			}
			opts.SkipFields = val
			i++
		case "-s":
			if i+1 >= len(args) {
				return opts, "", "", fmt.Errorf("missing value for -s")
			}
			val, err := strconv.Atoi(args[i+1])
			if err != nil {
				return opts, "", "", fmt.Errorf("invalid -s value: %v", args[i+1])
			}
			opts.SkipChars = val
			i++
		default:
			if inputFile == "" {
				inputFile = arg
			} else {
				outputFile = arg
			}
		}
		i++
	}

	flagCount := 0
	if opts.Count {
		flagCount++
	}
	if opts.Duplicate {
		flagCount++
	}
	if opts.Unique {
		flagCount++
	}
	if flagCount > 1 {
		return opts, "", "", fmt.Errorf("Error: can't use -c, -d, -u together")
	}

	return opts, inputFile, outputFile, nil
}

func readLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func writeLines(w io.Writer, lines []string) error {
	writer := bufio.NewWriter(w)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}
