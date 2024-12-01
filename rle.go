package rle

import (
	"errors"
	"fmt"
	"math"
	"slices"
)

var (
	ErrDelimitterFound = errors.New("delimitter value found in the input")
	ErrMalformedInput  = errors.New("malformed input")
	ErrInvalidBase     = errors.New("must supply base between 2-62")
)

func Encode(input []byte, delim byte) ([]byte, error) {
	result := make([]byte, 0)

	if len(input) == 0 {
		return result, nil
	}

	last := input[0]
	count := 1

	for i := 1; i < len(input); i++ {
		if input[i] == last {
			count++
			continue
		}

		base62Val := ToBase62(count)

		if len(base62Val)+2 < count {
			result = append(result, fmt.Sprintf("%c%s%c%c", delim, base62Val, delim, last)...)
		} else {
			result = append(result, makeBytes(count, last)...)
		}

		count = 1
		last = input[i]
	}

	result = append(result, makeBytes(count, last)...)

	return result, nil
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

	length, err := FromBase62(string(lengthBytes))
	if err != nil {
		return 0, 0, err
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

var base62Digits = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D',
	'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
	'Y', 'Z',
}

func ToBase62(val int) string {
	// number of digits to represent a number in base 'b' is logb(number) + 1
	numDigits := int(math.Log10(float64(val))/math.Log10(float64(62))) + 1

	result := make([]byte, numDigits)
	remaining := val
	i := 0
	for remaining > 0 {
		result[numDigits-i-1] = base62Digits[remaining%62]
		remaining = remaining / 62
		i += 1
	}

	return string(result)
}

func FromBase62(val string) (int, error) {
	result := 0

	for i := 0; i < len(val); i++ {
		curr := slices.Index(base62Digits, val[i])
		if curr == -1 {
			return 0, ErrMalformedInput
		}
		result += int(math.Pow(float64(62), float64(len(val)-i-1))) * curr
	}

	return result, nil
}
