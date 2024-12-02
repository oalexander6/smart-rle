// Package rle implements the necessary functionality to encode and decode byte arrays
// using a smart run-length encoding scheme. In this scheme, consecutive identical
// bytes will be replaced by <delimiter><length><delimiter><original byte> if and only
// if this results in a shorter output.
//
// Examples:
//   asdf > asdf
//   asdddf > asdddf
//   asddddf > as.4.df
//
// In order for this encoding to work, there must be a byte value that is never present
// in the input to use as a delimiter. If this byte is ever found in the input to encode,
// it will return an error.
//
// The goal of this implementation is to optimize the worst-case compression ratio of
// this encoding. The result of encoding an array of bytes with this implementation
// will be, at worst, the same length as the input. In a traditional RLE scheme, each
// byte is prefixed by a count, which could result in an output that is up to double
// the size of the input for a highly random input.
//
// An additional optimization is used by encoding the run lengths using base-62, where
// the digits 0-9, a-z, and A-Z are all used. This increases the number of lengths that
// can be described with a given number of characters, allowing for the swapping of
// consecutive bytes with a run length to occur more frequently than with base-10. This
// also further improves the compression ratio.

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

// Encode encodes the input byte array using the provided delimiter byte. If the
// input byte array contains the delimiter byte, an error will be returned.
func Encode(input []byte, delim byte) ([]byte, error) {
	result := make([]byte, 0)

	if len(input) == 0 {
		return result, nil
	}

	last := input[0]
	count := 1

	// loop until length of input because we need to capture the final character
	// and process it
	for i := 1; i <= len(input); i++ {
		if i < len(input) && input[i] == delim {
			return []byte{}, ErrDelimitterFound
		}

		if i < len(input) && input[i] == last {
			count++
			continue
		}

		// we will always get to here for the last character
		base62Val := ToBase62(count)

		if len(base62Val)+2 < count {
			result = append(result, fmt.Sprintf("%c%s%c%c", delim, base62Val, delim, last)...)
		} else {
			result = append(result, makeBytes(count, last)...)
		}

		if i < len(input) {
			count = 1
			last = input[i]
		}
	}

	return result, nil
}

// Decode decodes the input byte array from Smart RLE encoding to plain bytes. It must
// be called with the same delimiter byte that the input was encoded with. If the input
// is malformed, an error will be returned.
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

		// this is to catch the case where we end on a run length without a trailing byte
		if i >= len(input) {
			return []byte{}, ErrMalformedInput
		}

		newData := makeBytes(runLength, input[i])
		result = append(result, newData...)
		i += 1
	}

	return result, nil
}

// getNextRunLength reads from a delimiter to the next delimiter, returning an
// error if there is no second delimiter or an invalid character for base-62 is
// found. It returns the found run length, the number of bytes in the run length
// and the opening and closing delimiter, and any errors.
func getNextRunLength(input []byte, delim byte) (runLength int, bytesRead int, err error) {
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

	runLength, err = FromBase62(string(lengthBytes))
	if err != nil {
		return 0, 0, err
	}

	bytesRead = 2 + len(lengthBytes)
	return runLength, bytesRead, nil
}

// ToBase62 takes the provided number and converts it to a base-62 number string.
// Base-62 uses the characters 0-9, a-z, and A-Z in that order to describe digits
// from 0-62. In most cases, you will not need to use this function.
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

// FromBase62 takes the provided base-62 number string and converts it to an integer.
// If there are invalid characters in the input, an error will be returned.
// Base-62 uses the characters 0-9, a-z, and A-Z in that order to describe digits
// from 0-62. In most cases, you will not need to use this function.
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

// makeBytes creates a slice of bytes of length n filled with bytes of val
func makeBytes(n int, val byte) []byte {
	result := make([]byte, n)
	for i := range result {
		result[i] = val
	}
	return result
}

// the base 62 digits laid out with their index as their underlying value
var base62Digits = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D',
	'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
	'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X',
	'Y', 'Z',
}
