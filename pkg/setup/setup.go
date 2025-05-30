package setup

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/emersion/go-imap/v2/imapclient"
)

type (
	SetupState uint
	EmailType  string
)

const Button = "[ submit ]"

const (
	GMAIL   EmailType = "gmail"
	OUTLOOK EmailType = "outlook"
	YAHOO   EmailType = "yahoo"
	Other   EmailType = "other"
)

const (
	EMAILTYPE SetupState = iota
	LOGIN
	AUTHENTICATE
	SAVE
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3489cf"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type Waiting struct {
	spinner spinner.Model
	sending bool
	reason  string
	quit    bool
}
type Setup struct {
	input   []textinput.Model
	waiting Waiting
	focus   int
	state   SetupState
}

func InitSetupModel() Setup {
	inputs := make([]textinput.Model, 3)

	inputs[0] = textinput.New()
	inputs[0].ShowSuggestions = true
	inputs[0].Focus()
	inputs[0].SetSuggestions([]string{"gmail", "outlook", "yahoo"})
	inputs[0].Placeholder = "Email Service?"
	inputs[0].PromptStyle = focusedStyle
	inputs[0].TextStyle = focusedStyle
	inputs[0].Width = 50

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Username"
	inputs[1].Width = 50

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Password"
	inputs[2].EchoMode = textinput.EchoPassword
	inputs[2].EchoCharacter = '.'
	inputs[2].Width = 50

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = focusedStyle

	return Setup{
		input: inputs,
		focus: 0,
		state: LOGIN,
		waiting: Waiting{
			spinner: spin,
			reason:  "",
		},
	}
}

func (s Setup) Init() tea.Cmd {
	return textinput.Blink
}

func (s Setup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch s.state {
	case LOGIN:
		return s.update_login(msg)
	case AUTHENTICATE:
		return s.update_auth(msg)
	case SAVE:
		return s, tea.Quit
	}
	return s, tea.Quit
}

func (s Setup) View() string {
	switch s.state {
	case LOGIN:
		buff := ""
		for _, input := range s.input {
			buff += input.View() + "\n"
		}
		if len(s.input) == s.focus {
			buff += focusedStyle.Render(Button)
		} else {
			buff += blurredStyle.Render(Button)
		}
		return buff
	case AUTHENTICATE:
		if s.waiting.reason != "" {
			return focusedStyle.Render(s.waiting.reason)
		}
		return s.waiting.spinner.View() + "Verifying Account"
	}
	return "Error, invalid state"
}

func (s Setup) update_login(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
				s.state = AUTHENTICATE

				return s, s.attempt_login(s.input[1].Value(), s.input[2].Value())
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

func (s Setup) update_auth(msg tea.Msg) (tea.Model, tea.Cmd) {
	waiting := &s.waiting

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctr+c":
			return s, tea.Quit
		case "enter":
			if s.waiting.reason != "" {
				var cmd tea.Cmd
				waiting.spinner, cmd = waiting.spinner.Update(msg)
				return s, cmd
			}
			return s, nil

		default:
			return s, nil
		}
	case error:
		if msg != nil {
			waiting.reason = "Error: " + msg.Error()
		} else {
			waiting.reason = "Succes"
		}
		waiting.sending = false
		return s, nil
	default:
		var cmd tea.Cmd
		waiting.spinner, cmd = waiting.spinner.Update(msg)
		return s, cmd
	}
}

func (s Setup) attempt_login(user string, pass string) tea.Cmd {
	return func() tea.Msg {
		client, err := imapclient.DialTLS("imap.gmail.com:993", nil)
		if err != nil {
			return err
		}

		if err = client.Login(user, pass).Wait(); err != nil {
			return err
		}

		if err = client.Logout().Wait(); err != nil {
			return err
		}
		return nil
	}
}
