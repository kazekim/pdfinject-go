/*
  GoLang code created by Jirawat Harnsiriwatanakit https://github.com/kazekim
*/

package pdfinject

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"path/filepath"
)

// Form represents fields from the PDF form.
// define in key value map.
type XFDFForm map[string]interface{}

const xfdfHeader = `<?xml version="1.0" encoding="UTF-8"?>
<xfdf xmlns="http://ns.adobe.com/xfdf/" xml:space="preserve">
   <fields>`

const xfdfFooter = `</fields>
</xfdf>`

// TempPDFDir manage PDF Inject process
type XFDFGenerator struct {
	path string
}

// NewXFDFGenerator Create a XFDF Generator.
func NewXFDFGenerator(dir, prefix string) (*XFDFGenerator, error) {

	tmpDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %v", err)
	}
	return &XFDFGenerator{
		tmpDir,
	}, nil
}


// Remove all temp file when finish process
func (t *XFDFGenerator) Remove() {
	errD := os.RemoveAll(t.path)
	// Log the error only.
	if errD != nil {
		log.Printf("pdfinjector: failed to remove temporary directory '%s' again: %v", t.path, errD)
	}
}

// CreateTempOutputFile Create a temporary output file
func (t *XFDFGenerator) CreateTempOutputFile() string {
	file := filepath.Clean(t.path + "/output.pdf")
	return file
}

// CreateXFDFFile Create a temporary fdf file
func (t *XFDFGenerator) CreateXFDFFile(form Form) (string, error) {
	fdfFile := filepath.Clean(t.path + "/data.xfdf")
	err := t.generateXFdfFile(form, fdfFile)
	if err != nil {
		return "", fmt.Errorf("failed to create fdf form data file: %v", err)
	}

	return fdfFile, nil
}

// Generate XFDF file with parameters to inject to PDF
func (t *XFDFGenerator) generateXFdfFile(form Form, path string) error {
	// Create the file.
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new writer.
	w := bufio.NewWriter(file)

	// Write the fdf header.
	_, _ = fmt.Fprintln(w, xfdfHeader)
	// Write the form data.
	for key, value := range form {
		fmt.Fprintf(w, "<field name=\"%s\">" +
			"<value>%s</value>" +
		"</field>", key, value)
	}

	// Write the fdf footer.
	_, _ = fmt.Fprintln(w, xfdfFooter)

	// Flush everything.
	return w.Flush()
}