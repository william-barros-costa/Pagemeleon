package main

import (
	"fmt"
	"io"
	"os"
)

func WritterWithOffsetSaver(file *os.File, offset *[]int64) func(string) {
	return func(s string) {
		pos, err := file.Seek(0, io.SeekCurrent)
		if err != nil {
			panic(fmt.Sprintf("Was not able to get current byte offset due to %s", err))
		}
		*offset = append(*offset, pos)
		fmt.Println(offset)
		file.WriteString(s)
	}
}

func main() {
	file, err := os.Create("test1.pdf")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	offsets := make([]int64, 0)
	write := WritterWithOffsetSaver(file, &offsets)
	write("%PDF-1.4\n")
	write("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")

	write("2 0 obj\n<< /Type /Pages /Count 1 /Kids [3 0 R] >>\nendobj\n")

	write("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R \n/Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n")

	stream := "BT /F1 12 Tf 72 720 Td (Hellow, World!) Tj ET"
	write(fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", len(stream), stream))

	write("5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n")

	xrefPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(fmt.Sprintf("Failed to seek current position due to %s", err))
	}
	file.WriteString("xref\n")
	file.WriteString(fmt.Sprintf("0 %d\n", len(offsets)+1))
	file.WriteString("0000000000 65535 f\n")
	for _, offset := range offsets {
		file.WriteString(fmt.Sprintf("%010d 00000 n\n", offset))
	}
	file.WriteString("trailer\n")
	file.WriteString(fmt.Sprintf("<< /Size %d /Root 1 0 R >>\n", len(offsets)))
	file.WriteString(fmt.Sprintf("startxref\n%d\n", xrefPos))
	file.WriteString("%%EOF\n")


	newOffsets := make([]int64, 0)
	newWritter := WritterWithOffsetSaver(file, &newOffsets)

	newWritter("6 0 obj\n<< /Type /Catalog /Pages 7 0 R >>\nendobj\n")

	newWritter("7 0 obj\n<< /Type /Pages /Count 1 /Kids [8 0 R] >>\nendobj\n")

	newWritter("8 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 9 0 R \n/Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n")

	newStream := "BT /F1 12 Tf 72 720 Td (Test change) Tj ET"
	newWritter(fmt.Sprintf("9 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", len(newStream), newStream))

	newXrefPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(fmt.Sprintf("Failed to seek current position due to %s", err))
	}
	file.WriteString("xref\n")
	file.WriteString(fmt.Sprintf("0 %d\n", len(newOffsets)))
	for _, offset := range newOffsets {
		file.WriteString(fmt.Sprintf("%010d 00000 n\n", offset))
	}
	file.WriteString("trailer\n")
	file.WriteString(fmt.Sprintf("<< /Size %d /Root 6 0 R /Prev %d >>\n", len(offsets), xrefPos))
	file.WriteString(fmt.Sprintf("startxref\n%d\n", newXrefPos))
	file.WriteString("%%EOF\n")
}

/*
%PDF-1.4
1 0 obj
<< /Type /Catalog /Pages 2 0 R >>
endobj
2 0 obj
<< /Type /Pages /Count 1 /kids [3 0 R] >>
endobj
3 0 obj
<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R >>
endobj
4 0 obj
<< /Length 44 >>
stream
BT /F1 12 Tf 72 720 Td (Hello, World!) Tj ET
endstream
endobj
xref
0 5
0000000000 65535 f
0000000010 00000 n
0000000060 00000 n
0000000110 00000 n
0000000190 00000 n
trailer
<< /Size 5 /Root 1 0 R >>
startxref
230
%%EOF
*/
