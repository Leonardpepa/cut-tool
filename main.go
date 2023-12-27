package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	toManyListArguments = errors.New("error more tha tow numbers in range")
	decreasingRage      = errors.New("invalid decreasing range")
	invalidRangeFormat  = errors.New("invalid range format")
)

type List struct {
	// set of numbers
	numbers    map[int]int
	sortedKeys []int
}

func (list *List) SortKeys() {
	keys := make([]int, 0)
	for start := range list.numbers {
		keys = append(keys, start)
	}

	sort.Sort(sort.IntSlice(keys))
	list.sortedKeys = keys
}

func (list *List) SortedKeys() []int {
	return list.sortedKeys
}

func (list *List) appendNumber(from int, to int) {
	for start, end := range list.numbers {
		if start > from && (end < to || to == -1) {
			delete(list.numbers, start)
			list.numbers[from] = to
		}
	}
	list.numbers[from] = to
}

func parseList(data string) List {
	list := List{
		numbers:    make(map[int]int),
		sortedKeys: make([]int, 0),
	}

	data = strings.TrimFunc(data, func(r rune) bool {
		return r == '"'
	})

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
	list.SortKeys()
	return list
}

func parseRange(data string) (start, end int) {

	values := strings.Split(data, "-")

	if len(values) != 2 {
		log.Fatal(invalidRangeFormat)
	}

	start = parseEmptyNumber(values[0], 1)
	end = parseEmptyNumber(values[1], -1)

	if end != -1 && start > end {
		log.Fatal(decreasingRage)
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
	b := flag.String("b", "", "bytes_list")
	c := flag.String("c", "", "characters_list")
	d := flag.String("d", "\t", "delimiter")

	flag.Parse()

	if *f != "" {
		if *b != "" || *c != "" {
			log.Fatal(toManyListArguments)
		}
	}
	if *b != "" {
		if *f != "" || *c != "" {
			log.Fatal(toManyListArguments)
		}
	}

	if *c != "" {
		if *b != "" || *f != "" {
			log.Fatal(toManyListArguments)
		}
	}

	fieldsList := parseList(*f)

	filename := flag.Args()[0]

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	traverseFileByLine(file, *d, fieldsList, fieldsWorker)

}

func fieldsWorker(fields []string, list List) {
	for _, from := range list.SortedKeys() {
		to := list.numbers[from]
		if to == -1 {
			to = len(fields)
		}
		for i := from; i <= to; i++ {
			fmt.Printf("%s ", fields[i-1])
		}
	}
}

func traverseFileByLine(file *os.File, delimiter string, list List, work func(fields []string, list List)) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, delimiter)
		work(fields, list)
		fmt.Println()
	}

}
