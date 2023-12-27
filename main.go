package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type List struct {
	// set of numbers
	numbers map[int]int
}

func (list *List) getSortedKeys() []int {
	keys := make([]int, 0)
	for start, _ := range list.numbers {
		keys = append(keys, start)
	}

	sort.Sort(sort.IntSlice(keys))
	return keys
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

	fieldsList := parseList(*f)
	keys := fieldsList.getSortedKeys()

	filename := flag.Args()[0]

	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		for _, from := range keys {
			to := fieldsList.numbers[from]
			if to == -1 {
				to = len(fields)
			}
			for i := from; i <= to; i++ {
				fmt.Printf("%s\t", fields[i-1])
			}
			fmt.Println()
		}
	}

}
