package rle

import (
	"errors"
	"strconv"
)

var (
	ErrDelimitterFound = errors.New("delimitter value found in the input")
	ErrMalformedInput  = errors.New("malformed input")
)

func Encode(input []byte, delim byte) ([]byte, error) {
	return []byte{}, nil
}

func Decode(input []byte, delim byte) ([]byte, error) {
	result := make([]byte, 0)

	i := 0

	for i < len(input) {
		if input[i] != delim {
			result = append(result, input[i])
			i += 1
			continue
		}

		runLength, bytesRead, err := getNextRunLength(input[i:], delim)
		if err != nil {
			return nil, err
		}

		i += bytesRead

		newData := makeBytes(runLength, input[i])

		result = append(result, newData...)

		i += 1
	}

	return result, nil
}

func getNextRunLength(input []byte, delim byte) (int, int, error) {
	if len(input) < 3 || input[0] != delim {
		return 0, 0, ErrMalformedInput
	}

	lengthBytes := make([]byte, 0)
	foundClosingDelim := false

	for i := 1; i < len(input); i++ {
		if input[i] == delim {
			foundClosingDelim = true
			break
		}

		lengthBytes = append(lengthBytes, input[i])
	}

	if !foundClosingDelim || len(lengthBytes) == 0 {
		return 0, 0, ErrMalformedInput
	}

	length, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return 0, 0, ErrMalformedInput
	}

	bytesRead := 2 + len(lengthBytes)

	return length, bytesRead, nil
}

func makeBytes(n int, val byte) []byte {
	result := make([]byte, n)
	for i := range result {
		result[i] = val
	}
	return result
}
