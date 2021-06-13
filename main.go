package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var subdirs bool
var trim_spaces bool
var many_spaces_to_singles bool
var spaces_to_underscores bool
var remove_accents bool
var helper bool

type Option struct {
	call string
	help string
	plug *bool
}

var options = []Option{
	{
		call: "sd",
		help: "Sunamer in the subdirectories.",
		plug: &subdirs,
	},
	{
		call: "ts",
		help: "Trim spaces from begin and end.",
		plug: &trim_spaces,
	},
	{
		call: "ss",
		help: "Replace many spaces to singles.",
		plug: &many_spaces_to_singles,
	},
	{
		call: "su",
		help: "Replace spaces to underscores.",
		plug: &spaces_to_underscores,
	},
	{
		call: "ra",
		help: "Remove accents from name.",
		plug: &remove_accents,
	},
	{
		call: "h",
		help: "Show the usage help message.",
		plug: &helper,
	},
}

func main() {
	for _, option := range options {
		flag.BoolVar(option.plug, option.call, false, option.help)
	}
	flag.Usage = func() {
		flagSet := flag.CommandLine
		for _, option := range options {
			flag := flagSet.Lookup(option.call)
			fmt.Printf("-%s\n", flag.Name)
			fmt.Printf("  %s\n", flag.Usage)
		}
	}
	flag.Parse()
	if helper {
		fmt.Println("Usage of sunamer [OPTION]... [INPUT]...")
		fmt.Println("Options:")
		flag.Usage()
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

func sunamer(origin string) {
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
			sunamer(path.Join(origin, inside.Name()))
		}
	} else {
		fmt.Println("File: " + origin)
		oldName := info.Name()
		newName := oldName
		if trim_spaces {
			newName = doTrimSpaces(newName)
		}
		if many_spaces_to_singles {
			newName = doManySpacesToSingles(newName)
		}
		if spaces_to_underscores {
			newName = doSpacesToUnderscores(newName)
		}
		if remove_accents {
			newName = doRemoveAccents(newName)
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
				panic(err)
			}
		}
	}
}

func doTrimSpaces(name string) string {
	return strings.TrimSpace(name)
}

func doManySpacesToSingles(name string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(name, " ")
}

func doSpacesToUnderscores(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}

func doRemoveAccents(name string) string {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(transformer, name)
	if err != nil {
		panic(err)
	}
	return result
}
