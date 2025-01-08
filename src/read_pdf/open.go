package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	START_OBJECT_REGEX      = `^\d+ \d+ obj$`
	END_OBJECT_REGEX        = `^endobj$`
	START_XREF_REGEX        = `^xref$`
	XREF_GROUP_HEADER_REGEX = `^\d+ \d+$`
	XREF_GROUP_ENTRY_REGEX  = `^\d+ \d+ [fn]$`
	END_XREF_REGEX          = `^trailer$`
	START_TRAILER_REGEX     = `^trailer$`
	END_TRAILER_REGEX       = `%%EOF`
	START_DICTIONARY        = `^\s*<<`
	END_DICTIONARY          = `.*\s*>>`

	TRAILER_REFEX    = `(?:^|\s+)trailer([\s\S]*?)%%EOF`
	XREF_REGEX       = `(?:^|\s+)xref([\s\S]*?)trailer`
	OBJECT_REGREX    = `(\d+) (\d+) obj([\s\S]+?)endobj`
	OUTER_DICTIONARY = `(\/\w+)?\s*<<(.*)>>`
	STREAM           = `stream([\s\S]+?)endstream`
)

type Header struct {
	version string
	size    uint
}

type Object struct {
	ID       uint
	Revision uint
	Data     interface{}
}

type Page struct {
	Number  int
	Content []byte
}

type PDF struct {
	Header
	Objects   []Object
	Pages     []Page
	XrefTable map[uint]Xref
}

type Trailer struct {
	startXref uint
	size      uint
	root      *Object
	Info      *Object
}

type PdfScanner struct {
	bufio.Scanner
	pdf *PDF
}

type Xref struct {
	id         uint
	byteoffset uint
	generation uint
	free       bool
}

type Dictionary struct {
	Entries map[string]interface{}
	Parent  *Dictionary
}

func (p *PdfScanner) ExtractObject(line string) Object {
	found := false
	id, revision := parseStartObject(line)
	obj := Object{
		ID:       id,
		Revision: revision,
	}
	parts := []string{line}
	for p.Scanner.Scan() {
		line := p.Scanner.Text()
		parts = append(parts, line)
		if verifyRegex(line, END_OBJECT_REGEX) {
			found = true
			break
		}
	}
	if !found {
		panic("Was unable to find ending object")
	}
	obj.Data = parts
	return obj
}

func open(location string) {
	file, err := os.Open(location)
	if err != nil {
		panic(err)
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fileContent := string(contents)
	objectRegex := regexp.MustCompile(OBJECT_REGREX)
	objects := objectRegex.FindAllStringSubmatch(fileContent, -1)

	// xrefRegex := regexp.MustCompile(XREF_REGEX)
	// xrefTable := xrefRegex.FindAllString(fileContent, -1)

	// trailerRegex := regexp.MustCompile(TRAILER_REFEX)
	// trailerTable := trailerRegex.FindAllStringSubmatch(fileContent, -1)

	for _, entry := range objects[:1] {
		dictionaryRegex := regexp.MustCompile(OUTER_DICTIONARY)
		text := entry[0]
		fmt.Println(text)
		var key, value string
		objectDict := Dictionary{
			Entries: make(map[string]interface{}),
		}
		currentDict := objectDict
		for dictionary := dictionaryRegex.FindAllStringSubmatch(text, -1); dictionary != nil; text = value {
			elements := dictionary[0]
			if len(elements) == 3 {
				key = elements[1]
				value = elements[2]
				
			} else {
				value = elements[1]
				currentDict.Entries["entries"] = value
			}
			dictionary = dictionaryRegex.FindAllStringSubmatch(text, -1)
			parentDict := currentDict
			currentDict = Dictionary{
				Entries: make(map[string]interface{}),
				Parent: &parentDict,
			}
		}
		fmt.Println(key)
		// //
		// 		object := Object{
		// 			ID:       uint(getInt(entry[1])),
		// 			Revision: uint(getInt(entry[2])),
		// 			Data:     objects[3],
		// 		}
		// 		fmt.Println(object)
		// |(?:<<.*?>>)
		// |(?:\[.*?\])
		// (?m)\/(\w+)((?:\(.*?\)))
		// \/(\w+)\s*((?:\(.*?\))|(?:\[.*?\])|(?:<<.*?>>)|(?:))
	}
}

func (p *PdfScanner) GetXrefGroup(offset int, size int) {
	for i := 0; i < size; i++ {
		if !p.Scanner.Scan() {
			panic(fmt.Sprintf("Missing line %d to %d", i, size))
		}
		line := p.Scanner.Text()
		if !verifyRegex(line, XREF_GROUP_ENTRY_REGEX) {
			panic(fmt.Sprintf("Expected line of type 'offset generation status'. Instead got: %s", line))
		}
		parts := strings.Split(line, " ")
		id := uint(offset + i)
		if !strings.Contains("nf", parts[2]) {
			panic(fmt.Sprintf("Expect status of type 'n' or 'f' but, got '%s'", parts[2]))
		}
		p.pdf.XrefTable[id] = Xref{
			id:         uint(id),
			byteoffset: uint(getInt(parts[0])),
			generation: uint(getInt(parts[1])),
			free:       parts[2] == "f",
		}
	}
}

func (p *PdfScanner) ExtractXref(line string) {
	for p.Scanner.Scan() {
		line := p.Scanner.Text()
		parts := strings.Split(line, " ")
		if verifyRegex(line, XREF_GROUP_HEADER_REGEX) {
			offset := getInt(parts[0])
			size := getInt(parts[1])
			p.GetXrefGroup(int(offset), size)
		} else if verifyRegex(line, END_XREF_REGEX) {
			break
		} else {
			panic("Xref table is corrupted. It expected line of type '<offset> <size>' but got: " + line)
		}
	}
}

func getInt(value string) int {
	if val, err := strconv.ParseInt(value, 10, 0); err != nil {
		panic(fmt.Sprintf("Expect integer, got '%s'", value))
	} else {
		return int(val)
	}
}

func parseStartObject(line string) (uint, uint) {
	parts := strings.Split(line, " ")
	id, err := strconv.ParseUint(parts[0], 10, 0)
	if err != nil {
		panic(err)
	}
	revision, err := strconv.ParseUint(parts[1], 10, 0)
	if err != nil {
		panic(err)
	}
	return uint(id), uint(revision)
}

func (p *PdfScanner) Scan() {
	for p.Scanner.Scan() {
		line := strings.TrimSpace(p.Text())
		if verifyRegex(line, START_OBJECT_REGEX) {
			obj := p.ExtractObject(line)
			p.pdf.Objects = append(p.pdf.Objects, obj)
		}
		if verifyRegex(line, START_XREF_REGEX) {
			p.ExtractXref(line)
		}
		if verifyRegex(line, START_TRAILER_REGEX) {
			p.ExtractTrailer(line)
		}
	}
}

func (p *PdfScanner) ExtractTrailer(line string) Trailer {
	for p.Scanner.Scan() {
		line := p.Scanner.Text()
		if verifyRegex(line, END_TRAILER_REGEX) {
			break
		}
		if verifyRegex(line, `^<<`) {
			fmt.Println(line)
		}
	}
	return Trailer{
		startXref: uint(0),
	}
}

func verifyRegex(line string, regex string) bool {
	re, err := regexp.Compile(regex)
	if err != nil {
		return false
	}
	return re.MatchString(line)
}
