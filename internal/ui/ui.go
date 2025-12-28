package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knbr13/glow-net/internal/scanner"
	"github.com/knbr13/glow-net/internal/state"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tickMsg time.Time

type Model struct {
	table     table.Model
	state     *state.AppState
	bwScan    *scanner.BandwidthScanner
	width     int
	height    int
	err       error
	sortBy    string
	filter    string
	filtering bool
}

func NewModel(s *state.AppState) Model {
	columns := []table.Column{
		{Title: "PID", Width: 7},
		{Title: "Process", Width: 15},
		{Title: "Local Address", Width: 25},
		{Title: "Remote Address", Width: 25},
		{Title: "Remote Host", Width: 30},
		{Title: "Country", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	sStyle := table.DefaultStyles()
	sStyle.Header = sStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	sStyle.Selected = sStyle.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(sStyle)

	return Model{
		table:  t,
		state:  s,
		bwScan: &scanner.BandwidthScanner{},
		sortBy: "pid",
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.filtering {
			switch msg.String() {
			case "enter", "esc":
				m.filtering = false
				return m, nil
			case "backspace":
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
				}
				m.updateTable()
				return m, nil
			default:
				if len(msg.String()) == 1 {
					m.filter += msg.String()
					m.updateTable()
				}
				return m, nil
			}
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "f":
			m.filtering = true
			m.filter = ""
			return m, nil
		case "s":
			// Cycle sorting
			switch m.sortBy {
			case "pid":
				m.sortBy = "name"
			case "name":
				m.sortBy = "host"
			default:
				m.sortBy = "pid"
			}
			m.updateTable()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetHeight(m.height - 10)
		return m, nil
	case tickMsg:
		err := m.state.Refresh(m.bwScan)
		if err != nil {
			m.err = err
		}
		m.updateTable()
		return m, tick()
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *Model) updateTable() {
	rowsData := m.state.GetTableRows()

	// Filter
	if m.filter != "" {
		var filtered []state.RowData
		f := strings.ToLower(m.filter)
		for _, r := range rowsData {
			if strings.Contains(strings.ToLower(r.ProcessName), f) ||
				strings.Contains(strings.ToLower(r.RemoteHost), f) ||
				strings.Contains(strings.ToLower(r.RemoteAddr), f) {
				filtered = append(filtered, r)
			}
		}
		rowsData = filtered
	}

	// Sort
	sort.Slice(rowsData, func(i, j int) bool {
		switch m.sortBy {
		case "name":
			return rowsData[i].ProcessName < rowsData[j].ProcessName
		case "host":
			return rowsData[i].RemoteHost < rowsData[j].RemoteHost
		default:
			return rowsData[i].PID < rowsData[j].PID
		}
	})

	var rows []table.Row
	for _, r := range rowsData {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", r.PID),
			r.ProcessName,
			r.LocalAddr,
			r.RemoteAddr,
			r.RemoteHost,
			r.Country,
		})
	}
	m.table.SetRows(rows)
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	header := m.renderHeader()

	footer := "\n  q: quit • j/k: up/down • s: sort • f: filter"
	if m.filtering {
		footer = "\n  FILTER: " + m.filter + "█ (esc to close)"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		baseStyle.Render(m.table.View()),
		footer,
	)
}

func (m Model) renderHeader() string {
	m.state.RLock()
	stats := m.state.GlobalStats
	m.state.RUnlock()

	download := formatBytes(stats.DownloadSpeed) + "/s"
	upload := formatBytes(stats.UploadSpeed) + "/s"

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render(" GlowNet ")

	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Network Monitoring Tool")

	statsStr := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Render(fmt.Sprintf("Download: %s", download)) +
		"  " +
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Render(fmt.Sprintf("Upload: %s", upload))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, desc, "  ", statsStr)
}

func formatBytes(b float64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%.2f B", b)
	}
	div, exp := float64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", b/div, "KMGTPE"[exp])
}
