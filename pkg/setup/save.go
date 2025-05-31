package setup

import tea "github.com/charmbracelet/bubbletea"

// saves login information to encrypted file
func (s Setup) save_login() tea.Cmd
func (s Setup) saving(msg tea.Msg) (tea.Model, tea.Cmd)
