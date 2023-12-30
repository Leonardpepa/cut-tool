package main

import (
	"errors"
	"testing"
)

func TestValidateFlags(t *testing.T) {
	t.Run("should validate the flags", func(t *testing.T) {
		defaultDelimiter := "\t"
		testValidFlags(t, "-f", "", "", defaultDelimiter)
		testValidFlags(t, "", "-b", "", defaultDelimiter)
		testValidFlags(t, "", "", "-c", defaultDelimiter)
		testValidFlags(t, "-c", "", "", ",")

		testInvalidFlags(t, "-f", "-b", "-c", defaultDelimiter, toManyListArguments)
		testInvalidFlags(t, "", "-b", "", ",", delimiterError)
		testInvalidFlags(t, "", "", "", defaultDelimiter, noFlagSpecified)

	})
}

func testValidFlags(t *testing.T, fields, bytes, chars, delimiter string) {
	err := validateFlags(&fields, &bytes, &chars, &delimiter)
	if err != nil {
		t.Error(err)
	}
}

func testInvalidFlags(t *testing.T, fields, bytes, chars, delimiter string, expected error) {
	err := validateFlags(&fields, &bytes, &chars, &delimiter)
	if err == nil {
		t.Errorf("parsed wrong flags, expected %v", expected)
	}

	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}
