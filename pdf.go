package main

import "io"

const (
	// Continue
	continue int = iota
	beginString
	beginDictionary
	endDictionary
	beginKey
	endKey
	beginValue
	endValue
	beginArray
	endArray
	beginHex
	endHex

	// Stop
	error
	eof
)

type scanState struct  {
  isHex bool
	isKey bool
	Key string
	value string
	reader io.Reader
	offset int64
	stateStack []int
	err error
}

func (s *scanState) init(reader io.Reader) {
	s.stateStack = make([]int, 0)
	s.reader = reader
	s.offset = 0
	s.isHex = false
	s.isKey = false
	s.Key = ""
	s.value = ""
	s.err = nil
}

func (s *scanState) scan(){
	s.reader.	
}

func (s *scanState) step(){
	s.reader.
	if isSpace()
}
