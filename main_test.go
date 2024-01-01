package main

import (
	"bufio"
	"cut-tool/internal"
	"errors"
	"strings"
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
		testInvalidFlags(t, "", "", "-c", ",", delimiterError)
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

func TestExtractFields(t *testing.T) {
	t.Run("should return the correct fields", func(t *testing.T) {
		input := `f0	f1	f2	f3	f4
0	1	2	3	4
5	6	7	8	9
10	11	12	13	14
15	16	17	18	19
20	21	22	23	24`

		expected := "f1\n1\n6\n11\n16\n21\n"

		testOutput(t, input, "\t", "2", expected, extractFields)

		expected = "f0\tf1\n0\t1\n5\t6\n10\t11\n15\t16\n20\t21\n"
		testOutput(t, input, "\t", "1,2", expected, extractFields)

		input = `Song title,Artist,Year,Progression,Recorded Key
"10000 Reasons (Bless the Lord)",Matt Redman and Jonas Myrin,2012,IV–I–V–vi,G major
"20 Good Reasons",Thirsty Merc,2007,I–V–vi–IV,D♭ major
"Adore You",Harry Styles,2019,vi−I−IV−V,C minor
"Africa",Toto,1982,vi−IV–I–V (chorus),F♯ minor (chorus)
`

		expected = "Song title\n\"10000 Reasons (Bless the Lord)\"\n\"20 Good Reasons\"\n\"Adore You\"\n\"Africa\"\n"

		testOutput(t, input, ",", "1", expected, extractFields)

		expected = "Song title,Artist\n\"10000 Reasons (Bless the Lord)\",Matt Redman and Jonas Myrin\n\"20 Good Reasons\",Thirsty Merc\n\"Adore You\",Harry Styles\n\"Africa\",Toto\n"

		testOutput(t, input, ",", "1,2", expected, extractFields)

		expected = "Song title,Artist,Year,Progression,Recorded Key\n\"10000 Reasons (Bless the Lord)\",Matt Redman and Jonas Myrin,2012,IV–I–V–vi,G major\n\"20 Good Reasons\",Thirsty Merc,2007,I–V–vi–IV,D♭ major\n\"Adore You\",Harry Styles,2019,vi−I−IV−V,C minor\n\"Africa\",Toto,1982,vi−IV–I–V (chorus),F♯ minor (chorus)\n"
		testOutput(t, input, ",", "2,2,3,1-", expected, extractFields)

	})
}

func TestExtractBytes(t *testing.T) {
	t.Run("should return the correct bytes", func(t *testing.T) {
		defaultDelimiter := "\t"
		input := `f0	f1	f2	f3	f4
0	1	2	3	4
5	6	7	8	9
10	11	12	13	14
15	16	17	18	19
20	21	22	23	24`

		expected := "0f1\n\t\t2\n\t\t7\n011\n516\n021\n"

		testOutput(t, input, defaultDelimiter, "2,4-5", expected, extractBytes)

		expected = "f0\tf1\tf2\tf3\tf4\n0\t1\t2\t3\t4\n5\t6\t7\t8\t9\n10\t11\t12\t13\t14\n15\t16\t17\t18\t19\n20\t21\t22\t23\t24\n"

		testOutput(t, input, defaultDelimiter, "1-", expected, extractBytes)

	})
}

func testOutput(t *testing.T, input string, delimiter string, args string, expected string, worker func(line string, delimiter string, list *internal.List) (string, error)) {
	list, err := internal.ParseList(args)

	if err != nil {
		t.Error(err)
	}

	output, err := traverseFileByLine(bufio.NewScanner(strings.NewReader(input)), delimiter, list, worker)

	if err != nil {
		t.Error(err)
	}

	if expected != output {
		t.Errorf("expected %#v, got %#v", expected, output)
	}
}
