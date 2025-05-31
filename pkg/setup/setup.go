package setup

import (
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

type Setup struct {
	login   Login
	waiting Waiting
	state   SetupState
}

func InitSetupModel() Setup {
	return Setup{
		state:   LOGIN,
		login:   InitLoginModel(),
		waiting: InitWaitingModel(),
	}
}

func (s Setup) Init() tea.Cmd {
	return s.login.Init()
}

func (s Setup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch s.state {
	case LOGIN:
		return s.get_login(msg)
	case AUTHENTICATE:
		return s.verifing(msg)
	case SAVE:
		return s, tea.Quit
	default:
		return s, tea.Quit
	}
}

func (s Setup) View() string {
	switch s.state {
	case LOGIN:
		return s.login.View()
	case AUTHENTICATE:
		return s.waiting.View()
	case SAVE:
		return s.waiting.View()
	default:
		return "Error, invalid state"
	}
}

func attempt_login(user string, pass string) tea.Cmd {
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
