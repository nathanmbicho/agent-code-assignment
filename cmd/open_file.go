package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

var openFileName string

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

	fmt.Printf("opening file %s... \n", openFileName)

	err = exec.Command("open", path).Start()
	if err != nil {
		fmt.Printf("error opening file %s - %v\n", path, err)
		return
	}
}
