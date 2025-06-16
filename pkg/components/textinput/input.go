package textinput

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nathanmbicho/agent-code-assignment/pkg/ui"
)

type (
	errMsg error
)

type Output struct {
	Output string
	Quit   bool
}

func (o *Output) Update(value string) {
	o.Output = value
}

func (o *Output) QuitCmd() {
	o.Quit = true
}

type Model struct {
	textInput    textinput.Model
	err          error
	output       *Output
	header       string
	validateFunc func(string) (bool, error)
}

func InitialTextInputModel(output *Output, header string, validateFunc func(string) (bool, error)) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter something here..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 250

	return Model{
		textInput:    ti,
		err:          nil,
		output:       output,
		header:       ui.RenderHeader(header),
		validateFunc: validateFunc,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:

			input := m.textInput.Value()

			// VALIDATION HERE

			// if not validation
			if m.validateFunc == nil {
				m.output.Update(input)
				return m, tea.Quit
			}

			// validate input
			valid, err := m.validateFunc(input)
			if valid {
				m.output.Update(input)
				return m, tea.Quit
			}

			m.err = err
			m.textInput.SetValue("")
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	view := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		m.header,
		ui.InputStyle.Render(m.textInput.View()),
		ui.RenderInfo("(press enter to confirm, esc/ctrl+c to quit)"),
	) + "\n\n"

	if m.err != nil {
		view += ui.RenderError(fmt.Sprintf("%s\n Please try again!", m.err.Error()))
	}

	return view + "\n\n"
}
