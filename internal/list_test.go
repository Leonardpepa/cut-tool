package internal

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseList(t *testing.T) {
	t.Run("should parse the list from args", func(t *testing.T) {

		testList(t, "1,2,3", map[int]int{1: 3})

		testList(t, "1-", map[int]int{1: EndOfTheList})

		testList(t, "3-4,1-1,5-", map[int]int{1: 1, 3: EndOfTheList})

		testList(t, "1,1,1,1-1", map[int]int{1: 1})

		testList(t, "2-4,2-3", map[int]int{2: 4})

		testList(t, "1-2,1-4", map[int]int{1: 4})

		testList(t, "1-2,2-3,4,5", map[int]int{1: 5})

		testList(t, `"1-2,2-3,4,5"`, map[int]int{1: 5})

		testList(t, "2-5,1-1", map[int]int{1: 5})

		testErrorList(t, "3-1", decreasingRage)

		testErrorList(t, "3r", invalidNumberFormat)

		testErrorList(t, "1-  ,   2-", invalidNumberFormat)

		testErrorList(t, "-", invalidRangeWithNoEndPoint)

		testErrorList(t, "1,2,3-,-", invalidRangeWithNoEndPoint)
		testErrorList(t, `"1,2,3-,-"`, invalidRangeWithNoEndPoint)

	})
}

func testList(t *testing.T, arg string, expected map[int]int) {
	list, err := ParseList(arg)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, list.ranges) {
		t.Errorf("expected %v, got %v", expected, list.ranges)
	}
}

func testErrorList(t *testing.T, arg string, expected error) {
	_, err := ParseList(arg)

	if err == nil {
		t.Error("wrong list parsed")
	}

	if !errors.Is(err, expected) {
		t.Errorf("expected %v, got %v", expected, err)
	}
}
