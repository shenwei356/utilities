package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/shenwei356/util/bytesize"
	"github.com/shenwei356/util/sortitem"
)

var usage = `    
Usage: SortYumPackageSize datafile [datafile...]

By Wei Shen < shenwei356@gmail.com http://shenwei.me >

Update: 2014-3-13

Contents in data file are copied when running yum update:

    analitza      x86_64   4.12.3-1.fc20     updates       312 k
    ark           x86_64   4.12.3-1.fc20     updates       278 k
    ark-libs      x86_64   4.12.3-1.fc20     updates       138 k

`
var BytesizeRegexp = regexp.MustCompile(`(?i)([\d\.]+\s*(?:[KMGTPEZY]?B|[BKMGTPEZY]))\s*$`)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(0)
	}

	for _, file := range os.Args[1:] {
		err := HandleFile(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

}

func HandleFile(file string) error {
	fh, err := os.Open(file)
	if err != nil {
		return errors.New("Failed to open file: " + file)
	}
	defer func() {
		fh.Close()
	}()

	reader := bufio.NewReader(fh)
	data := make([]sortitem.Item, 0)
	var sum float64 = 0

	var line string
	for {
		line, err = reader.ReadString('\n')
		line = strings.TrimRight(line, "\r?\n")
		if err == io.EOF {
			break
		}

		if !BytesizeRegexp.Match([]byte(line)) {
			continue
		}

		subs := BytesizeRegexp.FindSubmatch([]byte(line))
		size, err := bytesize.Parse(subs[1])
		if err != nil {
			return errors.New("Failed to parse bytesize: " + string(subs[1]))
		}

		sum += float64(size)

		data = append(data, sortitem.Item{line, float64(size)})

	}

	sort.Sort(sortitem.Reverse{sortitem.ByValue{data}})

	for _, item := range data {
		fmt.Printf("%v\t%s\n", bytesize.ByteSize(item.Value), item.Key)
	}

	fmt.Printf("\nSum: %s\n", bytesize.ByteSize(sum))
	return nil
}
