package main

import (
	"testing"
)

func TestBasicUniq(t *testing.T) {
	input := []string{
		"line1",
		"line1",
		"line2",
		"line3",
		"line3",
	}

	result := processLines(input, false, false, false, false, 0, 0)

	if len(result) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(result))
	}

	if result[0] != "line1" {
		t.Errorf("Expected 'line1', got '%s'", result[0])
	}
	if result[1] != "line2" {
		t.Errorf("Expected 'line2', got '%s'", result[1])
	}
	if result[2] != "line3" {
		t.Errorf("Expected 'line3', got '%s'", result[2])
	}
}

func TestCountFlag(t *testing.T) {
	input := []string{
		"hello",
		"hello",
		"world",
	}

	result := processLines(input, true, false, false, false, 0, 0)

	if len(result) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(result))
	}

	if result[0] != "2 hello" {
		t.Errorf("Expected '2 hello', got '%s'", result[0])
	}
}

func TestDuplicateFlag(t *testing.T) {
	input := []string{
		"dup",
		"dup",
		"unique",
	}

	result := processLines(input, false, true, false, false, 0, 0)

	if len(result) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result))
	}

	if result[0] != "dup" {
		t.Errorf("Expected 'dup', got '%s'", result[0])
	}
}

func TestUniqueFlag(t *testing.T) {
	input := []string{
		"dup",
		"dup",
		"unique",
	}

	result := processLines(input, false, false, true, false, 0, 0)

	if len(result) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result))
	}

	if result[0] != "unique" {
		t.Errorf("Expected 'unique', got '%s'", result[0])
	}
}

func TestIgnoreCaseFlag(t *testing.T) {
	input := []string{
		"Hello",
		"hello",
		"HELLO",
	}

	result := processLines(input, false, false, false, true, 0, 0)

	if len(result) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result))
	}
}

func TestSkipFieldsFlag(t *testing.T) {
	input := []string{
		"a test",
		"b test",
		"c other",
	}

	result := processLines(input, false, false, false, false, 1, 0)

	if len(result) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(result))
	}
}

func TestSkipCharsFlag(t *testing.T) {
	input := []string{
		"atest",
		"btest",
		"ctest",
	}

	result := processLines(input, false, false, false, false, 0, 1)

	if len(result) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result))
	}
}

func TestEmptyInput(t *testing.T) {
	input := []string{}
	result := processLines(input, false, false, false, false, 0, 0)

	if len(result) != 0 {
		t.Errorf("Expected 0 lines, got %d", len(result))
	}
}
