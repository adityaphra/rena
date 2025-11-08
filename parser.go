package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
)

func ParseCommands(commands []string) ([]command, error) {
	parsed := make([]command, len(commands))
	for i, c := range commands {
		if c == "" {
			continue
		}
		result, err := parseCommand(c)
		if err != nil {
			return nil, err
		}
		parsed[i] = result
	}
	return parsed, nil
}

var parsers = map[byte]func([]string) (command, error){
	's': parseSearch,
	'd': parseDelete,
	't': parseTemplate,
	'm': parseMove,
}

func parseCommand(text string) (command, error) {
	if len(text) < 2 {
		return nil, fmt.Errorf("Invalid format: %v", text)
	}

	var separator byte
	for _, sep := range []byte(",;/|\t") {
		if sep == text[1] {
			separator = sep
			break
		}
	}
	if separator == 0 {
		return nil, fmt.Errorf("Separator is not supported: %v", text[1])
	}

	reader := csv.NewReader(strings.NewReader(text))
	reader.Comma = rune(separator)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Read error: '%v'", text)
	}

	commandType := text[0]
	for _, record := range records {
		parser, ok := parsers[commandType]
		if ok {
			result, err := parser(record)
			if err != nil {
				return nil, fmt.Errorf("%v: '%v'", err, text)
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("Uknown command: '%v'", text)
}

// format: s, search, replace, flags
// index   0  1       2        3
func parseSearch(record []string) (command, error) {
	length := len(record)
	if length < 3 || length > 4 {
		return nil, errors.New("Invalid format")
	}

	search := record[1]
	replace := record[2]
	if length == 3 {
		return searchCommand{search: search, replace: replace}, nil
	}

	flags := record[3]
	matchCase := strings.Contains(flags, "m")
	regexMode := strings.Contains(flags, "r")
	whiteSpace := strings.Contains(flags, "w")
	return searchCommand{search, replace, matchCase, regexMode, whiteSpace}, nil
}

// format: d, value, flags
// index   0  1      2
func parseDelete(record []string) (command, error) {
	length := len(record)
	if length < 2 || length > 3 {
		return nil, errors.New("Invalid format")
	}

	value := record[1]
	if length == 2 {
		return deleteCommand{value: value}, nil
	}

	flags := record[2]
	matchCase := strings.Contains(flags, "m")
	regexMode := strings.Contains(flags, "r")
	whiteSpace := strings.Contains(flags, "w")
	return deleteCommand{value, matchCase, regexMode, whiteSpace}, nil
}

// format: t, template
// index   0  1
func parseTemplate(record []string) (command, error) {
	if len(record) != 2 {
		return nil, errors.New("Invalid format")
	}
	return templateCommand{record[1]}, nil
}

// format: m, pattern, destinationDirectory, flags
// index   0  1        2                     3
func parseMove(record []string) (command, error) {
	length := len(record)
	if length < 3 || length > 4 {
		return nil, errors.New("Invalid format")
	}

	pattern := record[1]
	destinationDir := record[2]
	if length == 3 {
		return moveCommand{pattern: pattern, destinationDir: destinationDir}, nil
	}

	flags := record[3]
	matchCase := strings.Contains(flags, "m")
	regexMode := strings.Contains(flags, "r")
	return moveCommand{pattern, destinationDir, matchCase, regexMode}, nil
}
