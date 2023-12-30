package internal

import (
	"errors"
	"log"
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
	alreadyAdded := false
	// merge ranges if needed
	for start, end := range list.ranges {
		if from-1 <= start && (to > end && end != -1 || to == -1) {
			if from-1 == start || from == start {
				list.ranges[start] = to
			} else {
				delete(list.ranges, start)
				list.ranges[from] = to
			}
			alreadyAdded = true
		} else if from >= start && ((to < end && to != -1) || (to == end)) {
			alreadyAdded = true
		}
	}

	if !alreadyAdded {
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

func parseRange(data string) (start, end int, err error) {

	if data == "-" {
		err = invalidRangeWithNoEndPoint
	}

	values := strings.Split(data, "-")

	if len(values) != 2 {
		log.Fatal(invalidRangeFormat)
	}

	// default 0 means from the beginning of the line
	start, err = parseEmptyNumberInRange(values[0], 1)

	// default -1 means end of the line
	end, err = parseEmptyNumberInRange(values[1], -1)

	if end != -1 && start > end {
		err = decreasingRage
	}

	return
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
