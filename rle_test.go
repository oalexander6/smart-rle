package rle

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func ExampleDecode() {
	input := []byte("as.3.df")
	delimiter := byte('.')

	decoded, err := Decode(input, delimiter)

	fmt.Println(decoded, err)

	// Output:
	// asdddf nil
}

func ExampleEncode() {
	input := []byte("asdddf")
	delimiter := byte('.')

	encoded, err := Encode(input, delimiter)

	fmt.Println(encoded, err)

	// Output:
	// as.3.df nil
}

func TestDecode(t *testing.T) {
	result, err := Decode([]byte("as3df"), '.')
	if err != nil || !bytes.Equal(result, []byte("as3df")) {
		t.FailNow()
	}

	result, err = Decode([]byte("as.3.df"), '.')
	if err != nil || !bytes.Equal(result, []byte("asdddf")) {
		t.FailNow()
	}

	result, err = Decode([]byte("as.12.df"), '.')
	if err != nil || !bytes.Equal(result, []byte("asddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddf")) {
		t.FailNow()
	}

	result, err = Decode([]byte("as.3.d.4.f"), '.')
	if err != nil || !bytes.Equal(result, []byte("asdddffff")) {
		t.FailNow()
	}

	_, err = Decode([]byte("as.3df"), '.')
	if !errors.Is(err, ErrMalformedInput) {
		t.FailNow()
	}

	_, err = Decode([]byte("as.3?.df"), '.')
	if !errors.Is(err, ErrMalformedInput) {
		t.FailNow()
	}

	_, err = Decode([]byte("asdf.4."), '.')
	if !errors.Is(err, ErrMalformedInput) {
		t.FailNow()
	}

	result, err = Decode([]byte(".4.a.4.b"), '.')
	if err != nil || !bytes.Equal(result, []byte("aaaabbbb")) {
		t.FailNow()
	}

	result, err = Decode([]byte(".1Q.aewfffiohwef.C.b"), '.')
	if err != nil || !bytes.Equal(result, []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaewfffiohwefbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")) {
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

	result, err = Encode([]byte("aaaabbbb"), '.')
	if err != nil || !bytes.Equal(result, []byte(".4.a.4.b")) {
		t.FailNow()
	}

	result, err = Encode([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaewfffiohwefbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), '.')
	if err != nil || !bytes.Equal(result, []byte(".1Q.aewfffiohwef.C.b")) {
		t.FailNow()
	}

	_, err = Encode([]byte("asdf.dfd"), '.')
	if !errors.Is(err, ErrDelimitterFound) {
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
