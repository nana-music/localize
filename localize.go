package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Item struct {
	File         string
	Key          string
	Comment      string
	Localization map[string]string
}

func NewItem(filename string) *Item {
	return &Item{
		File:         filename,
		Localization: map[string]string{},
	}
}

func (i Item) String() string {
	return strings.Join([]string{
		i.File,
		i.Key,
		i.Comment,
		"ja = " + i.Localization["ja"],
		"en = " + i.Localization["en"],
	}, "\n")
}

func MergeItems(dst map[string]Item, items []Item) {
	for _, i := range items {
		item, ok := dst[i.Key]
		if ok {
			for k, v := range i.Localization {
				item.Localization[k] = v
			}
		} else {
			dst[i.Key] = i
		}
	}
}

func main() {
	find := flag.Bool("find", false, "find .strings files")
	prints := flag.Bool("print", false, "print .strings files")
	csv := flag.Bool("csv", false, "convert .strings to .csv")
	strings := flag.Bool("strings", false, "convert .csv to .strings")
	outputdir := flag.String("o", ".", ".strings output root directory")
	flag.Parse()

	if *find {
		f, err := FindStrings(".")
		if err != nil {
			log.Fatal(err)
		}

		for _, path := range f {
			fmt.Println(path)
		}
		return
	}

	if *prints {
		f, err := FindStrings(".")
		if err != nil {
			log.Fatal(err)
		}

		for _, path := range f {
			items, err := LoadStrings(path)
			if err != nil {
				log.Fatal(err)
			}

			for _, item := range items {
				fmt.Println(item)
			}
			fmt.Println()
		}
		return
	}

	if *csv {
		dst := map[string]Item{}

		f, err := FindStrings(".")
		if err != nil {
			log.Fatal(err)
		}

		for _, path := range f {
			items, err := LoadStrings(path)
			if err != nil {
				log.Fatal(err)
			}
			MergeItems(dst, items)
		}

		items := make([]Item, len(dst))
		i := 0
		for _, item := range dst {
			items[i] = item
			i++
		}

		var out io.Writer
		if len(flag.Args()) == 0 {
			out = os.Stdout
		} else {
			filename := flag.Arg(0)
			file, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			out = file
		}

		err = WriteCsv(out, items)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *strings {
		var in io.Reader
		if len(flag.Args()) == 0 {
			in = os.Stdin
		} else {
			filename := flag.Arg(0)
			file, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			in = file
		}

		items, err := LoadCsv(in)
		if err != nil {
			log.Fatal(err)
		}

		err = WriteStrings(*outputdir, items)
		if err != nil {
			log.Fatal(err)
		}

		return
	}
}
