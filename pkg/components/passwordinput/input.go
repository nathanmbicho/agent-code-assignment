package passwordinput

import "C"
import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nathanmbicho/agent-code-assignment/pkg/ui"
	"os"
	"os/user"
	"strings"
)

var (
	targetPath string
	force      bool
)

type deleteResult struct {
	err     error
	message string
}

// state
type state int

const (
	confirmationState state = iota
	passwordState
	processingState
	completedState
	cancelledState
	errorState
)

type Model struct {
	state      state
	targetPath string
	isDir      bool
	textInput  textinput.Model
	password   string
	err        error
	message    string
}

// InitialPasswordInputModel - initialize the model
func InitialPasswordInputModel(path string, isDirectory bool) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your password"
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	ti.CharLimit = 256

	return Model{
		state:      confirmationState,
		targetPath: path,
		isDir:      isDirectory,
		textInput:  ti,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case confirmationState:
			switch msg.String() {
			case "y", "Y":
				m.state = passwordState
				m.textInput.Focus()
				return m, textinput.Blink
			case "n", "N":
				m.state = cancelledState
				m.message = "operation cancelled by user."
				return m, tea.Quit
			case "ctrl+c", "esc":
				return m, tea.Quit
			}

		case passwordState:
			switch msg.String() {
			case "enter":
				m.password = m.textInput.Value()
				if m.password == "" {
					m.err = fmt.Errorf("password cannot be empty")
					return m, nil
				}

				// Authenticate and delete
				m.state = processingState
				return m, tea.Batch(
					func() tea.Msg { return authenticateAndDelete(m.targetPath, m.password) },
				)
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.state = confirmationState
				m.textInput.SetValue("")
				m.textInput.Blur()
				return m, nil
			}

			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd

		case processingState, completedState, cancelledState, errorState:
			switch msg.String() {
			case "ctrl+c", "enter", "esc":
				return m, tea.Quit
			}
		}

	case deleteResult:
		if msg.err != nil {
			m.state = errorState
			m.err = msg.err
		} else {
			m.state = completedState
			m.message = msg.message
		}
		return m, nil
	}

	return m, cmd
}

// View implements tea.Model
func (m Model) View() string {
	var s strings.Builder
	s.WriteString(ui.RenderHeader("Secure Delete") + "\n\n")

	itemType := "file"
	if m.isDir {
		itemType = "directory"
	}
	s.WriteString(fmt.Sprintf("path %s: %s\n\n", itemType, m.targetPath))

	switch m.state {
	case confirmationState:
		if m.isDir {
			s.WriteString(ui.RenderInfo("warning: This is a directory!") + "\n")
			s.WriteString("all contents will be permanently deleted.\n\n")
		}

		s.WriteString("are you sure you want to delete this " + itemType + "?" + "\n\n")
		s.WriteString(ui.RenderInfo("press [Y] to continue, [N] to cancel\n press [Esc] or [Ctrl+C] to quit"))

	case passwordState:
		s.WriteString(ui.RenderInfo("Authentication Required") + "\n\n")

		currentUser, _ := user.Current()
		if currentUser != nil {
			s.WriteString(fmt.Sprintf("User: %s\n", currentUser.Username))
		}

		s.WriteString(ui.RenderHeader("Enter your password to proceed:\n\n"))
		s.WriteString(m.textInput.View() + "\n\n")

		if m.err != nil {
			s.WriteString(ui.RenderError(m.err.Error()) + "\n\n")
		}

		s.WriteString(ui.RenderInfo("Press [Enter] to confirm, [Esc] to go back, [Ctrl+C] to quit"))

	case processingState:
		s.WriteString(ui.TextStyle.Render("Processing...") + "\n\n")
		s.WriteString("Authenticating and deleting " + itemType + "...\n")

	case completedState:
		s.WriteString(ui.RenderSuccess("Success!") + "\n\n")
		s.WriteString(m.message + "\n\n")
		s.WriteString(ui.RenderInfo("Press [Enter] or [Esc] to exit"))

	case cancelledState:
		s.WriteString(ui.RenderError("Cancelled") + "\n\n")
		s.WriteString(m.message + "\n\n")
		s.WriteString(ui.RenderInfo("Press [Enter] or [Esc] to exit"))

	case errorState:
		s.WriteString(ui.RenderError("Error") + "\n\n")
		s.WriteString(m.err.Error() + "\n\n")
		s.WriteString(ui.RenderInfo("Press [Enter] or [Esc] to exit"))
	}

	return s.String()
}

// authenticate and delete function
func authenticateAndDelete(path, password string) deleteResult {
	// verify password (simplified for demo)
	if err := verifyPassword(password); err != nil {
		return deleteResult{err: fmt.Errorf("authentication failed: %w", err)}
	}

	// perform deletion
	info, err := os.Stat(path)
	if err != nil {
		return deleteResult{err: fmt.Errorf("error accessing path: %w", err)}
	}

	if info.IsDir() {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}

	if err != nil {
		return deleteResult{err: fmt.Errorf("deletion failed: %w", err)}
	}

	itemType := "file"
	if info.IsDir() {
		itemType = "directory"
	}

	return deleteResult{message: fmt.Sprintf("Successfully deleted %s: %s", itemType, path)}
}

func verifyPassword(password string) error {
	if len(password) == 0 {
		return fmt.Errorf("password cannot be empty")
	}

	// get current user
	_, err := user.Current()
	if err != nil {
		return fmt.Errorf("unable to get current user: %w", err)
	}

	// authenticate using system-specific auth (to work on)
	return nil
}
