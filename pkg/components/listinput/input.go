package listinput

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nathanmbicho/agent-code-assignment/pkg/ui"
)

var (
	result string
	valid  bool
)

type Selection struct {
	Choice string
	Quit   bool
}

func (s *Selection) Update(choice string) {
	s.Choice = choice
}

func (s *Selection) QuitCmd() {
	s.Quit = true
}

type Model struct {
	err          error
	cursor       int
	fileName     string
	choices      []string
	selected     map[int]struct{}
	choice       *Selection
	header       string
	validateFunc func(string, string) (string, bool, error)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func InitialListInputModel(choices []string, fileName string, selection *Selection, header string, validateFunc func(string, string) (string, bool, error)) Model {
	return Model{
		err:          nil,
		choices:      choices,
		fileName:     fileName,
		selected:     make(map[int]struct{}),
		header:       ui.RenderHeader(header),
		choice:       selection,
		validateFunc: validateFunc,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			//validate options
			result, valid, m.err = m.validateFunc(m.fileName, m.choices[m.cursor])
			if valid {
				if len(m.selected) == 1 {
					m.selected = make(map[int]struct{})
				}

				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}

				return m, tea.Quit
			}

			m.selected[m.cursor] = struct{}{}
			return m, nil
		case "y":
			if len(m.selected) == 1 {
				return m, tea.Quit
			}
		case "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	view := m.header + "\n"

	// handle selected choice
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ui.SuccessStyle2.Render(">")
		}

		view += fmt.Sprintf("%s [%d] %s\n\n", cursor, i+1, ui.TextStyle.Render(choice))

	}

	if m.err != nil {
		view += ui.RenderError(fmt.Sprintf("%s. Please try again!", m.err.Error()))
	}

	view += fmt.Sprintf("\n%s\n", ui.RenderInfo("(press enter to confirm choice, esc to quit)"))

	// show code if default is selected
	if _, ok := m.selected[m.cursor]; ok && m.cursor == 0 {
		view += fmt.Sprintf("\n%s\n", ui.RenderCode(fmt.Sprintf("\n%+v\n", result)))
	}

	return fmt.Sprintf("%s\n", view)
}
