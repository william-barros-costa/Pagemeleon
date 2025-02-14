/*
This package will contain all the functions required to Scan a PDF file.

Functions:

	[x] VerifyFile;
	[x] OpenFile;
	[x] VerifyFileIsPDF;
	[ ] ExtractTrailer;
	[ ] ExtractCrossReferenceTable;
	[ ] ExtractObject;
	[ ] ExtractVersion;
*/
package scan

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	PDF_HEADER       = "%PDF-x.x"
	BLOCK_SIZE int64 = 4096
)

// Verifies if file exists, has size > 0 and, is not a directory
func VerifyFile(fileLocation string) (bool, error) {
	fileInfo, err := os.Stat(fileLocation)
	if err != nil {
		return false, err
	}
	if fileInfo.Size() == 0 {
		return false, errors.New("File is empty")
	}
	if fileInfo.IsDir() {
		return false, fmt.Errorf("%q is a directory, not a file", fileLocation)
	}
	return true, nil
}

func VerifyFileIsPDF(pdfReader *os.File) (bool, error) {
	_, err := pdfReader.Seek(0, io.SeekStart)
	if err != nil {
		return false, err
	}
	buffer := make([]byte, len(PDF_HEADER))
	read, err := pdfReader.Read(buffer)
	if err != nil {
		return false, err
	}
	if read != len(PDF_HEADER) {
		return false, errors.New("Can't read PDF Header")
	}
	if !strings.HasPrefix(string(buffer), "%PDF-") {
		return false, fmt.Errorf("Expected string of type \"%%PDF-x.x\", got %q", buffer)
	}
	return true, nil
}

func OpenFile(fileLocation string) (*os.File, error) {
	if valid, err := VerifyFile(fileLocation); !valid {
		return nil, err
	}
	file, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func ExtractTrailer(pdfReader *os.File) (
	Trailer,
	error,
) {
	// If block is bigger than file size an error occurs so, we address this below
	block := BLOCK_SIZE
	size := Size(pdfReader)
	if BLOCK_SIZE > size {
		block = size
	}
	_, err := pdfReader.Seek(int64(-block), io.SeekEnd)
	if err != nil {
		return Trailer{}, err
	}
	buffer := make([]byte, BLOCK_SIZE)
	var content []byte
	found := false
	for {
		_, err := pdfReader.Read(buffer)
		if err != nil {
			return Trailer{}, err
		}
		content = append(buffer, content...)
		if strings.Contains(string(content), "trailer") {
			found = true
			break
		}
		read, err := pdfReader.Seek(int64(-block), io.SeekCurrent)
		if read == 0 {
			break
		}
		if err == nil {
			continue
		}
		if err == io.EOF {
			break
		}
		return Trailer{}, err
	}

	if !found {
		return Trailer{}, errors.New("PDF is missing trailer keyword")
	}

	index := bytes.Index(content, []byte("trailer"))
	content = content[index:]

	for _, member := range []string{"%%EOF", "startxref"} {
		if !bytes.Contains(content, []byte(member)) {
			return Trailer{}, fmt.Errorf("Trailer is missing %s keyword", member)
		}
	}

	startDictionary := bytes.Index(content, []byte("trailer")) + len("trailer")
	endDictionary := bytes.Index(content, []byte("startxref"))

	if startDictionary >= endDictionary {
		return Trailer{}, errors.New("Trailer is missing dictionary")
	}

	dictionary := bytes.TrimSpace(content[startDictionary:endDictionary])
	if len(dictionary) == 0 {
		return Trailer{}, errors.New("Trailer is missing dictionary")
	}

	if !bytes.Contains(dictionary, []byte("/Size")) {
		return Trailer{}, errors.New("trailer's dictionary is missing size keyword")
	}

	//	var parts [][]byte = SplitDictionary(dictionary)
	//	for _, part := range parts {
	//		fmt.Println(string(part))
	//	}

	return Trailer{}, nil
}

func Size(f *os.File) int64 {
	fileInfo, err := f.Stat()
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}

func SplitDictionary(dictionary []byte) [][]byte {
	elements := make([][]byte, 0)
	current := make([]byte, 0)
	for _, b := range dictionary {
		if b == byte('<') || b == byte('>') {
			continue
		}
		if isWhitespace(b) || b == byte('/') {
			if len(current) > 0 {
				elements = append(elements, current)
				current = make([]byte, 0)
			}
			continue
		}
		current = append(current, b)
	}
	return append(elements, current)
}

//func ScanParts(parts [][]byte) {
//	for _, part := range parts {
//		if bytes.HasPrefix(part, []byte("/")) {
//			fmt.Println(part)
//		}
//	}
//}

func Scan(dictionary []byte) []Object {
	objects := make([]Object, 0)
	currentObject := make([][]byte, 0)
	current := make([]byte, 0)
	for _, b := range dictionary {
		if b == byte('/') {
			currentObject = make([][]byte, 0)
			current = make([]byte, 0)
		} else if isWhitespace(b) {
			currentObject = append(currentObject, current)
			current = make([]byte, 0)
		} else if b == byte('R') && len(currentObject) == 3 {
			fmt.Println(len(currentObject))
			fmt.Println(currentObject[2])
			objects = append(objects, Object{
				Name:       string(currentObject[0]),
				Id:         toUint(currentObject[1]),
				Generation: toUint(currentObject[2]),
			})
			currentObject = make([][]byte, 0)
			current = make([]byte, 0)
		} else {
			current = append(current, b)
		}
	}
	return objects
}

func toUint(b []byte) uint {
	number, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0
	}
	return uint(number)
}
