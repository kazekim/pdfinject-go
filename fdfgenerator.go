/*
  GoLang code created by Jirawat Harnsiriwatanakit https://github.com/kazekim
*/

package pdfinject

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)


const fdfHeader = `%FDF-1.2
%,,oe"
1 0 obj
<<
/FDF << /Fields [`

const fdfFooter = `]
>>
>>
endobj
trailer
<<
/Root 1 0 R
>>
%%EOF`

// TempPDFDir manage PDF Inject process
type FDFGenerator struct {
	file TempFile
}

// NewFDFGenerator Create a FDF Generator.
func NewFDFGenerator(dir, prefix string) (*FDFGenerator, error) {

	file, err := NewTempFile(dir, prefix)
	if err != nil {
		return nil, err
	}
	return &FDFGenerator{
		*file,
	}, nil
}

// Remove delete temp directory when finish process
func (t *FDFGenerator) Remove() {
	t.file.Remove()
}

// CreateFDFFile Create a temporary fdf file
func (t *FDFGenerator) CreateFDFFile(form Form) (string, error) {
	fdfFile := filepath.Clean(t.path() + "/data.fdf")
	err := t.generateFdfFile(form, fdfFile)
	if err != nil {
		return "", fmt.Errorf("failed to create fdf form data file: %v", err)
	}

	return fdfFile, nil
}

// path get temp file path
func (t *FDFGenerator) path() string {
	return t.file.path
}

// GetTempOutputFile get temp output file for pdf
func (t *FDFGenerator) GetTempOutputFile() string {
	return t.file.GetTempOutputFilePath()
}

// Generate FDF file with parameters to inject to PDF
func (t *FDFGenerator) generateFdfFile(form Form, path string) error {
	// Create the file.
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new writer.
	w := bufio.NewWriter(file)

	// Write the fdf header.
	_, _ = fmt.Fprintln(w, fdfHeader)

	// Write the form data.
	for key, _ := range form {
		fmt.Fprintf(w, "<< /T (%s) /V (\xfe\xff\xbd \xe01) >>\n", key)
	}

	// Write the fdf footer.
	_, _ = fmt.Fprintln(w, fdfFooter)

	// Flush everything.
	return w.Flush()
}
