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
	toManyListArguments        = errors.New(" only one type of list may be specified")
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

func main() {

	f := flag.String("f", "", "fields_list")
	b := flag.String("b", "", "bytes_list")
	c := flag.String("c", "", "characters_list")
	d := flag.String("d", "\t", "delimiter")

	flag.Parse()

	err := validateFlags(f, b, c)
	if err != nil {
		log.Fatal(err)
	}

	var list *List
	var worker func(line string, delimiter string, list *List)

	delimiter := *d

	//  use the first character
	if delimiter != "\t" {
		delimiter = string(delimiter[0])
	}

	if *f != "" {
		list, err = parseList(*f)
		worker = fieldsWorker
	}

	if *b != "" {
		list, err = parseList(*b)
		worker = bytesWorker
	}

	// same as bytes
	// doesn't support multibyte chars for now
	if *c != "" {
		list, err = parseList(*c)
		worker = bytesWorker
	}

	if err != nil {
		log.Fatal(err)
	}

	run(delimiter, list, worker)
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
	// merge ranges if needed
	for start, end := range list.ranges {
		if start > from && (end < to || to == -1) {
			delete(list.ranges, start)
			list.ranges[from] = to
		}
	}
	list.ranges[from] = to
}

func parseList(data string) (*List, error) {
	list := &List{
		ranges:     make(map[int]int),
		sortedKeys: make([]int, 0),
	}

	values := prepareTheArguments(data)

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
func prepareTheArguments(data string) []string {
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

func run(delimiter string, list *List, worker func(line string, delimiter string, list *List)) {
	filenames := flag.Args()

	if len(filenames) == 0 || (len(filenames) == 1 && filenames[0] == "-") {
		traverseFileByLine(bufio.NewScanner(os.Stdin), delimiter, list, worker)
		return
	}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		traverseFileByLine(bufio.NewScanner(file), delimiter, list, worker)

		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

func validateFlags(f *string, b *string, c *string) error {
	if *f != "" {
		if *b != "" || *c != "" {
			return toManyListArguments
		}
	}
	if *b != "" {
		if *f != "" || *c != "" {
			return toManyListArguments
		}
	}

	if *c != "" {
		if *b != "" || *f != "" {
			return toManyListArguments
		}
	}
	return nil
}

func fieldsWorker(line string, delimiter string, list *List) {
	fields := strings.Split(line, delimiter)
	for index, from := range list.SortedKeys() {
		to := list.ranges[from]
		if to == -1 || to > len(fields) {
			to = len(fields)
		}
		// dont print the comma in the end
		if index == len(list.SortedKeys())-1 {
			delimiter = ""
		}
		for i := from; i <= to; i++ {
			fmt.Printf("%s%s", fields[i-1], delimiter)
		}
	}
}

func bytesWorker(line string, _ string, list *List) {
	reader := strings.NewReader(line)

	for _, from := range list.SortedKeys() {
		to := list.ranges[from]
		if to == -1 || to > int(reader.Size()) {
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

func traverseFileByLine(scanner *bufio.Scanner, delimiter string, list *List, work func(line string, delimiter string, list *List)) {
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		work(line, delimiter, list)
		fmt.Println()
	}

}
