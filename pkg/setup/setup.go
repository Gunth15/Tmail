package setup

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	SetupState uint
	EmailType  uint
)

const (
	GMAIL EmailType = iota
	OUTLOOK
	Other
)

const (
	EMAILTYPE SetupState = iota
	USERNAME
	PASSWORD
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type Setup struct {
	input  []textinput.Model
	cursor cursor.Mode
	focus  int
}

func InitSetupModel() Setup {
	inputs := make([]textinput.Model, 3)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Username"
	inputs[0].Focus()
	inputs[0].PromptStyle = focusedStyle
	inputs[0].TextStyle = focusedStyle

	inputs[1] = textinput.New()
	inputs[1].ShowSuggestions = true
	inputs[1].SetSuggestions([]string{"Gmail", "Outlook", "Other"})
	inputs[1].Placeholder = "Email Service?"

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Password"
	inputs[2].EchoMode = textinput.EchoPassword
	inputs[2].EchoCharacter = '.'

	return Setup{
		input: inputs,
		focus: 0,
	}
}

func (s Setup) Init() tea.Cmd {
	return textinput.Blink
}

func (s Setup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Default
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			for i := range s.input {
				if s.input[i].ShowSuggestions && i == s.focus {
					s.input[i].SetValue(s.input[i].CurrentSuggestion())
					s.input[i].CursorEnd()
				}
			}
		case "ctrl+c", "esc":
			return s, tea.Quit
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			str := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if str == "enter" && s.focus == len(s.input) {
				return s, tea.Quit
			}

			// Cycle indexes
			if str == "up" || str == "shift+tab" {
				s.focus--
			} else {
				s.focus++
			}

			if s.focus > len(s.input) {
				s.focus = 0
			} else if s.focus < 0 {
				s.focus = len(s.input)
			}

			cmds := make([]tea.Cmd, len(s.input))
			for i := range s.input {
				if i == s.focus {
					cmds[i] = s.input[i].Focus()
					s.input[i].PromptStyle = focusedStyle
					s.input[i].TextStyle = focusedStyle
					continue
				}
				s.input[i].Blur()
				s.input[i].PromptStyle = blurredStyle
				s.input[i].TextStyle = blurredStyle
			}

			return s, tea.Batch(cmds...)
		}
	}
	cmds := make([]tea.Cmd, len(s.input))
	for i := range cmds {
		s.input[i], cmds[i] = s.input[i].Update(msg)
	}
	return s, tea.Batch(cmds...)
}

func (s Setup) View() string {
	buff := ""
	for _, input := range s.input {
		buff += input.View() + "\n"
	}

	return buff
}
