package pdfinject

import (
	"fmt"
)

const (
	pdfFormPkgName = "pdftk"
	dirPath = ""
	prefixFileName = "pdfinj-"
)

type PDFInject struct {
	form Form
	destPDFFile string
	overWrite bool
}


func New() PDFInject {
	return PDFInject{
		overWrite: true,
	}
}

func NewWithDestFile(destPDFFile string) PDFInject {
	return PDFInject{
		destPDFFile: destPDFFile,
		overWrite: true,
	}
}

// SetOverWrite allow overWrite to Destination file
func (pdf PDFInject) SetOverWrite(canOverwrited bool) {
	pdf.overWrite = canOverwrited
}

// Fill specified PDF form fields with the specified form values and export to a filled PDF file.
// One variadic boolean specifies, whenever to overwrite the destination file if it exists.

func (pdf PDFInject) Fill(form Form, formPDFFile string) error {

	// Get the absolute paths.
	formPDFFile, err := absPath(formPDFFile)
	if err != nil {
		return err
	}

	if pdf.destPDFFile == "" {
		return fmt.Errorf("dest file is not defined")
	}

	destPDFFile, err := absPath(pdf.destPDFFile)
	if err != nil {
		return err
	}

	// Check if the form file exists.
	isExist, err := checkExist(formPDFFile)
	if err != nil {
		return err
	}else if !isExist {
		return fmt.Errorf("PDF file does not exists: '%s'", formPDFFile)
	}

	// Check if the dest file exists
	isExist, err = checkExist(destPDFFile)
	if err != nil {
		return err
	}else if isExist {
		if !pdf.overWrite {
			return fmt.Errorf("destination PDF file already exists: '%s'", destPDFFile)
		}

		err = removeFile(destPDFFile)
		if err != nil {
			return fmt.Errorf("%s before create new one", err.Error())
		}
	}

	// Check if the pdftk utility exists.
	err = checkPkgExist(pdfFormPkgName)
	if err != nil {
		return err
	}

	// Create a temporary directory.
	tmpDir, err := NewTempDir(dirPath, prefixFileName)
	if err != nil {
		return err
	}

	// Remove the temporary directory on defer again.
	defer tmpDir.Remove()

	// Create the temporary output file path.
	outputFile := tmpDir.CreateTempOutputFile()

	// Create the fdf data file.
	fdfFile, err := tmpDir.CreateFDFFile(form)
	if err != nil {
		return err
	}

	// Create the pdftk command line arguments.
	args := pdf.createArgsTextOnly(formPDFFile, fdfFile, outputFile)

	// Run PDF Injector
	err = pdf.runInjector(tmpDir, args)
	if err != nil {
		return err
	}

	// On success, copy the output file to the final destination.
	err = copyFile(outputFile, destPDFFile)
	if err != nil {
		return err
	}

	return nil
}

func (pdf PDFInject) FillWithDestFile(form Form, formPDFFile, destPDFFile string) error {
	pdf.destPDFFile = destPDFFile

	return pdf.Fill(form, formPDFFile)
}

func (pdf PDFInject) createArgsTextOnly(formPDFFile, fdfFile, outputFile string) []string {
	args := []string{
		formPDFFile,
		"fill_form", fdfFile,
		"output", outputFile,
		"flatten",
	}

	return args
}

// runInjector Run the pdftk utility.
func (pdf PDFInject) runInjector(tmpDir *TempPDFDir, args []string) error {

	cmd := NewShellCommand(pdfFormPkgName)
	err := cmd.RunInPath(tmpDir.path, args...)
	if err != nil {
		return fmt.Errorf("pdftk error: %v", err)
	}

	return nil
}