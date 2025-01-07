package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	ALPHANUMERIC = "[a-zA-Z0-9]"
	SYMBOL       = "[.,_\\-\\!\"'#$%&*+]"
)

const (
	// Continue
	continueScan int = iota
	beginString
	str
	endString
	beginDictionary
	endDictionary
	beginKey
	key
	endKey
	beginValue
	value
	endValue
	beginArray
	endArray
	beginHex
	hex
	endHex

	// Stop
	error
	eof
)

func clarify(signal int) string {
	return []string{
		"continueScan",
		"beginString",
		"str",
		"endString",
		"beginDictionary",
		"endDictionary",
		"beginKey",
		"key",
		"endKey",
		"beginValue",
		"value",
		"endValue",
		"beginArray",
		"endArray",
		"beginHex",
		"hex",
		"endHex",
		"error",
		"eof",
	}[signal]
}

type scanState struct {
	isHex        bool
	isValue      bool
	isString     bool
	isEndingDict bool
	isKey        bool
	hexStarted   bool
	Key          string
	value        string
	reader       io.Reader
	offset       int64
	stateStack   []int
	err          string
	depth        int
}

func (s *scanState) init(reader io.Reader) {
	s.stateStack = make([]int, 0)
	s.depth = 0
	s.reader = reader
	s.offset = 0
	s.isHex = false
	s.hexStarted = false
	s.isKey = false
	s.Key = ""
	s.value = ""
	s.err = ""
}

func export(signal int, char string) {
	fmt.Println(clarify(signal), char)
}

func (s *scanState) scan() {
	byte := make([]byte, 1)
	for _, err := s.reader.Read(byte); err != io.EOF; _, err = s.reader.Read(byte) {
		char := string(byte)
		switch {
		case char == " ":
			if s.isKey {
				s.isKey = false
				export(endKey, char)
			} else if s.isValue {
				export(value, char)
			} else if s.isString {
				export(str, char)
			} else {
				export(continueScan, char)
			}
		case char == "(":
			s.isString = true
			export(beginString, char)
		case char == ")":
			s.isString = false
			export(endString, char)
		case char == "<":
			if s.isHex {
				s.isHex = false
				s.depth += 1
				export(beginDictionary, char)
			} else {
				s.isHex = true
				export(continueScan, char)
			}
		case char == "/":
			if s.isKey {
				export(endKey, char)
			} else if s.isValue {
				s.isValue = false
				export(endValue, char)
			}
			s.isKey = true
			export(beginKey, char)
		case char == ">":
			if s.isValue {
				s.isValue = false
				export(endValue, char)
			} else if s.isKey {
				s.isKey = false
				export(endKey, char)
			}
			if s.isHex {
				s.isHex = false
				export(endHex, char)
			} else if s.isEndingDict && s.depth > 0 {
				s.depth--
				s.isEndingDict = false
				export(endDictionary, char)
			} else if s.isEndingDict {
				s.err = "Cannot exit top-level dictionary"
				export(error, char)
				panic(s.err)
			} else {
				s.isEndingDict = true
			}
		case regexp.MustCompile(ALPHANUMERIC).Match(byte):
			if s.isKey {
				export(key, char)
			} else if s.isHex {
				if !s.hexStarted {
					export(beginHex, char)
					s.hexStarted = true
				}
				export(hex, char)
			} else if s.isString {
				export(str, char)
			} else if s.isValue {
				export(value, char)
			} else {
				s.isValue = true
				export(beginValue, char)
			}
		case regexp.MustCompile(SYMBOL).Match(byte):
			if s.isValue {
				export(value, char)
			}
		}
	}
}

func (s *scanState) step() {
}

func main() {
	s := scanState{}
	s.init(strings.NewReader(`<<
  /Type /Catalog
  /Version 1.7
  /Pages 3 0 R
  /Metadata 4 0 R
  /Info <<
    /Title (Sample PDF)
    /Author (John Doe)
    /CreationDate (D:20250101000000Z)
    /Producer (Acrobat Distiller 2025)
    /Keywords (PDF, Dictionary, Example, Metadata)
  >>
  /Contents 5 0 R
  /Fonts <<
    /F1 <<
      /Type /Font
      /Subtype /Type1
      /BaseFont /Times-Roman
      /Encoding /WinAnsiEncoding
    >>
  >>
  /Resources <<
    /ProcSet [ /PDF /Text /ImageB /ImageC ]
    /XObject <<
      /Im1 6 0 R
      /Im2 7 0 R
    >>
  >>
  /Annots [
    <<
      /Type /Annot
      /Subtype /Text
      /Contents (This is a sample annotation)
      /Rect [ 100 200 200 300 ]
      /Color [ 0.0 1.0 0.0 ]
      /Open true
    >>
    <<
      /Type /Annot
      /Subtype /Link
      /Rect [ 150 250 250 350 ]
      /A <<
        /Type /Action
        /S /URI
        /URI (http://example.com)
      >>
    >>
  ]
  /Metadata <<
    /Type /Metadata
    /ContentLength 350
    /Content << /HexData <FEEDFACE> >>
  >>
  /CustomData <<
    /Version 1
    /Options <<
      /EnableFeatureX true
      /Threshold 0.75
      /Data [ 12 34 56 78 90 ]
      /HexValue <9A8F4D1E>
    >>
    /Notes (This section contains custom data)
  >>
  /StructureTreeRoot 8 0 R
>>`))
	s.scan()
}
