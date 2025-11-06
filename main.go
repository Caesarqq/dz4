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
	run()
}

func run() {
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

func parseArgs(args []string) (uniq.Options, string, string, error) {
	opts := uniq.Options{}
	inputFile, outputFile := "", ""

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isFlag(arg) {
			var err error
			i, err = handleFlag(arg, &opts, args, i)
			if err != nil {
				return opts, "", "", err
			}
		} else {
			if inputFile == "" {
				inputFile = arg
			} else {
				outputFile = arg
			}
		}
	}

	if err := checkExclusiveFlags(opts); err != nil {
		return opts, "", "", err
	}

	return opts, inputFile, outputFile, nil
}

func isFlag(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

func handleFlag(arg string, opts *uniq.Options, args []string, i int) (int, error) {
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
			return i, fmt.Errorf("отсутствует значение для -f")
		}
		val, err := strconv.Atoi(args[i+1])
		if err != nil {
			return i, fmt.Errorf("неверное значение для -f: %v", args[i+1])
		}
		opts.SkipFields = val
		i++
	case "-s":
		if i+1 >= len(args) {
			return i, fmt.Errorf("отсутствует значение для -s")
		}
		val, err := strconv.Atoi(args[i+1])
		if err != nil {
			return i, fmt.Errorf("неверное значение для -s: %v", args[i+1])
		}
		opts.SkipChars = val
		i++
	}
	return i, nil
}

func checkExclusiveFlags(opts uniq.Options) error {
	count := 0
	if opts.Count {
		count++
	}
	if opts.Duplicate {
		count++
	}
	if opts.Unique {
		count++
	}
	if count > 1 {
		return fmt.Errorf("Ошибка: флаги -c, -d и -u нельзя использовать одновременно")
	}
	return nil
}
