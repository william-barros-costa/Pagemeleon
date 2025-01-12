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
	block := BLOCK_SIZE
	size := Size(pdfReader)
	if BLOCK_SIZE > size {
		block = size
	}
	_, err := pdfReader.Seek(int64(-block), io.SeekEnd)
	if err != nil {
		fmt.Println(err)
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
		content = append(content, buffer...)
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
	return Trailer{}, nil
}

func Size(f *os.File) int64 {
	fileInfo, err := f.Stat()
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}
