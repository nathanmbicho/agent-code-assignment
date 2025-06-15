package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	dirPath    string
	showHidden bool
)

var readDirCmd = &cobra.Command{
	Use:   "read",
	Short: "Read directory and list its content",
	Long:  `Read directory and list its content, both its files and other directories in tree like structure.`,
	Run:   readDirectory,
}

func init() {
	rootCmd.AddCommand(readDirCmd)

	readDirCmd.Flags().StringVarP(&dirPath, "path", "p", ".", "path name with current directory as default")
	readDirCmd.Flags().BoolVarP(&showHidden, "all", "a", false, "show hidden files and directories")
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

	fmt.Printf("absolute path %s\n", path)

	// print directory details
	err = printDirectory(path, "")
	if err != nil {
		fmt.Printf("error getting dir contents %v \n", path)
		return
	}

	fmt.Printf("\n")
}

// print directory contents
func printDirectory(dirPath, prefix string) error {
	// read directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	// filter hidden files if not showing them
	var filteredEntries []os.DirEntry
	for _, entry := range entries {
		if !showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		filteredEntries = append(filteredEntries, entry)
	}

	// sort entries - directories first, then files
	sort.Slice(filteredEntries, func(i, j int) bool {
		if filteredEntries[i].IsDir() && !filteredEntries[j].IsDir() {
			return true
		}
		if !filteredEntries[i].IsDir() && filteredEntries[j].IsDir() {
			return false
		}
		return filteredEntries[i].Name() < filteredEntries[j].Name()
	})

	// loop entries creating a tree like directory design
	for i, entry := range filteredEntries {
		isLast := i == len(filteredEntries)-1

		var connector, childPrefix string
		if isLast {
			connector = "└── "
			childPrefix = prefix + "    "
		} else {
			connector = "├── "
			childPrefix = prefix + "│   "
		}

		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		fmt.Printf("%s%s%s\n", prefix, connector, name)

		// recursively print subdirectories
		if entry.IsDir() {
			subPath := filepath.Join(dirPath, entry.Name())
			err := printDirectory(subPath, childPrefix)
			if err != nil {
				fmt.Printf("%s%s[Error: %v]\n", childPrefix, "├── ", err)
			}
		}
	}

	return nil
}
