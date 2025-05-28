package setup

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

type Setup struct {
	Type     textinput.Model
	Username textinput.Model
	Password textinput.Model
	Cursor   cursor.Mode
}

func InitSetupModel() Setup {
	setup := Setup{
		Type:     textinput.New(),
		Username: textinput.New(),
		Password: textinput.New(),
		Cursor:   0,
	}
	setup.Type.ShowSuggestions = true

	email_suggestions := make([]string, 3)
	email_suggestions = append(email_suggestions, "Gmail", "Outlook", "Other")

	setup.Type.SetSuggestions(email_suggestions)

	setup.Type.Placeholder = "What email Service would you like to use with Tmail"

	setup.Username.Placeholder = "Username"

	setup.Password.Placeholder = "Password"
	setup.Password.EchoCharacter = '.'

	return setup
}

func (s Setup) Init() tea.Cmd {
	return textinput.Blink
}

func (s Setup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Default
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return s, tea.Quit
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			str := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if str == "enter" && s.Cursor == 4 {
				return s, tea.Quit
			}

			// Cycle indexes
			if str == "up" || str == "shift+tab" {
				s.Cursor--
			} else {
				s.Cursor++
			}

			if s.Cursor > 3 {
				s.Cursor = 0
			} else if s.Cursor < 0 {
				s.Cursor = 3
			}

		}
	}

	cmds := make([]tea.Cmd, 3)
	_, cmds[0] = s.Username.Update(msg)
	_, cmds[1] = s.Password.Update(msg)
	_, cmds[2] = s.Type.Update(msg)
	return s, tea.Batch(cmds...)
}

func (s Setup) View() string {
	username := s.Username.View() + "\n"
	password := s.Password.View() + "\n"
	email_type := s.Type.View() + "\n"

	return email_type + username + password
}
