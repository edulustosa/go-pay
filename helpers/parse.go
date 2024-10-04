package helpers

import (
	"errors"
	"strconv"
	"strings"
)

func calculateDigit(document string, factor int) (int, error) {
	total := 0
	for i := 0; i < factor-1; i++ {
		n, err := strconv.Atoi(string(document[i]))
		if err != nil {
			return 0, errors.New("invalid document")
		}
		total += n * (factor - i)
	}

	rest := (total * 10) % 11
	if rest == 10 {
		rest = 0
	}
	return rest, nil
}

// Verify if a document is valid CPF
func ParseDocument(document string) error {
	document = strings.ReplaceAll(document, ".", "")
	document = strings.ReplaceAll(document, "-", "")

	if len(document) != 11 {
		return errors.New("document must have 11 characters")
	}

	if strings.Repeat(string(document[0]), 11) == document {
		return errors.New("document must not have all characters equal")
	}

	digit1, err := calculateDigit(document, 10)
	if err != nil {
		return err
	}
	digit2, err := calculateDigit(document[:9]+strconv.Itoa(digit1), 11)
	if err != nil {
		return err
	}

	if digit1 != int(document[9]-'0') || digit2 != int(document[10]-'0') {
		return errors.New("invalid document")
	}

	return nil
}
