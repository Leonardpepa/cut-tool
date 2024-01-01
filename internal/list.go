package internal

import (
	"errors"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

var (
	decreasingRage             = errors.New("invalid decreasing range")
	invalidRangeFormat         = errors.New("invalid range format")
	invalidNumberFormat        = errors.New("invalid number format")
	invalidRangeWithNoEndPoint = errors.New("invalid range with no endpoint: -")
)

const EndOfTheList = math.MaxInt32

type List struct {
	// set of ranges
	ranges     map[int]int
	sortedKeys []int
}

func (list *List) Range(start int) int {
	return list.ranges[start]
}

func (list *List) SortKeys() {
	keys := make([]int, 0)
	for start := range list.ranges {
		keys = append(keys, start)
	}

	sort.Sort(sort.IntSlice(keys))
	list.sortedKeys = keys
}

func (list *List) SortedKeys() []int {
	return list.sortedKeys
}

func (list *List) appendNumber(from int, to int) {
	hasBeenAdded := false

	for rangeStart, rangeStop := range list.ranges {

		// bigger range already exists
		if rangeStart <= from && rangeStop >= to {
			// nothing
			hasBeenAdded = true
			continue
		}

		// case bigger range
		// merge with the old one
		if from <= rangeStart && to > rangeStop {
			if from != rangeStart {
				delete(list.ranges, rangeStart)
			}
			list.ranges[from] = to
			hasBeenAdded = true
			continue
		}

		//the continuation of a list
		if (from == rangeStop || from-1 == rangeStop) && to >= rangeStop {
			list.ranges[rangeStart] = to
			hasBeenAdded = true
			continue
		}

		// merge from beginning
		if from == rangeStart-1 {
			delete(list.ranges, rangeStart)
			if to < rangeStop {
				list.ranges[from] = rangeStop
			} else {
				list.ranges[from] = to
			}
			hasBeenAdded = true
			continue
		}

		if from-1 == rangeStart && to > rangeStop {
			list.ranges[rangeStart] = to
			hasBeenAdded = true
			continue
		}
	}
	if !hasBeenAdded {
		list.ranges[from] = to
	}

}

func ParseList(data string) (*List, error) {
	list := &List{
		ranges:     make(map[int]int),
		sortedKeys: make([]int, 0),
	}

	values := prepareListArguments(data)

	for _, val := range values {

		isRange := strings.Contains(val, "-")

		if isRange {
			start, end, err := parseRange(val)

			if err != nil {
				return nil, err
			}
			list.appendNumber(start, end)
			continue
		}

		num, err := strconv.Atoi(val)

		if err != nil {
			return nil, invalidNumberFormat
		}

		list.appendNumber(num, num)
	}

	list.SortKeys()

	return list, nil
}

// split the list items
// ranges, ranges with -, ranges with space
func prepareListArguments(data string) []string {
	args := make([]string, 0)
	data = strings.TrimFunc(data, func(r rune) bool {
		return r == '"'
	})

	values := strings.Split(data, ",")

	for _, val := range values {
		args = append(args, strings.Split(val, " ")...)
	}

	return args
}

func parseRange(data string) (int, int, error) {
	if data == "-" {
		return 0, 0, invalidRangeWithNoEndPoint
	}
	values := strings.Split(data, "-")

	if len(values) != 2 {
		log.Fatal(invalidRangeFormat)
	}

	// default 0 means from the beginning of the line
	start, err := parseEmptyNumberInRange(values[0], 1)

	if err != nil {
		return 0, 0, err
	}

	// default -1 means end of the line
	end, err := parseEmptyNumberInRange(values[1], EndOfTheList)

	if err != nil {
		return 0, 0, err
	}

	if end != -1 && start > end {
		return 0, 0, decreasingRage
	}

	return start, end, nil
}

func parseEmptyNumberInRange(value string, defaultValue int) (int, error) {
	if value != "" {
		num, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}

		return num, nil
	}
	return defaultValue, nil
}
