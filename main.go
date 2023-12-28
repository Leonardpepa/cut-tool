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

	validateFlags(f, b, c)

	var list List
	var worker func(line string, delimiter string, list List)

	delimiter := *d

	//  use the first character
	if delimiter != "\t" {
		delimiter = string(delimiter[0])
	}

	if *f != "" {
		list = parseList(*f)
		worker = fieldsWorker
	}

	if *b != "" {
		list = parseList(*b)
		worker = bytesWorker
	}

	// same as bytes
	// doesn't support multibyte chars for now
	if *c != "" {
		list = parseList(*c)
		worker = bytesWorker
	}

	run(delimiter, list, worker)
}

func run(delimiter string, list List, worker func(line string, delimiter string, list List)) {
	filenames := flag.Args()

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		traverseFileByLine(file, delimiter, list, worker)

		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

func validateFlags(f *string, b *string, c *string) {
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
}

func fieldsWorker(line string, delimiter string, list List) {
	fields := strings.Split(line, string(delimiter[0]))
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

func bytesWorker(line string, _ string, list List) {
	reader := strings.NewReader(line)

	for _, from := range list.SortedKeys() {
		to := list.numbers[from]
		if to == -1 {
			to = int(reader.Size())
		}
		for i := from; i <= to; i++ {
			b := make([]byte, 1)
			_, err := reader.ReadAt(b, int64(i-1))
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%s", string(b))
		}
	}
}

func traverseFileByLine(file *os.File, delimiter string, list List, work func(line string, delimiter string, list List)) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		work(line, delimiter, list)
		fmt.Println()
	}

}
