package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

// NOTE: May need to explore possible custom chracterset readers which may make parsing more effiucient for future customizations

type Mailboxes struct {
	cursor    int
	mailboxes []string
}
type MainView struct {
	msg     string
	sidebar Mailboxes
	isMail  bool
}

func (v MainView) Init() tea.Cmd {
	return getMailboxes
}

func (v MainView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Mailboxes:
		v.isMail = true
		v.sidebar = msg
		return v, tea.ClearScreen
	case string:
		v.msg = msg
		return v, tea.ClearScreen
	case tea.KeyMsg:
		////////////
		switch msg.Type {
		case tea.KeyCtrlC:
			return v, tea.Quit
		case tea.KeyUp:
			v.sidebar.cursor = v.sidebar.cursor % (len(v.sidebar.mailboxes) - 1)
			return v, nil
		case tea.KeyDown:
			v.sidebar.cursor = (v.sidebar.cursor + 1) % len(v.sidebar.mailboxes)
			return v, nil
		default:
			v.msg = "Coming soon"
			return v, nil
		}
		////////////
	}
	return v, nil
}

func (v MainView) View() string {
	if !v.isMail {
		return v.msg
	}

	alignmanet := lipgloss.NewStyle().Align(lipgloss.Left).Border(lipgloss.RoundedBorder(), false, true).Width(8).Height(2).Margin(8)
	return list.New(v.sidebar.mailboxes).ItemStyleFunc(func(items list.Items, index int) lipgloss.Style {
		if index == v.sidebar.cursor {
			return lipgloss.NewStyle().Background(lipgloss.Color("#FFFF00")).Inherit(alignmanet)
		}
		return lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Inherit(alignmanet)
	}).Enumerator(func(items list.Items, index int) string { return "" }).String()
}

func getMailboxes() tea.Msg {
	client, err := imapclient.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return err.Error()
	}

	if err = client.Login(testUsername, testPass).Wait(); err != nil {
		return err.Error()
	}

	// ReturnStatus requires server support for IMAP4rev2 or LIST-STATUS
	listCmd := client.List("", "%", &imap.ListOptions{
		ReturnStatus: &imap.StatusOptions{
			NumMessages: true,
			NumUnseen:   true,
		},
	})

	mboxes, err := listCmd.Collect()
	if err != nil {
		return err.Error()
	}

	mailboxes := make([]string, 0)
	for _, mbox := range mboxes {

		if mbox == nil {
			break
		}
		mailboxes = append(mailboxes, mbox.Mailbox)
	}

	if err = client.Logout().Wait(); err != nil {
		return err.Error()
	}

	return Mailboxes{
		cursor:    0,
		mailboxes: mailboxes,
	}
}

// TODO: Make a nice error interface

func main() {
	if _, err := tea.NewProgram(MainView{msg: "Waiting on Server"}).Run(); err != nil {
		fmt.Println("Eww error", err)
		os.Exit(1)
	}
}
