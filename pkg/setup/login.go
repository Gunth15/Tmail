package setup

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Login struct {
	input []textinput.Model
	focus int
}

func InitLoginModel() Login {
	inputs := make([]textinput.Model, 3)

	for i := range inputs {
		inputs[i] = textinput.New()
		inputs[i].Width = 50
		inputs[i].PromptStyle = blurredStyle
		inputs[i].TextStyle = blurredStyle
		switch i {
		case 0:
			inputs[i].ShowSuggestions = true
			inputs[i].Focus()
			inputs[i].SetSuggestions([]string{"gmail", "outlook", "yahoo"})
			inputs[i].Placeholder = "Email Service?"
			inputs[i].PromptStyle = focusedStyle
			inputs[i].TextStyle = focusedStyle

		case 1:
			inputs[i].Placeholder = "Username"

		case 2:
			inputs[i].Placeholder = "Password"
			inputs[i].EchoMode = textinput.EchoPassword
			inputs[i].EchoCharacter = '.'
		}
	}

	return Login{
		input: inputs,
		focus: 0,
	}
}

func (l Login) Init() tea.Cmd {
	return textinput.Blink
}

func (s Setup) get_login(msg tea.Msg) (tea.Model, tea.Cmd) {
	l := &s.login
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			for i := range l.input {
				if l.input[i].ShowSuggestions && i == l.focus {
					l.input[i].SetValue(l.input[i].CurrentSuggestion())
					l.input[i].CursorEnd()
				}
			}
		case "ctrl+c", "esc":
			return s, tea.Quit
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			str := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so authenticate their login.
			if str == "enter" && l.focus == len(l.input) {
				s.state = AUTHENTICATE
				s.waiting = s.waiting.Wait("Verifying Account", func() (tea.Model, tea.Cmd) {
					s.state = SAVE
					s.waiting.reason = "Succes"
					return s, nil
				}, func() (tea.Model, tea.Cmd) {
					s.state = LOGIN
					s.login = InitLoginModel()
					return s, nil
				})
				return s, tea.Batch(attempt_login(s.login.input[1].Value(), s.login.input[2].Value()), s.waiting.Init())
			}

			// Cycle indexes
			if str == "up" || str == "shift+tab" {
				l.focus--
			} else {
				l.focus++
			}
			if l.focus > len(l.input) {
				l.focus = 0
			} else if l.focus < 0 {
				l.focus = len(l.input)
			}

			// change focus
			cmds := make([]tea.Cmd, len(l.input))
			for i := range l.input {
				if i == l.focus {
					cmds[i] = l.input[i].Focus()
					l.input[i].PromptStyle = focusedStyle
					l.input[i].TextStyle = focusedStyle
					continue
				}
				l.input[i].Blur()
				l.input[i].PromptStyle = blurredStyle
				l.input[i].TextStyle = blurredStyle
			}

			return s, tea.Batch(cmds...)
		}
	}
	cmds := make([]tea.Cmd, len(l.input))
	for i := range cmds {
		l.input[i], cmds[i] = l.input[i].Update(msg)
	}
	return s, tea.Batch(cmds...)
}

func (l Login) View() string {
	buff := ""
	for _, input := range l.input {
		buff += input.View() + "\n"
	}
	if len(l.input) == l.focus {
		buff += focusedStyle.Render(Button)
	} else {
		buff += blurredStyle.Render(Button)
	}
	return buff
}
