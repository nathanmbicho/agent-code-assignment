package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var fileName string

// createFileCmd - create a new file
var createFileCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new file for a given programming language",
	Long:  `Creating a new file for a given programming language. You can create a file of any of the following languages: go, js, py, php`,
	Run:   createFile,
}

func init() {
	rootCmd.AddCommand(createFileCmd)

	createFileCmd.Flags().StringVarP(&fileName, "file", "f", "", "The file name to create")
	err := createFileCmd.MarkFlagRequired("file")
	if err != nil {
		fmt.Println("Error while marking flag as required")
		return
	}
}

func createFile(cmd *cobra.Command, args []string) {

	allowedExtensions := []string{".go", ".js", ".py", ".php"}
	ext := filepath.Ext(fileName)

	// check if file created extention exists
	if !isValidExtension(ext, allowedExtensions) {
		fmt.Printf("Invalid file extension '%s'. Allowed languages are %s \n", ext, strings.Join(allowedExtensions, ","))
		return
	}

	// check if the file exists
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		fmt.Printf("Oops! File '%s' already exists.\n", fileName)
		return
	}

	// create a file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file %s ::- %v\n", fileName, err)
		return
	}

	// defer close
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file %s ::- %v\n", fileName, err)
		}

		fmt.Printf("File %s created successfully. \n", fileName)
	}(file)

	// generate file template
	temp := generateFileTemplate(fileName)
	if temp != "" {
		if _, err := file.WriteString(temp); err != nil {
			fmt.Printf("Error creating file %s ::- %v\n", fileName, err)
			return
		}
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

	fmt.Printf("Generating file %s ... \n", fileName)

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
