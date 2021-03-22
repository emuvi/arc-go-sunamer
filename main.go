package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func removeAccents(s string) string {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(transformer, s)
	if err != nil {
		panic(err)
	}
	return result
}

func replaceSpaces(s string) string {
	result := strings.TrimSpace(s)
	result = strings.ReplaceAll(result, " ", "_")
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return result
}

var helper bool
var subdirs bool
var input string
var accents bool
var spaces bool

func process(origin string) {
	info, err := os.Stat(origin)
	if err != nil {
		panic(err)
	}
	if info.IsDir() && subdirs {
		fmt.Println("Folder: " + origin)
		files, err := ioutil.ReadDir(origin)
		if err != nil {
			panic(err)
		}
		for _, inside := range files {
			process(path.Join(origin, inside.Name()))
		}
	} else {
		fmt.Println("File: " + origin)
		oldName := info.Name()
		newName := oldName
		if accents {
			newName = removeAccents(newName)
		}
		if spaces {
			newName = replaceSpaces(newName)
		}
		if newName == oldName {
			fmt.Println("Nothing to do.")
		} else {
			root := path.Dir(origin)
			destiny := path.Join(root, newName)
			index := 1
			for {
				if _, err := os.Stat(destiny); os.IsNotExist(err) {
					break
				} else {
					index++
					extension := filepath.Ext(newName)
					justName := newName[0 : len(newName)-len(extension)]
					destiny = path.Join(root, justName+" ("+strconv.Itoa(index)+")"+extension)
				}
			}
			fmt.Println(destiny)
			err := os.Rename(origin, destiny)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func main() {
	flag.BoolVar(&helper, "h", false, "Show the help.")
	flag.BoolVar(&subdirs, "s", false, "Process the subdirs also.")
	flag.StringVar(&input, "i", ".", "File or folder to be processed.")
	flag.BoolVar(&accents, "a", false, "Remove accents from name.")
	flag.BoolVar(&spaces, "s", false, "Replace spaces to underscores.")
	flag.Parse()
	if helper {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	process(input)
}
