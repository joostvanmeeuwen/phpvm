package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joostvanmeeuwen/phpvm/internal/php"
)

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	focusedBorderStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))

	blurredBorderStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	infoTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	infoValueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("246"))

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
)

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

type App struct {
	phpManager   *php.Manager
	list         list.Model
	details      viewport.Model
	selectedItem Item
	width        int
	height       int
	focused      int // 0: list, 1: details
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

	detailsView := viewport.New(0, 0)

	var selectedItem Item
	if len(items) > 0 {
		selectedItem = items[0].(Item)
	}
	return &App{
		phpManager:   manager,
		list:         l,
		details:      detailsView,
		selectedItem: selectedItem,
		focused:      0,
	}
}

func (a *App) Start() error {
	p := tea.NewProgram(a, tea.WithAltScreen())
	_, err := p.Run()

	return err
}

func (a *App) Init() tea.Cmd {
	a.updateDetails()
	return nil
}

func (a *App) updateDetails() {
	if i, ok := a.list.SelectedItem().(Item); ok {
		a.selectedItem = i

		var detailsContent strings.Builder

		detailsContent.WriteString(infoTitleStyle.Render("PHP Version Details"))
		detailsContent.WriteString("\n\n")

		detailsContent.WriteString(infoTitleStyle.Render("Version: "))
		detailsContent.WriteString(infoValueStyle.Render(a.selectedItem.version.Version))
		detailsContent.WriteString("\n\n")

		detailsContent.WriteString(infoTitleStyle.Render("Path: "))
		detailsContent.WriteString(infoValueStyle.Render(a.selectedItem.version.Path))
		detailsContent.WriteString("\n\n")

		detailsContent.WriteString(infoTitleStyle.Render("Status: "))
		if a.selectedItem.version.Active {
			detailsContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("Active"))
		} else {
			detailsContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Inactive"))
		}

		a.details.SetContent(detailsContent.String())
	}
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit

		case "tab":
			a.focused = (a.focused + 1) % 2
			return a, nil

		case "enter":
			return a, nil
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		columnWidth := (a.width / 2) - 5

		a.list.SetWidth(columnWidth)
		a.list.SetHeight(a.height - 12)

		a.details.Width = columnWidth
		a.details.Height = a.height - 10
	}

	if a.focused == 0 {
		newList, cmd := a.list.Update(msg)
		a.list = newList
		cmds = append(cmds, cmd)

		if a.list.SelectedItem() != nil {
			a.updateDetails()
		}
	} else {
		newViewport, cmd := a.details.Update(msg)
		a.details = newViewport
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) View() string {
	var leftStyle, rightStyle lipgloss.Style

	if a.focused == 0 {
		leftStyle = focusedBorderStyle.Copy().Padding(1)
		rightStyle = blurredBorderStyle.Copy().Padding(1)
	} else {
		leftStyle = blurredBorderStyle.Copy().Padding(1)
		rightStyle = focusedBorderStyle.Copy().Padding(1)
	}

	leftWidth := a.width/2 - 4
	rightWidth := a.width/2 - 4

	listView := a.list.View()

	helpText := fmt.Sprintf("\n  %s", helpStyle.Render("tab: switch panels • enter: activate version • q: quit"))

	leftContent := leftStyle.Width(leftWidth).Render(listView + helpText)
	rightContent := rightStyle.Width(rightWidth).Render(a.details.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent)
}
