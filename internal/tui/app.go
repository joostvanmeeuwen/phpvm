package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	viewport viewport.Model
}

func NewApp() *App {
	return &App{
		viewport: viewport.New(80, 24),
	}
}

func (a *App) Start() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}

func (a *App) Init() tea.Cmd {
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a *App) View() string {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Render("PHP Version Manager!")
}
