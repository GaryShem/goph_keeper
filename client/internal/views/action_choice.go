package views

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	listTitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	listItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	listSelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	listPaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	listHelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	listQuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := listItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return listSelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

type ActionChoice struct {
	list     list.Model
	choice   string
	quitting bool
	header   string
	footer   string
}

func NewActionChoice(header, footer string) ActionChoice {
	items := []list.Item{
		item("Upload Text"),
		item("Upload Card"),
		item("Upload Binary"),
		item("Download Text"),
		item("Download Card"),
		item("Download Binary"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Choose action:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = listTitleStyle
	l.Styles.PaginationStyle = listPaginationStyle
	l.Styles.HelpStyle = listHelpStyle

	m := ActionChoice{
		header: header,
		footer: footer,
		list:   l,
	}

	return m
}

func (m ActionChoice) Init() tea.Cmd {
	return nil
}

func (m ActionChoice) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			switch m.choice {
			case "Upload Text":
				return NewUploadTextModel(), nil
			case "Upload Card":
				return NewUploadCardModel(), nil
			case "Upload Binary":
				return NewUploadBinaryModel(), nil
			case "Download Text":
				return NewDownloadTextModel(), nil
			case "Download Card":
				return NewDownloadCardModel(), nil
			case "Download Binary":
				return NewDownloadBinaryModel(), nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ActionChoice) View() string {
	if m.quitting {
		return listQuitTextStyle.Render("Quitting～")
	}
	return fmt.Sprintf("%s\n\n", m.header) + m.list.View() + fmt.Sprintf("%s\n\n", m.footer)
}
