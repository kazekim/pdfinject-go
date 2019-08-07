package pdfinject

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)
// Form represents fields from the PDF form.
// define in key value map.
type Form map[string]interface{}

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
type TempPDFDir struct {
	path string
}

// NewTempDir Create a temporary directory.
func NewTempDir(dir, prefix string) (*TempPDFDir, error) {

	tmpDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %v", err)
	}
	return &TempPDFDir{
		tmpDir,
	}, nil
}

// Remove delete temp directory when finish process
func (t *TempPDFDir) Remove() {
	errD := os.RemoveAll(t.path)
	// Log the error only.
	if errD != nil {
		log.Printf("fillpdf: failed to remove temporary directory '%s' again: %v", t.path, errD)
	}
}

// CreateTempOutputFile Create a temporary output file
func (t *TempPDFDir) CreateTempOutputFile() string {
	file := filepath.Clean(t.path + "/output.pdf")
	return file
}

// CreateFDFFile Create a temporary fdf file
func (t *TempPDFDir) CreateFDFFile(form Form) (string, error) {
	fdfFile := filepath.Clean(t.path + "/data.fdf")
	err := t.generateFdfFile(form, fdfFile)
	if err != nil {
		return "", fmt.Errorf("failed to create fdf form data file: %v", err)
	}

	return fdfFile, nil
}

// Generate FDF file with parameters to inject to PDF
func (t *TempPDFDir) generateFdfFile(form Form, path string) error {
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
	for key, value := range form {
		fmt.Fprintf(w, "<< /T (%s) /V (%v)>>\n", key, value)
	}

	// Write the fdf footer.
	_, _ = fmt.Fprintln(w, fdfFooter)

	// Flush everything.
	return w.Flush()
}