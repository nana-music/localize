package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/scanner"
)

func FindStrings(root string) ([]string, error) {
	result := []string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".strings" {
			result = append(result, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func LoadStrings(name string) ([]Item, error) {
	lang := ""
	dir := filepath.Dir(name)
	for dir != "." {
		base := filepath.Base(dir)
		if strings.HasSuffix(base, ".lproj") {
			lang = strings.Split(base, ".")[0]
			break
		}
		dir = filepath.Dir(dir)
	}

	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	path := strings.Split(name, string(filepath.Separator))
	for i, _ := range path {
		if filepath.Ext(path[i]) == ".lproj" {
			path[i] = "*.lproj"
		}
	}
	name = filepath.Join(path...)

	result := []Item{}
	item := NewItem(name)
	key := true
	s := new(scanner.Scanner).Init(f)
	s.Mode = scanner.ScanStrings | scanner.ScanComments
	t := s.Scan()
	for t != scanner.EOF {
		text := s.TokenText()

		if strings.HasPrefix(text, "/*") || strings.HasPrefix(text, "//") {
			text = strings.Trim(text, "//")
			text = strings.Trim(text, "/*")
			text = strings.Trim(text, "*/")
			text = strings.TrimSpace(text)
			if len(item.Comment) > 0 {
				item.Comment += "\n" + text
			} else {
				item.Comment = text
			}
		} else if key && strings.HasPrefix(text, "\"") {
			item.Key = strings.Trim(text, "\"")
		} else if text == "=" {
			key = false
		} else if !key && strings.HasPrefix(text, "\"") {
			item.Localization[lang] = strings.Trim(text, "\"")
		} else if text == ";" {
			result = append(result, *item)
			item = NewItem(name)
			key = true
		}

		t = s.Scan()
	}

	return result, nil
}

func WriteStrings(root string, items []Item) error {
	files := map[string]*os.File{}
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	for _, item := range items {
		for lang := range item.Localization {
			paths := strings.Split(item.File, string(filepath.Separator))
			for i := range paths {
				if paths[i] == "*.lproj" {
					paths[i] = lang + ".lproj"
					break
				}
			}
			filename := filepath.Join(paths...)

			f := files[filename]
			if f == nil {
				path := filepath.Join(root, filename)
				err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755))
				if err != nil {
					return err
				}
				f, err = os.Create(path)
				if err != nil {
					return err
				}
				files[filename] = f
			}

			_, err := f.WriteString("/* " + item.Comment + " */\n")
			_, err = f.WriteString("\"" + item.Key + "\" = \"" + item.Localization[lang] + "\";\n\n")

			if err != nil {
				return err
			}
		}
	}
	return nil
}
