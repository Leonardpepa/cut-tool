package main

import (
	"bufio"
	"cut-tool/internal"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	toManyListArguments = errors.New(" only one type of list may be specified")
	delimiterError      = errors.New("an input delimiter may be specified only when operating on fields")
	noFlagSpecified     = errors.New("no flag specified")
)

var Empty = ""

func main() {

	var delimiter string
	var help bool
	var fields string
	var bytesFlag string
	var chars string

	flag.StringVar(&fields, "f", Empty, "fields_list")
	flag.StringVar(&fields, "fields", Empty, "fields_list")

	flag.StringVar(&bytesFlag, "b", Empty, "bytes_list")
	flag.StringVar(&bytesFlag, "bytes", Empty, "bytes_list")

	flag.StringVar(&chars, "c", Empty, "characters_list")
	flag.StringVar(&chars, "characters", Empty, "characters_list")

	flag.StringVar(&delimiter, "d", "\t", "delimiter")
	flag.StringVar(&delimiter, "delimiter", "\t", "delimiter")

	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&help, "help", false, "help")

	flag.Parse()

	if help {
		usage()
		os.Exit(0)
	}

	err := validateFlags(fields, bytesFlag, chars, &delimiter)

	if err != nil {
		log.Fatal(err)
	}

	var list *internal.List
	var worker func(line string, delimiter string, list *internal.List) (string, error)

	if fields != Empty {
		list, err = internal.ParseList(fields)
		worker = extractFields
	}

	if bytesFlag != Empty {
		list, err = internal.ParseList(bytesFlag)
		worker = extractBytes
	}

	// same as bytes
	// doesn't support multibyte chars for now
	if chars != Empty {
		list, err = internal.ParseList(chars)
		worker = extractBytes
	}

	if err != nil {
		log.Fatal(err)
	}

	run(delimiter, list, worker)
}

func run(delimiter string, list *internal.List, worker func(line string, delimiter string, list *internal.List) (string, error)) {
	filenames := flag.Args()

	if len(filenames) == 0 || (len(filenames) == 1 && filenames[0] == "-") {
		output, err := traverseFileByLine(bufio.NewScanner(os.Stdin), delimiter, list, worker)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(output)

		return
	}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		output, err := traverseFileByLine(bufio.NewScanner(file), delimiter, list, worker)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(output)

		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

func validateFlags(f, b, c string, d *string) error {

	if f == Empty && b == Empty && c == Empty {
		return noFlagSpecified
	}

	if f != Empty {
		if b != Empty || c != Empty {
			return toManyListArguments
		}
	}

	if b != Empty {
		if f != Empty || c != Empty {
			return toManyListArguments
		}
	}

	if c != Empty {
		if b != Empty || f != Empty {
			return toManyListArguments
		}
	}

	if *d != "\t" && f == Empty {
		return delimiterError
	}

	// use the first char as delimiter
	if *d != "\t" {
		*d = string(string(*d)[0])
	}

	return nil
}

func extractFields(line string, delimiter string, list *internal.List) (string, error) {
	fields := strings.Split(line, delimiter)

	var builder strings.Builder

	for index, from := range list.SortedKeys() {
		to := list.Range(from)
		if to == internal.EndOfTheList || to > len(fields) {
			to = len(fields)
		}
		// don't print the delimiter in the end
		for i := from; i <= to; i++ {
			if index == len(list.SortedKeys())-1 && i == to {
				delimiter = Empty
			}
			builder.WriteString(fmt.Sprintf("%s%s", fields[i-1], delimiter))
		}
	}

	return builder.String(), nil
}

func extractBytes(line string, _ string, list *internal.List) (string, error) {
	reader := strings.NewReader(line)
	var builder strings.Builder
	for _, from := range list.SortedKeys() {
		to := list.Range(from)
		if to == internal.EndOfTheList || to > int(reader.Size()) {
			to = int(reader.Size())
		}
		for i := from; i <= to; i++ {
			b := make([]byte, 1)
			_, err := reader.ReadAt(b, int64(i-1))
			if err != nil {
				return Empty, err
			}

			builder.WriteString(fmt.Sprintf("%s", string(b)))
		}
	}

	return builder.String(), nil
}

func traverseFileByLine(scanner *bufio.Scanner, delimiter string, list *internal.List, work func(line string, delimiter string, list *internal.List) (string, error)) (string, error) {
	scanner.Split(bufio.ScanLines)
	var builder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		s, err := work(line, delimiter, list)
		if err != nil {
			return Empty, err
		}
		builder.WriteString(s + "\n")
	}

	return builder.String(), nil
}

func usage() {
	fmt.Printf(`Usage: %s OPTION... [FILE]...
Print selected parts of lines from each FILE to standard output.

With no FILE, or when FILE is -, read standard input.

Mandatory arguments to long options are mandatory for short options too.
  -b, --bytes=LIST        select only these bytes
  -c, --characters=LIST   select only these characters
  -d, --delimiter=DELIM   use DELIM instead of TAB for field delimiter
  -f, --fields=LIST       select only these fields;  also print any line

Use one, and only one of -b, -c or -f.  Each LIST is made up of one
range, or many ranges separated by commas.  Selected input is written
in the same order that it is read, and is written exactly once.
Each range is one of:

  N     N'th byte, character or field, counted from 1
  N-    from N'th byte, character or field, to end of line
  N-M   from N'th to M'th (included) byte, character or field
  -M    from first to M'th (included) byte, character or field`, filepath.Base(os.Args[0]))
}
