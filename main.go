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
		printUsage()
		return
	}

	lines, err := getLines(inputFile)
	if err != nil {
		fmt.Println("Ошибка при чтении входного файла:", err)
		return
	}

	result := uniq.ProcessLines(lines, opts)

	if err := writeOutput(outputFile, result); err != nil {
		fmt.Println("Ошибка при записи в выходной файл:", err)
	}
}

func printUsage() {
	fmt.Println("Использование: uniq [-c | -d | -u] [-i] [-f число] [-s символы] [входной_файл [выходной_файл]]")
}

func getLines(inputFile string) ([]string, error) {
	var reader io.Reader
	if inputFile != "" {
		file, err := os.Open(inputFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}
	return readLines(reader)
}

func writeOutput(outputFile string, lines []string) error {
	var writer io.Writer
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = file
	} else {
		writer = os.Stdout
	}
	return writeLines(writer, lines)
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
				return opts, "", "", fmt.Errorf("отсутствует значение для -f")
			}
			val, err := strconv.Atoi(args[i+1])
			if err != nil {
				return opts, "", "", fmt.Errorf("неверное значение для -f: %v", args[i+1])
			}
			opts.SkipFields = val
			i++
		case "-s":
			if i+1 >= len(args) {
				return opts, "", "", fmt.Errorf("отсутствует значение для -s")
			}
			val, err := strconv.Atoi(args[i+1])
			if err != nil {
				return opts, "", "", fmt.Errorf("неверное значение для -s: %v", args[i+1])
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
		return opts, "", "", fmt.Errorf("Ошибка: флаги -c, -d и -u нельзя использовать одновременно")
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
		if _, err := fmt.Fprintln(writer, line); err != nil {
			return err
		}
	}
	return writer.Flush()
}
