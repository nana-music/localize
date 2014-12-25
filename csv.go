package main

import (
	"encoding/csv"
	"io"
	"sort"
)

type SortableItems []Item

func (s SortableItems) Len() int {
	return len(s)
}

func (s SortableItems) Less(i, j int) bool {
	if s[i].File < s[j].File {
		return true
	} else if s[i].File == s[j].File {
		return s[i].Key < s[j].Key
	} else {
		return false
	}
}

func (s SortableItems) Swap(i, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

func WriteCsv(w io.Writer, items []Item) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	//header
	writer.Write([]string{
		"File",
		"Key",
		"Comment",
		"ja",
		"en",
	})

	sort.Sort(SortableItems(items))

	for _, item := range items {
		err := writer.Write([]string{
			item.File,
			item.Key,
			item.Comment,
			item.Localization["ja"],
			item.Localization["en"],
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func LoadCsv(r io.Reader) ([]Item, error) {
	result := []Item{}
	reader := csv.NewReader(r)

	//skip header
	reader.Read()

	record, err := reader.Read()
	for err == nil {
		i := NewItem(record[0])
		i.Key = record[1]
		i.Comment = record[2]
		i.Localization["ja"] = record[3]
		i.Localization["en"] = record[4]

		result = append(result, *i)
		record, err = reader.Read()
	}

	if err != nil && err != io.EOF {
		return nil, err
	}
	return result, nil
}
