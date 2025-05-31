package setup

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type WaitingCallBack func() (tea.Model, tea.Cmd)

type Waiting struct {
	spinner    spinner.Model
	sending    bool
	reason     string
	error      bool
	success_cb WaitingCallBack
	error_cb   WaitingCallBack
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

// WARNING: Set the error and success callback before use
func (w Waiting) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctr+c":
			return w, tea.Quit

		case "enter":
			if w.error {
				return w.error_cb()
			}
			return w, nil

		default:
			return w, nil
		}

	case error:
		w.sending = false
		if msg != nil {
			w.reason = "(Error)" + msg.Error() + "\n\n(Press enter to retry)"
			w.error = true
			return w, nil
		} else {
			return w.success_cb()
		}

	default:
		var cmd tea.Cmd
		w.spinner, cmd = w.spinner.Update(msg)
		return w, cmd
	}
}

func (w Waiting) View() string {
	if !w.sending {
		return focusedStyle.Render(w.reason)
	}
	return w.spinner.View() + focusedStyle.Render(w.reason)
}

func (w Waiting) Wait(reason string, success WaitingCallBack, error WaitingCallBack) Waiting {
	w.reason = reason
	w.sending = true
	w.success_cb = success
	w.error_cb = error
	return w
}
