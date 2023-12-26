package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	//filename := "tests/sample.tsv"
	filename := "tests/sample.tsv"

	scanner := getScanner(openFile(filename))

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		r := strings.NewReader(line)
		b1 := make([]byte, 1)
		b2 := make([]byte, 1)
		b5 := make([]byte, 1)

		_, err := r.ReadAt(b1, 0)
		if err != nil {
			log.Fatal(err)
		}

		_, err = r.ReadAt(b2, 1)
		if err != nil {
			log.Fatal(err)
		}

		_, err = r.ReadAt(b5, 4)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s%s%s\n", string(b1), string(b2), string(b5))
	}
}

func getScanner(file *os.File) *bufio.Scanner {
	return bufio.NewScanner(file)
}
func openFile(filename string) *os.File {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalln(err)
	}

	return file
}
