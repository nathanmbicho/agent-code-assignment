package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nathanmbicho/agent-code-assignment/pkg/components/passwordinput"
	"github.com/nathanmbicho/agent-code-assignment/pkg/components/textinput"
	"github.com/nathanmbicho/agent-code-assignment/pkg/ui"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

type options struct {
	FileName *textinput.Output
}

// deleteFileCmd - delete an existing file or directory
var deleteFileCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing file or folder",
	Long:  `You can delete an existing file or directory of given valid path.`,
	Run:   deleteFile,
}

func init() {
	rootCmd.AddCommand(deleteFileCmd)
}

func deleteFile(cmd *cobra.Command, args []string) {
	//input command
	option := options{
		FileName: &textinput.Output{},
	}

	// handle program create, passing values
	tProgram := tea.NewProgram(textinput.InitialTextInputModel(
		option.FileName,
		"Enter directory or file name to delete ...",
		func(input string) (bool, error) {
			return validateSearchFile(input)
		},
	))

	// run bubbletea program
	if _, err := tProgram.Run(); err != nil {
		cobra.CheckErr(err)
		return
	}

	if option.FileName.Quit {
		fmt.Println("\n ❌Open file operation cancelled.")
	}

	fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("You are about to delete %s", option.FileName.Output)))

	absPath, isDir, err := validateDeleteFile(option.FileName.Output)
	if err != nil {
		cobra.CheckErr(err)
		return
	}

	// Start Bubble Tea program
	tProgram = tea.NewProgram(passwordinput.InitialPasswordInputModel(absPath, isDir),
		tea.WithAltScreen(),
	)

	if _, err := tProgram.Run(); err != nil {
		cobra.CheckErr(err)
		return
	}

}

func validateDeleteFile(targetPath string) (string, bool, error) {
	if targetPath == "" {
		return "", false, fmt.Errorf("path flag is required")
	}

	// Resolve the absolute path
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return "", false, fmt.Errorf("error resolving path: %w", err)
	}

	// Check if the path exists
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return "", false, fmt.Errorf("path does not exist: %s", absPath)
	}
	if err != nil {
		return "", false, fmt.Errorf("error accessing path: %w", err)
	}

	return absPath, info.IsDir(), nil
}
