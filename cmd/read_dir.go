package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sort"
)

var dirPath string

var readDirCmd = &cobra.Command{
	Use:   "read",
	Short: "Read directory and list its content",
	Long:  `Read directory and list its content, both its files and other directories in tree like structure.`,
	Run:   readDirectory,
}

func init() {
	rootCmd.AddCommand(readDirCmd)

	readDirCmd.Flags().StringVarP(&dirPath, "path", "p", ".", "path name with current directory as default")
	readDirCmd.MarkFlagRequired("file")
}

func readDirectory(cmd *cobra.Command, args []string) {
	// check if the path exists
	fileInfo, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		fmt.Printf("error - path %s does not exist \n", dirPath)
		return
	}

	// check if the path is a valid directory
	if !fileInfo.IsDir() {
		fmt.Printf("error - path %s is an ivalid directory \n", dirPath)
		return
	}

	path, err := filepath.Abs(dirPath)
	if err != nil {
		fmt.Printf("error getting path %v \n", path)
		return
	}

	// print directory details
	err = printDirectory(path)
	if err != nil {
		fmt.Printf("error getting dir contents %v \n", path)
		return
	}

	fmt.Printf("absolute path %s\n", path)
}

// print directory contents
func printDirectory(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// sort entries
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() && !entries[j].IsDir() {
			return true
		}

		if !entries[i].IsDir() && entries[j].IsDir() {
			return false
		}

		return entries[i].Name() < entries[j].Name()
	})

	fmt.Printf("found %d entries in %s\n", len(entries), dirPath)

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("DIR:: \t%s\n", entry.Name())
		} else {
			fmt.Printf("FILE:: \t%s\n", entry.Name())
		}
	}

	return nil
}
