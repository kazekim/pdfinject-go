/*
  GoLang code created by Jirawat Harnsiriwatanakit https://github.com/kazekim
*/

package pdfinject

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"log"
)

type TempFile struct {
	path string
}
// NewTempFile Create a FDF Generator.
func NewTempFile(dir, prefix string) (*TempFile, error) {

	tmpDir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %v", err)
	}
	return &TempFile{
		tmpDir,
	}, nil
}

// GetTempOutputFilePath Get a temporary output file path
func (t *TempFile) GetTempOutputFilePath() string {
	file := filepath.Clean(t.path + "/output.pdf")
	return file
}

// Remove delete temp directory when finish process
func (t *TempFile) Remove() {
	errD := os.RemoveAll(t.path)
	// Log the error only.
	if errD != nil {
		log.Printf("pdfinject: failed to remove temporary directory '%s' again: %v", t.path, errD)
	}
}