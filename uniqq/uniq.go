package uniq

import (
	"fmt"
	"strings"
)

type Options struct {
	Count      bool
	Duplicate  bool
	Unique     bool
	IgnoreCase bool
	SkipFields int
	SkipChars  int
}

func ProcessLines(lines []string, opts Options) []string {
	if len(lines) == 0 {
		return []string{}
	}
	var result []string
	currentLine := lines[0]
	currentKey := makeKey(lines[0], opts)
	count := 1
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		key := makeKey(line, opts)

		if key == currentKey {
			count++
		} else {
			if shouldPrint(count, opts) {
				result = append(result, formatLine(currentLine, count, opts))
			}
			currentLine = line
			currentKey = key
			count = 1
		}
	}

	if shouldPrint(count, opts) {
		result = append(result, formatLine(currentLine, count, opts))
	}
	return result
}

func makeKey(line string, opts Options) string {
	key := line
	if opts.SkipFields > 0 {
		parts := strings.Fields(key)
		if opts.SkipFields < len(parts) {
			key = strings.Join(parts[opts.SkipFields:], " ")
		} else {
			key = ""
		}
	}

	if opts.SkipChars > 0 {
		if len(key) > opts.SkipChars {
			key = key[opts.SkipChars:]
		} else {
			key = ""
		}
	}

	if opts.IgnoreCase {
		key = strings.ToLower(key)
	}

	return key
}

func shouldPrint(count int, opts Options) bool {
	if opts.Duplicate {
		return count > 1
	}
	if opts.Unique {
		return count == 1
	}
	return true
}

func formatLine(line string, count int, opts Options) string {
	if opts.Count {
		return fmt.Sprintf("%d %s", count, line)
	}
	return line
}
