package rle_test

import (
	"bytes"
	"errors"
	"testing"

	rle "github.com/oalexander6/smart-rle"
)

func TestDecode(t *testing.T) {
	input := []byte("as3df")
	result, err := rle.Decode(input, '.')
	if err != nil || !bytes.Equal(result, input) {
		t.FailNow()
	}

	input = []byte("as.3.df")
	result, err = rle.Decode(input, '.')
	if err != nil || !bytes.Equal(result, []byte("asdddf")) {
		t.FailNow()
	}

	input = []byte("as.3.d.4.f")
	result, err = rle.Decode(input, '.')
	if err != nil || !bytes.Equal(result, []byte("asdddffff")) {
		t.FailNow()
	}

	input = []byte("as.3df")
	_, err = rle.Decode(input, '.')
	if !errors.Is(err, rle.ErrMalformedInput) {
		t.FailNow()
	}

	input = []byte("as.3a.df")
	_, err = rle.Decode(input, '.')
	if !errors.Is(err, rle.ErrMalformedInput) {
		t.FailNow()
	}
}
