package scan

import (
	"testing"
)

const (
	SamplePDF             = "../../pdfs/sample.pdf"
	WrongHeader           = "../../pdfs/wrong_header.pdf"
	TRAILER_NO_EOF        = "../../pdfs/trailer_no_EOF.pdf"
	TRAILER_NO_ROOT       = "../../pdfs/trailer_no_root.pdf"
	TRAILER_NO_STARTXREF  = "../../pdfs/trailer_no_startxref.pdf"
	TRAILER_NO_TRAILER    = "../../pdfs/trailer_no_trailer.pdf"
	TRAILER_NO_DICTIONARY = "../../pdfs/trailer_no_dictionary.pdf"
	SmallFile             = "../../pdfs/small.pdf"
	SampleDirectory       = "../../pdfs"
)

func TestVerifyFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		fileLocation string
		want         bool
		wantErr      bool
		error        string
	}{
		{name: "Sample PDF", fileLocation: SamplePDF, want: true, wantErr: false},
		{name: "Missing PDF", fileLocation: "unknown", wantErr: true, error: "stat unknown: no such file or directory"},
		{name: "Sample Directory", fileLocation: SampleDirectory, wantErr: true, error: "\"../../pdfs\" is a directory, not a file"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := VerifyFile(tt.fileLocation)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("VerifyFile() failed: %v", gotErr)
				} else if tt.error != gotErr.Error() {
					t.Errorf("VerifyFile() = %v, want %v", gotErr, tt.error)
				}
				return
			}
			if got != tt.want {
				t.Errorf("VerifyFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyFileIsPDF(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		location string
		want     bool
		wantErr  bool
		error    string
	}{
		{name: "Sample PDF", location: SamplePDF, want: true},
		{name: "Wrong Header", location: WrongHeader, wantErr: true, error: "Expected string of type \"%PDF-x.x\", got \"%PDA-1.3\""},
		{name: "Small PDF", location: SmallFile, wantErr: true, error: "Can't read PDF Header"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := OpenFile(tt.location)
			got, gotErr := VerifyFileIsPDF(file)

			if gotErr != nil && !tt.wantErr {
				t.Errorf("VerifyFile() failed: %v", gotErr)
			} else if gotErr != nil && tt.error != gotErr.Error() {
				t.Errorf("VerifyFile() Error = %v, want %v", gotErr, tt.error)
			} else if tt.wantErr && gotErr == nil {
				t.Error("Expected Error but did not get it")
			} else if got != tt.want {
				t.Errorf("VerifyFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTrailer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		location string
		want     Trailer
		wantErr  bool
		error    string
	}{
		{name: "Sample PDF", location: SamplePDF, want: Trailer{Root: Object{}, Encrypt: Object{}, Info: Object{}, Ids: [][]byte{}}, wantErr: false},
		{name: "Trailer has no Trailer Keyword", location: TRAILER_NO_TRAILER, wantErr: true, error: "PDF is missing trailer keyword"},
		{name: "Trailer has no EOF", location: TRAILER_NO_EOF, wantErr: true, error: "Trailer is missing %%EOF keyword"},
		{name: "Trailer has no Root Object", location: TRAILER_NO_ROOT, wantErr: true, error: "Trailer is missing Root Object"},
		{name: "Trailer has no startxref Keyword", location: TRAILER_NO_STARTXREF, wantErr: true, error: "Trailer is missing startxref keyword"},
		{name: "Trailer has no Dictionary", location: TRAILER_NO_DICTIONARY, wantErr: true, error: "Trailer is missing dictionary"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := OpenFile(tt.location)
			got, gotErr := ExtractTrailer(file)

			if gotErr != nil && !tt.wantErr {
				t.Errorf("VerifyFile() failed: %v", gotErr)
			} else if gotErr != nil && tt.error != gotErr.Error() {
				t.Errorf("VerifyFile() Error = %v, want %v", gotErr, tt.error)
			} else if tt.wantErr && gotErr == nil {
				t.Error("Expected Error but did not get it")
			} else if !got.Equal(tt.want) {
				t.Errorf("VerifyFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
