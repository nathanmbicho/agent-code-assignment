package cmd

import (
	"bufio"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nathanmbicho/agent-code-assignment/pkg/components/listinput"
	"github.com/nathanmbicho/agent-code-assignment/pkg/components/textinput"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	openFileName    string
	showLineNumbers bool
)

type InputOptions struct {
	FileName *textinput.Output
}

type ListOptions struct {
	ListOptions *listinput.Selection
}

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
	//input command
	inputOptions := InputOptions{
		FileName: &textinput.Output{},
	}

	// handle program create, passing values
	tProgram := tea.NewProgram(textinput.InitialTextInputModel(
		inputOptions.FileName,
		"Enter file name to open ...",
		func(input string) (bool, error) {
			return validateSearchFile(input)
		},
	))

	// run bubbletea program
	if _, err := tProgram.Run(); err != nil {
		cobra.CheckErr(err)
		return
	}

	if inputOptions.FileName.Quit {
		fmt.Println("\n ❌Open file operation cancelled.")
	}

	// list command
	listOptions := ListOptions{
		ListOptions: &listinput.Selection{},
	}

	listOfOpenFileTools := []string{
		"Default",
		"Code",
	}

	tProgram = tea.NewProgram(listinput.InitialListInputModel(
		listOfOpenFileTools,
		inputOptions.FileName.Output,
		listOptions.ListOptions,
		"Select a tool to open with...",
		func(path, choice string) (string, bool, error) {
			return validateOpenFile(path, choice)
		},
	))

	if _, err := tProgram.Run(); err != nil {
		cobra.CheckErr(err)
		return
	}

	if listOptions.ListOptions.Quit {
		fmt.Println("\n ❌Create file operation cancelled.")
	}
}

// display file content in the cli
func displayFileContents(fileName string) (string, error) {
	// open file
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("error closing file %s - %v\n", fileName, err)
			return
		}
	}(file)

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	var list []string
	for scanner.Scan() {
		line := scanner.Text()
		text := fmt.Sprintf(fmt.Sprintf("%4d | %s", lineNumber, line))
		lineNumber++

		list = append(list, text)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(list, "\n"), nil
}

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
