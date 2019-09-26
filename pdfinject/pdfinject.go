package pdfinject

import (
	"errors"
	"fmt"
)

const (
	pdfFormPkgName = "pdftk"
	javaPkgName = "java"
	dirPath = ""
	prefixFileName = "pdfinj-"
)

type PDFInject struct {
	form map[string]interface{}
	destPDFFile string
	overWrite bool
	inputType InputType
}


func New() PDFInject {
	return PDFInject{
		overWrite: true,
		inputType: XFDF,
	}
}

func NewWithDestFile(destPDFFile string) PDFInject {
	return PDFInject{
		destPDFFile: destPDFFile,
		overWrite: true,
		inputType: XFDF,
	}
}

// SetOverWrite allow overWrite to Destination file
func (pdf PDFInject) SetOverWrite(canOverwrited bool) {
	pdf.overWrite = canOverwrited
}

func (pdf PDFInject) SetInputType(in InputType) {
	pdf.inputType = in
}

// Fill specified PDF form fields with the specified form values and export to a filled PDF file.
// One variadic boolean specifies, whenever to overwrite the destination file if it exists.

func (pdf PDFInject) Fill(form map[string]interface{}, formPDFFile string) (*string, error) {

	// Check if the form file exists.
	formPDFFile, err := checkFileExist(formPDFFile)
	if err != nil {
		return nil, err
	}

	// Check if the form file exists.
	isExist, err := checkExist(formPDFFile)
	if err != nil {
		return nil, err
	}else if !isExist {
		return nil, fmt.Errorf("PDF file does not exists: '%s'", formPDFFile)
	}

	// Check if the dest file exists
	destPDFFile, err := pdf.checkDestFileExist()
	if err != nil {
		return nil, err
	}

	// Check if the pdftk utility exists.
	err = checkPkgExist(pdfFormPkgName)
	if err != nil {
		return nil, err
	}

	outputFile, inputFile, tempFile, err := pdf.generateInputDataFile(form)
	if err != nil {
		return nil, err
	}
	// Remove the temporary directory on defer again.
	defer tempFile.Remove()

	// Create the pdftk command line arguments.
	args := pdf.createArgsTextOnly(formPDFFile, *inputFile, *outputFile)

	// Run PDF Injector
	err = pdf.runInjectorWithStd(javaPkgName, tempFile.path, *inputFile, *outputFile, args)
	if err != nil {
		return nil, err
	}

	// On success, copy the output file to the final destination.
	err = copyFile(*outputFile, destPDFFile)
	if err != nil {
		return nil, err
	}

	return &destPDFFile, nil
}

func (pdf PDFInject) FillWithDestFile(form map[string]interface{}, formPDFFile, destPDFFile string) (*string, error) {
	pdf.destPDFFile = destPDFFile

	newForm := prepareMap(form)

	return pdf.Fill(*newForm, formPDFFile)
}

func (pdf PDFInject) FillModel(model interface{}, formPDFFile string) (*string, error) {

	form := structToForm(model)

	return pdf.Fill(form, formPDFFile)
}

func (pdf PDFInject) FillModelWithDestFile(model interface{}, formPDFFile string, destPDFFile string) (*string, error) {

	pdf.destPDFFile = destPDFFile
	form := structToForm(model)

	return pdf.Fill(form, formPDFFile)
}

