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

var helper bool
var subdirs bool
var remove_accents bool
var spaces_to_underscores bool

func sunamer(origin string) {
	info, err := os.Stat(origin)
	if err != nil {
		log.Println(err)
		return
	}
	if info.IsDir() && subdirs {
		fmt.Println("Folder: " + origin)
		files, err := ioutil.ReadDir(origin)
		if err != nil {
			log.Println(err)
			return
		}
		for _, inside := range files {
			sunamer(path.Join(origin, inside.Name()))
		}
	} else {
		fmt.Println("File: " + origin)
		oldName := info.Name()
		newName := oldName
		if remove_accents {
			newName = doRemoveAccents(newName)
		}
		if spaces_to_underscores {
			newName = doSpacesToUnderscores(newName)
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
					destiny = path.Join(root, justName+"("+strconv.Itoa(index)+")"+extension)
				}
			}
			fmt.Println(destiny)
			err := os.Rename(origin, destiny)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func main() {
	flag.BoolVar(&helper, "h", false, "Show the usage help message.")
	flag.BoolVar(&subdirs, "d", false, "Sunamer in the subdirectories.")
	flag.BoolVar(&remove_accents, "ra", false, "Remove accents from name.")
	flag.BoolVar(&spaces_to_underscores, "su", false, "Replace spaces to underscores.")
	flag.Parse()
	if helper {
		fmt.Printf("Usage of sunamer [OPTION]... [INPUT]...\n")
		flag.PrintDefaults()
		os.Exit(0)
	}
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("You must pass at least one input file.")
		os.Exit(-1)
	}
	for _, input := range args {
		sunamer(input)
	}
}

func doRemoveAccents(s string) string {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(transformer, s)
	if err != nil {
		log.Println(err)
		return s
	}
	return result
}

func doSpacesToUnderscores(s string) string {
	result := strings.TrimSpace(s)
	result = strings.ReplaceAll(result, " ", "_")
	for strings.Contains(result, "__") {
		result = strings.ReplaceAll(result, "__", "_")
	}
	return result
}
