package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nathanmbicho/agent-code-assignment/pkg/components/textinput"
	"github.com/nathanmbicho/agent-code-assignment/pkg/ui"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var fileName string

type CreateOptions struct {
	FileName *textinput.Output
}

// createFileCmd - create a new file
var createFileCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new file for a given programming language",
	Long:  `Creating a new file for a given programming language. You can create a file of any of the following languages: go, js, py, php`,
	Run:   createFile,
}

func init() {
	rootCmd.AddCommand(createFileCmd)
}

func createFile(cmd *cobra.Command, args []string) {
	allowedExtensions := []string{".go", ".js", ".py", ".php"}

	options := CreateOptions{
		FileName: &textinput.Output{},
	}

	// handle program create, passing values
	tProgram := tea.NewProgram(textinput.InitialTextInputModel(
		options.FileName,
		fmt.Sprintf("Create a new file. Allowed languages are %s", strings.Join(allowedExtensions, ",")),
		func(input string) (bool, error) {
			return validateFileCreate(input, allowedExtensions)
		},
	))

	// run bubbletea program
	if _, err := tProgram.Run(); err != nil {
		cobra.CheckErr(err)
		return
	}

	if options.FileName.Quit {
		fmt.Println("\n ❌Create file operation cancelled.")
	}

	fileName = options.FileName.Output

	if fileName != "" {
		success := ui.RenderSuccess(fmt.Sprintf("file '%s' created successfully!", fileName))
		fmt.Printf(success)
	}
}

// check allowed file extensions
func isValidExtension(ext string, allowedExts []string) bool {
	for _, allowedExt := range allowedExts {
		if allowedExt == ext {
			return true
		}
	}
	return false
}

// generate file template
func generateFileTemplate(fileName string) string {
	ext := filepath.Ext(fileName)

	switch ext {
	case ".go":
		fmt.Printf("Generating file %s ... \n", ui.GoFileStyle.Render(fmt.Sprintf("%s", fileName)))
	case ".js":

		fmt.Printf("Generating file %s ... \n", ui.JSFileStyle.Render(fmt.Sprintf("%s", fileName)))
	case ".py":

		fmt.Printf("Generating file %s ... \n", ui.PythonFileStyle.Render(fmt.Sprintf("%s", fileName)))
	case ".php":

		fmt.Printf("Generating file %s ... \n", ui.PHPFileStyle.Render(fmt.Sprintf("%s", fileName)))
	}

	switch ext {
	case ".go":
		return `package main

import "fmt"

func main(){
	fmt.Println("Hello world")
}
`
	case ".js":
		return `console.log("Hello world");`
	case ".py":
		return `print ("Hello world")`
	case ".php":
		return `<?php
echo "Hello world"; 
?> 
`
	default:
		return ""
	}
}

// create directory
func generateFileDirectory(fileName string) error {
	dir := filepath.Dir(fileName)

	if dir != "." && dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}
	}

	return nil
}
