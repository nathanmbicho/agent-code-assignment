package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// validateSearchFile - validate file path input value
func validateSearchFile(fileName string) (bool, error) {
	// check filename if is empty
	if strings.TrimSpace(fileName) == "" {
		return false, fmt.Errorf("filename cannot be empty")
	}

	//check if the file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false, fmt.Errorf("file %s does not exist. incorrect path or file name \n", fileName)
	}

	// get the file path
	path, err := filepath.Abs(fileName)
	if err != nil {
		return false, fmt.Errorf("error getting absolute path %s\n", path)
	}

	return true, nil
}

// validation
func validateFileCreate(fileName string, allowedExtensions []string) (bool, error) {
	// check filename if is empty
	if strings.TrimSpace(fileName) == "" {
		return false, fmt.Errorf("filename cannot be empty")
	}

	ext := filepath.Ext(fileName)

	// check if the file extension is valid
	if !isValidExtension(ext, allowedExtensions) {
		return false, fmt.Errorf("invalid file extension. allowed: %s", strings.Join(allowedExtensions, ", "))
	}

	// check if the file already exists
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return false, fmt.Errorf("file '%s' already exists", fileName)
	}

	// generate directory if included in the file path
	if err := generateFileDirectory(fileName); err != nil {
		return false, fmt.Errorf("error creating directory: %v", err)
	}

	// create the file
	file, err := os.Create(fileName)
	if err != nil {
		return false, fmt.Errorf("error creating file: %v", err)
	}

	// generate file template and write to file
	temp := generateFileTemplate(fileName)
	if temp != "" {
		if _, err := file.WriteString(temp); err != nil {
			err := file.Close()
			if err != nil {
				return false, fmt.Errorf("error closing file: %v", err)
			}
			return false, fmt.Errorf("error writing to file: %v", err)
		}
	}

	// close file
	if err := file.Close(); err != nil {
		return false, fmt.Errorf("error closing file: %v", err)
	}

	return true, nil
}

// validateOpenFile - validate open file
func validateOpenFile(fileName, editor string) (string, bool, error) {
	var cmd *exec.Cmd

	if strings.TrimSpace(fileName) == "" {
		return "", false, fmt.Errorf("filename cannot be empty")
	}

	if strings.TrimSpace(editor) == "" {
		return "", false, fmt.Errorf("error reading open file %s\n", fileName)
	}

	// get the file path
	path, err := filepath.Abs(fileName)
	if err != nil {
		return "", false, fmt.Errorf("error getting absolute path %s\n", path)
	}

	switch strings.ToLower(editor) {
	case "code":
		cmd = exec.Command("code", path)
		return "", true, cmd.Start()
	case "vim":
		cmd = exec.Command("vim", path)
		return "", true, cmd.Start()
	default:
		// display file data
		code, err := displayFileContents(fileName)
		if err != nil {
			return "", false, fmt.Errorf("error opening file %s - %v\n", path, err)
		}
		return code, true, nil
	}
}
