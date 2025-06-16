package cmd

import (
	"bufio"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nathanmbicho/agent-code-assignment/pkg/components/textinput"
	"github.com/nathanmbicho/agent-code-assignment/pkg/ui"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
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
}

func openFile(cmd *cobra.Command, args []string) {
	options := Options{
		FileName: &textinput.Output{},
	}

	// handle program create, passing values
	tProgram := tea.NewProgram(textinput.InitialTextInputModel(
		options.FileName,
		"Enter file name to open ...",
		func(input string) (bool, error) {
			return validateFileOpen(input)
		},
	))

	// run bubbletea program
	if _, err := tProgram.Run(); err != nil {
		cobra.CheckErr(err)
		return
	}

	if options.FileName.Quit {
		fmt.Println("\n ❌Open file operation cancelled.")
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

	fmt.Printf("\nFile : %s\n\n", absolutePath)

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(ui.RenderCode(fmt.Sprintf("%4d | %s", lineNumber, line)))
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// validate file path input value
func validateFileOpen(fileName string) (bool, error) {
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

	// display file data
	err = displayFileContents(fileName, path)
	if err != nil {
		return false, fmt.Errorf("error opening file %s - %v\n", path, err)
	}

	return true, nil
}
