package scan

import (
	"testing"
)

const (
	SamplePDF       = "../../pdfs/sample.pdf"
	SampleDirectory = "../../pdfs"
)

func TestVerifyFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		fileLocation string
		want         bool
		wantErr      bool
		err          string
	}{
		{name: "Sample PDF", fileLocation: SamplePDF, want: true, wantErr: false},
		{name: "Missing PDF", fileLocation: "unknown", wantErr: true, err: "stat unknown: no such file or directory"},
		{name: "Sample Directory", fileLocation: SampleDirectory, wantErr: true, err: "\"../../pdfs\" is a directory, not a file"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := VerifyFile(tt.fileLocation)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("VerifyFile() failed: %v", gotErr)
				} else if tt.err != gotErr.Error() {
					t.Errorf("VerifyFile() = %v, want %v", gotErr, tt.err)
				}
				return
			}
			if got != tt.want {
				t.Errorf("VerifyFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
