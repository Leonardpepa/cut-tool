# cut tool

## Purpose
This project is a solution for [Write Your Own cut Tool](https://codingchallenges.fyi/challenges/challenge-cut)
build for my personal educational purposes

## Description
cut is a command line tool, read the [original specification](https://www.gnu.org/software/coreutils/manual/html_node/cut-invocation.html#cut-invocation) for more

## Usage
```terminal
Usage: cut-tool.exe OPTION... [FILE]...
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
  -M    from first to M'th (included) byte, character or field
  ```

## How to run
1. Clone the repo ```git clone https://github.com/Leonardpepa/cut-tool.git```
2. Build ```go build```
3. run on windows```cut-tool.exe [OPTIONS] [FILE]```
4. run on linux ```./cut-tool [OPTIONS] [FILE]```
5. run tests ```go test ./...```