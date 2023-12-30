package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	toManyListArguments        = errors.New(" only one type of list may be specified")
	decreasingRage             = errors.New("invalid decreasing range")
	invalidRangeFormat         = errors.New("invalid range format")
	invalidNumberFormat        = errors.New("invalid number format")
	invalidRangeWithNoEndPoint = errors.New("invalid range with no endpoint: -")
	delimiterError             = errors.New("an input delimiter may be specified only when operating on fields")
)

func main() {

	var delimiter string

	f := flag.String("f", "", "fields_list")
	b := flag.String("b", "", "bytes_list")
	c := flag.String("c", "", "characters_list")
	flag.StringVar(&delimiter, "d", "\t", "delimiter")

	flag.Parse()

	err := validateFlags(f, b, c, &delimiter)

	if err != nil {
		log.Fatal(err)
	}

	var list *List
	var worker func(line string, delimiter string, list *List) (string, error)

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

func run(delimiter string, list *List, worker func(line string, delimiter string, list *List) (string, error)) {
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

func validateFlags(f, b, c, d *string) error {
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

	if *d != "\t" && *f == "" {
		return delimiterError
	}

	// use the first char as delimiter
	if *d != "\t" {
		*d = string(string(*d)[0])
	}

	return nil
}

func fieldsWorker(line string, delimiter string, list *List) (string, error) {
	fields := strings.Split(line, delimiter)

	var builder strings.Builder

	for index, from := range list.SortedKeys() {
		to := list.ranges[from]
		if to == -1 || to > len(fields) {
			to = len(fields)
		}
		// don't print the comma in the end
		for i := from; i <= to; i++ {
			if index == len(list.SortedKeys())-1 && i == to {
				delimiter = ""
			}
			builder.WriteString(fmt.Sprintf("%s%s", fields[i-1], delimiter))
		}
	}

	return builder.String(), nil
}

func bytesWorker(line string, _ string, list *List) (string, error) {
	reader := strings.NewReader(line)
	var builder strings.Builder
	for _, from := range list.SortedKeys() {
		to := list.ranges[from]
		if to == -1 || to > int(reader.Size()) {
			to = int(reader.Size())
		}
		for i := from; i <= to; i++ {
			b := make([]byte, 1)
			_, err := reader.ReadAt(b, int64(i-1))
			if err != nil {
				return "", err
			}

			builder.WriteString(fmt.Sprintf("%s", string(b)))
		}
	}

	return builder.String(), nil
}

func traverseFileByLine(scanner *bufio.Scanner, delimiter string, list *List, work func(line string, delimiter string, list *List) (string, error)) {
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		s, err := work(line, delimiter, list)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(s)
	}

}
