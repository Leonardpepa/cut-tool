package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type List struct {
	// set of numbers
	numbers map[int]int
}

func (list *List) appendNumber(from int, to int) {
	list.numbers[from] = to
}

func (list *List) appendListOfNumber(nums []int) {
	for _, num := range nums {
		list.appendNumber(num, num)
	}
}

func parseList(data string) List {
	list := List{
		numbers: make(map[int]int),
	}

	values := strings.Split(data, ",")

	for _, val := range values {

		isRange := strings.Contains(val, "-")

		if isRange {
			start, end := parseRange(val)
			list.appendNumber(start, end)
			continue
		}

		num, err := strconv.Atoi(val)

		if err == nil {
			list.appendNumber(num, num)
		}
	}

	return list
}

func parseRange(data string) (start, end int) {

	values := strings.Split(data, "-")

	if len(values) != 2 {
		log.Fatal("error more tha tow numbers in range")
	}

	start = parseEmptyNumber(values[0], 1)
	end = parseEmptyNumber(values[1], -1)

	if end != -1 && start > end {
		log.Fatal("invalid decreasing range")
	}

	return
}

func parseEmptyNumber(value string, defaultValue int) int {
	if value != "" {
		num, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}

		return num
	}
	return defaultValue
}

func main() {

	f := flag.String("f", "", "fields_list")

	flag.Parse()

	if *f == "" {
		log.Fatal("no f provided")
	}

	fmt.Println(parseList(*f))
}