func (pdf PDFInject) Stamp(stampPDFFile, srcPDFFile string) error {

	// Check if the Stamp file exists.
	stampPDFFile, err := checkFileExist(stampPDFFile)
	if err != nil {
		return err
	}

	// Check if the Source file exists.
	srcPDFFile, err = checkFileExist(srcPDFFile)
	if err != nil {
		return err
	}

	// Check if the dest file exists
	destPDFFile, err := pdf.checkDestFileExist()
	if err != nil {
		return nil
	}

	// Check if the pdftk utility exists.
	err = checkPkgExist(pdfFormPkgName)
	if err != nil {
		return err
	}

	// Create a temporary directory.
	tmpFile, err := NewTempFile(dirPath, prefixFileName)
	if err != nil {
		return err
	}

	// Remove the temporary directory on defer again.
	defer tmpFile.Remove()

	// Create the temporary output file path.
	outputFile := tmpFile.path

	// Create the pdftk command line arguments.
	args := pdf.createArgsStampPDF(srcPDFFile, stampPDFFile, outputFile)

	// Run PDF Injector
	err = pdf.runInjector(pdfFormPkgName, tmpFile.path, args)
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

func (pdf PDFInject) StampWithDestFile(stampPDFFile, srcPDFFile, destPDFFile string) error {
	pdf.destPDFFile = destPDFFile

	return pdf.Stamp(stampPDFFile, srcPDFFile)
}

//createArgsTextOnly add text from struct to PDF
func (pdf PDFInject) createArgsTextOnly(formPDFFile, fdfFile, outputFile string) []string {

	home, err := HomeDirectory()
	if err != nil {
		optPath := "/opt"
		home = &optPath
	}
	args := []string{
		"-jar", *home + "/jar/mcpdf.jar",
		formPDFFile,
		"fill_form","-",
		"output", "-",
	}

	return args
}

//createArgsStampPDF stamp pdf on pdf
func (pdf PDFInject) createArgsStampPDF(srcPDFFile, stampPDFFile, outputFile string) []string {
	args := []string{
		srcPDFFile,
		"multistamp", stampPDFFile,
		"output", outputFile,
	}

	return args
}

// runInjector Run the pdftk utility.
func (pdf PDFInject) runInjector(cmdName string, tmpDir string, args []string) error {

	cmd := NewShellCommand(cmdName)
	err := cmd.RunInPath(tmpDir, args...)
	if err != nil {
		return fmt.Errorf("pdftk error: %v", err)
	}

	return nil
}

func (pdf PDFInject) runInjectorWithStd(cmdName string, tmpDir string, stdInPath string, stdOutPath string, args []string) error {

	cmd := NewShellCommand(cmdName)
	err := cmd.RunInPathWithStd(tmpDir, stdInPath, stdOutPath, args...)
	if err != nil {
		return fmt.Errorf("pdftk error: %v", err)
	}

	return nil
}

func (pdf PDFInject) checkDestFileExist() (string, error) {

	if pdf.destPDFFile == "" {
		return "", fmt.Errorf("dest file is not defined")
	}

	destPDFFile, err := absPath(pdf.destPDFFile)
	if err != nil {
		return "", err
	}

	// Check if the dest file exists
	isExist, err := checkExist(destPDFFile)
	if err != nil {
		return "", err
	}else if isExist {
		if !pdf.overWrite {
			return "", fmt.Errorf("destination PDF file already exists: '%s'", destPDFFile)
		}

		err = removeFile(destPDFFile)
		if err != nil {
			return "", fmt.Errorf("%s before create new one", err.Error())
		}
	}

	return destPDFFile, nil
}

func (pdf PDFInject) generateInputDataFile(form map[string]interface{}) (*string, *string, *TempFile, error){

	var outputFile string
	var inputFile string
	switch pdf.inputType {
	case XFDF:
		// Create a temporary directory.
		xfdf, err := NewXFDFGenerator(dirPath, prefixFileName)
		if err != nil {
			return nil, nil, nil, err
		}

		// Create the temporary output file path.
		outputFile = xfdf.GetTempOutputFile()

		// Create the fdf data file.
		inputFile, err = xfdf.CreateXFDFFile(form)
		if err != nil {
			return nil, nil, nil, err
		}

		return &outputFile, &inputFile, &xfdf.file, nil
	case FDF:
		// Create a temporary directory.
		fdf, err := NewFDFGenerator(dirPath, prefixFileName)
		if err != nil {
			return nil, nil, nil, err
		}

		// Create the temporary output file path.
		outputFile = fdf.GetTempOutputFile()

		// Create the fdf data file.
		inputFile, err = fdf.CreateFDFFile(form)
		if err != nil {
			return nil, nil, nil, err
		}

		return &outputFile, &inputFile, &fdf.file, nil
	}
	return nil, nil, nil, errors.New("undefined input file")
}