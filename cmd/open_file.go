package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	openFileName    string
	showLineNumbers bool
)

// openFileCmd - one file
var openFileCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the file in the current directory",
	Long:  `Open the file in the current specified directory. File opened must exist in the current directory and will open on the terminal.`,
	Run:   openFile,
}

func init() {
	rootCmd.AddCommand(openFileCmd)

	openFileCmd.Flags().StringVarP(&openFileName, "file", "f", "", "file to be opened")
	openFileCmd.Flags().BoolVarP(&showLineNumbers, "line-numbers", "l", false, "show line numbers - default: false")

	err := openFileCmd.MarkFlagRequired("file")
	if err != nil {
		fmt.Println("error while marking flag as required")
		return
	}
}

func openFile(cmd *cobra.Command, args []string) {
	//check if the file exists
	if _, err := os.Stat(openFileName); os.IsNotExist(err) {
		fmt.Printf("error - file %s does not exist. incorrect path or file name \n", openFileName)
		return
	}

	// get the file path
	path, err := filepath.Abs(openFileName)
	if err != nil {
		fmt.Printf("error - error getting path %s\n", openFileName)
		return
	}

	// display file data
	err = displayFileContents(openFileName, path)
	if err != nil {
		fmt.Printf("error opening file %s - %v\n", path, err)
		return
	}
}

// display file content in the cli
func displayFileContents(fileName, absolutePath string) error {
	// open file
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file %s - %v\n", fileName, err)
			return
		}
	}(file)

	fmt.Printf("\n")
	fmt.Printf("File : %s\n", absolutePath)
	fmt.Printf("\n")

	// scan read and display file content
	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()

		// check if show line numbers
		if showLineNumbers {
			fmt.Printf("%4d | %s\n", lineNumber, line)
		} else {
			fmt.Printf("%s\n", line)
		}

		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("\n\n")

	return nil
}
