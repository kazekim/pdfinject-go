package pdfinject

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

//absPath create abs path for file
func absPath(fileName string) (string, error) {
	formPDFFile, err := filepath.Abs(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create the absolute path: %v", err)
	}

	return formPDFFile, nil
}

//isExist returns whether the given file or directory exists or not
func checkExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check if form PDF file exists: %v", err)
}

func checkFileExist(filePath string) (string, error) {

	// Get the absolute paths.
	filePath, err := absPath(filePath)
	if err != nil {
		return "", err
	}

	// Check if the file exists.
	isExist, err := checkExist(filePath)
	if err != nil {
		return "", err
	}else if !isExist {
		return "", fmt.Errorf("PDF file does not exists: '%s'", filePath)
	}
	return filePath, nil
}

// checkPkgExist check if the pkgName utility exists.
func checkPkgExist(pkgName string) error {

	_, err := exec.LookPath(pkgName)
	if err != nil {
		return fmt.Errorf("%s utility is not installed!", pkgName)
	}
	return nil
}

// removeFile delete file with filePath
func removeFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to remove PDF file: %v", err)
	}
	return nil
}

// copyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cErr := out.Close()
		if err == nil {
			err = cErr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()

	if err != nil {
		err = fmt.Errorf("failed to copy created output PDF to final destination: %v", err)
	}
	return
}
