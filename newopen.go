package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	bufSize = 512
	TRAILER_REFEX    = `(?:^|\s+)trailer([\s\S]*?)%%EOF`
)

func openPDF(location string) {
	if size, err := os.Stat(location); err != nil {
		errors.New(fmt.Sprintf("Cannot open file due to %s", err))
	} else if size.Size() == 0 {
		errors.New("File is empty so we will not continue")
	}
	file, err := os.Open(location)
	if err != nil {
		errors.New(fmt.Sprintf("Cannot open file due to %s", err))
	}

	buffer := make([]byte, 0)
	for iteration := 0; ; iteration++ {
		CurrentBuffer := make([]byte, bufSize)
		buffer = append(buffer, CurrentBuffer...)
		file.Seek(-int64(iteration*bufSize), io.SeekEnd)
		file.Read(CurrentBuffer)
		index := strings.LastIndex(string(CurrentBuffer), "startxref") 
		if index == -1 {
			continue	
		}

		
		
		endTrailer = strings.Index(string(CurrentBuffer), "%%EOF")
		trailer = CurrentBuffer[]
		
	}
}
