package main

import (
	"bufio"
	// "fmt"
	// "io"
	"os"
	// "time"
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
	pdf := PDF{
		Header: Header{
			version: "PDF 1.7",
			size: PdfSize(location),
		},
		Objects: make([]Object, 0),
		Pages: make([]Page, 0),
		XrefTable: make(map[uint]Xref),
	}
	scanner := PdfScanner{
		Scanner: *bufio.NewScanner(file),		
		pdf: &pdf,
	}

	scanner.Scan()	

	// for _, value := range pdf.xreftable {
	// 	fmt.println(value)
	// 	time.sleep(1*time.second)
	// }

	return pdf
}


func main() {
	// fmt.Println(openPdf("./pdf.pdf"))
	// open("./test1.pdf")
	open("pdf.pdf")
}

