package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joostvanmeeuwen/phpvm/internal/php"
)

var (
	appStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1)

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginLeft(0)
)

type App struct {
	phpManager *php.Manager
	list       list.Model
	width      int
	height     int
}

type Item struct {
	version php.PHPVersion
}

func (i Item) FilterValue() string {
	return i.version.Version
}

func (i Item) Title() string {
	title := i.version.Version

	if i.version.Active {
		title += " (active)"
	}

	return title
}

func (i Item) Description() string {
	return i.version.Path
}

func NewApp() *App {
	manager := php.NewManager()
	versions := manager.GetVersions()

	items := make([]list.Item, len(versions))
	for i, v := range versions {
		items[i] = Item{version: v}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "PHP Versions"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	return &App{
		phpManager: manager,
		list:       l,
	}
}

func (a *App) Start() error {
	p := tea.NewProgram(a, tea.WithAltScreen())
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

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.list.SetWidth(msg.Width - 4) // Account for border padding
		a.list.SetHeight(msg.Height - 4)
	}

	var cmd tea.Cmd
	a.list, cmd = a.list.Update(msg)

	return a, cmd
}

func (a *App) View() string {
	return appStyle.Render(a.list.View())
}
