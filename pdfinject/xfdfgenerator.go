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

const xfdfHeader = `<?xml version="1.0" encoding="UTF-8"?>
<xfdf xmlns="http://ns.adobe.com/xfdf/" xml:space="preserve">
   <fields>`

const xfdfFooter = `</fields>
</xfdf>`

// TempPDFDir manage PDF Inject process
type XFDFGenerator struct {
	file TempFile
}

// NewXFDFGenerator Create a XFDF Generator.
func NewXFDFGenerator(dir, prefix string) (*XFDFGenerator, error) {

	file, err := NewTempFile(dir, prefix)
	if err != nil {
		return nil, err
	}
	return &XFDFGenerator{
		*file,
	}, nil
}


// Remove all temp file when finish process
func (t *XFDFGenerator) Remove() {
	t.file.Remove()
}

// CreateXFDFFile Create a temporary fdf file
func (t *XFDFGenerator) CreateXFDFFile(form map[string]interface{}) (string, error) {
	fdfFile := filepath.Clean(t.path() + "/data.xfdf")
	err := t.generateXFdfFile(form, fdfFile)
	if err != nil {
		return "", fmt.Errorf("failed to create fdf form data file: %v", err)
	}

	return fdfFile, nil
}

// path get temp file path
func (t *XFDFGenerator) path() string {
	return t.file.path
}

// GetTempOutputFile get temp output file for pdf
func (t *XFDFGenerator) GetTempOutputFile() string {
	return t.file.GetTempOutputFilePath()
}

// Generate XFDF file with parameters to inject to PDF
func (t *XFDFGenerator) generateXFdfFile(form map[string]interface{}, path string) error {

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