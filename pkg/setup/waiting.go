package setup

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Waiting struct {
	spinner  spinner.Model
	sending  bool
	reason   string
	error    bool
	callback func() (tea.Model, tea.Cmd)
}

func InitWaitingModel() Waiting {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = focusedStyle
	return Waiting{
		reason:  "",
		spinner: spin,
	}
}

func (w Waiting) Init() tea.Cmd {
	return w.spinner.Tick
}

func (s Setup) verifing(msg tea.Msg) (tea.Model, tea.Cmd) {
	w := &s.waiting
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctr+c":
			return s, tea.Quit

		case "enter":
			if w.error {
				s.state = LOGIN
				s.login = InitLoginModel()
				return s, nil
			}
			return s, nil

		default:
			return s, nil
		}

	case error:
		w.sending = false
		if msg != nil {
			w.reason = "(Error)" + msg.Error() + "\n\n(Press enter to retry)"
			w.error = true
			return s, nil
		} else {
			return w.callback()
		}

	default:
		var cmd tea.Cmd
		w.spinner, cmd = w.spinner.Update(msg)
		return s, cmd
	}
}

func (w Waiting) View() string {
	if !w.sending {
		return focusedStyle.Render(w.reason)
	}
	return w.spinner.View() + focusedStyle.Render(w.reason)
}

func (w *Waiting) Wait(reason string, cb func() (tea.Model, tea.Cmd)) {
	w.reason = reason
	w.sending = true
	w.callback = cb
}
