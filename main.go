package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func PdfSize(location string) uint {
	fileInfo, err := os.Stat(location)
	if err != nil {
		panic(err)
	}
	return uint(fileInfo.Size())
}

func openPdf(location string) PDF {
	file, err := os.Open(location)
	if err != nil {
		panic(err)
	}
	scanner := PdfScanner{
		Scanner: *bufio.NewScanner(file),		
		XrefTable: make(map[uint]Xref),
	}

	scanner.Scan()	

	for _, value := range scanner.XrefTable {
		fmt.Println(value)
		time.Sleep(1*time.Second)
	}


	return PDF{
		Header: Header{
			version: "PDF 1.7",
			size: PdfSize(location),
		},
	}
}

func main() {
	fmt.Println(openPdf("./pdf.pdf"))
}
