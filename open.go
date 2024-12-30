package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
	size uint
	root *Object
	Info *Object

}

type PdfScanner struct {
	bufio.Scanner
	Objects   []Object
	XrefTable map[uint]Xref
}

type Xref struct {
	id         uint
	byteoffset uint
	generation uint
	free       bool
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
		if isEndObject(line) {
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

func (p *PdfScanner) GetXrefGroup(offset int, size int) {
	for i := 0; i < size; i++ {
		if !p.Scanner.Scan() {
			panic(fmt.Sprintf("Missing line %d to %d", i, size))
		}
		line := p.Scanner.Text()
		if !isXrefEntry(line) {
			panic(fmt.Sprintf("Expected line of type 'offset generation status'. Instead got: %s", line))
		}
		parts := strings.Split(line, " ")
		id := uint(offset + i)
		if !strings.Contains("nf", parts[2]) {
			panic(fmt.Sprintf("Expect status of type 'n' or 'f' but, got '%s'", parts[2]))
		}
		p.XrefTable[id] = Xref{
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
		if isXrefGroupHeader(line) {
			offset := getInt(parts[0])
			size := getInt(parts[1])
			p.GetXrefGroup(int(offset), size)
		} else if isTrailer(line) {
			break
		} else {
			println(len(parts), line)
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
		if isStartObject(line) {
			obj := p.ExtractObject(line)
			p.Objects = append(p.Objects, obj)
		}
		if isXref(line) {
			p.ExtractXref(line)
		}
		if isTrailer(line) {
			p.ExtractTrailer(line)
		}
	}
}

func (p *PdfScanner) ExtractTrailer(line string) Trailer {
	for p.Scanner.Scan() {
		line := p.Scanner.Text()
		if isEOF(line) {
			break
		}
	}
}

func isStartObject(line string) bool {
	return verifyRegex(line, `^\d+ \d+ obj$`)
}

func verifyRegex(line string, regex string) bool {
	re, err := regexp.Compile(regex)
	if err != nil {
		return false
	}
	return re.MatchString(line)
}

func isEndObject(line string) bool {
	return verifyRegex(line, `^endobj$`)
}

func isXref(line string) bool {
	return verifyRegex(line, `^xref$`)
}

func isTrailer(line string) bool {
	return verifyRegex(line, `^trailer$`)
}

func isXrefGroupHeader(line string) bool {
	return verifyRegex(line, `^\d+ \d+$`)
}

func isXrefEntry(line string) bool {
	return verifyRegex(line, `^\d+ \d+ [fn]$`)
}

func isEOF(line string) bool {
	return verifyRegex(line, "%%EOF")
}
