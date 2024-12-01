package rle

import (
	"bytes"
	"errors"
	"testing"
)

func TestDecode(t *testing.T) {
	input := []byte("as3df")
	result, err := Decode(input, '.')
	if err != nil || !bytes.Equal(result, input) {
		t.FailNow()
	}

	input = []byte("as.3.df")
	result, err = Decode(input, '.')
	if err != nil || !bytes.Equal(result, []byte("asdddf")) {
		t.FailNow()
	}

	input = []byte("as.12.df")
	result, err = Decode(input, '.')
	if err != nil || !bytes.Equal(result, []byte("asddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddf")) {
		t.FailNow()
	}

	input = []byte("as.3.d.4.f")
	result, err = Decode(input, '.')
	if err != nil || !bytes.Equal(result, []byte("asdddffff")) {
		t.FailNow()
	}

	input = []byte("as.3df")
	_, err = Decode(input, '.')
	if !errors.Is(err, ErrMalformedInput) {
		t.FailNow()
	}

	input = []byte("as.3?.df")
	_, err = Decode(input, '.')
	if !errors.Is(err, ErrMalformedInput) {
		t.FailNow()
	}
}

func TestEncode(t *testing.T) {
	result, err := Encode([]byte("asdf"), '.')
	if err != nil || !bytes.Equal(result, []byte("asdf")) {
		t.FailNow()
	}

	result, err = Encode([]byte("asddf"), '.')
	if err != nil || !bytes.Equal(result, []byte("asddf")) {
		t.FailNow()
	}

	result, err = Encode([]byte("asdddf"), '.')
	if err != nil || !bytes.Equal(result, []byte("asdddf")) {
		t.FailNow()
	}

	result, err = Encode([]byte("asddddf"), '.')
	if err != nil || !bytes.Equal(result, []byte("as.4.df")) {
		t.FailNow()
	}
}

func TestToBase82(t *testing.T) {
	result := ToBase62(12)
	if result != "c" {
		t.FailNow()
	}

	result = ToBase62(94)
	if result != "1w" {
		t.FailNow()
	}

	result = ToBase62(12478)
	if result != "3fg" {
		t.FailNow()
	}
}

func TestFromBase62(t *testing.T) {
	if result, err := FromBase62("c"); err != nil || result != 12 {
		t.FailNow()
	}

	if result, err := FromBase62("1w"); err != nil || result != 94 {
		t.FailNow()
	}

	if result, err := FromBase62("3fg"); err != nil || result != 12478 {
		t.FailNow()
	}
}
