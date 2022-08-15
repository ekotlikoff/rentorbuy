package data

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle     = focusedStyle.Copy()
	noStyle         = lipgloss.NewStyle()
	helpStyle       = blurredStyle.Copy()
	validateNumeric = func(s string) error {
		_, err := strconv.ParseFloat(s, 64)
		return err
	}
	focusedButton = focusedStyle.Copy().Render("[ Open chart ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Open chart"))
)

type model struct {
	scenario           Scenario
	savingsAfterMonths int
	focusIndex         int
	inputs             []textinput.Model
	openedChart        *os.File
}

func InitialModel() model {
	m := model{
		scenario:           defaultScenario,
		savingsAfterMonths: 8 * 12,
		inputs:             make([]textinput.Model, 6),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Prompt = "Rent > "
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.SetValue(fmt.Sprintf("%.0f", m.scenario.Rent))
			t.Validate = validateNumeric
		case 1:
			t.Prompt = "Cost > "
			t.CharLimit = 20
			t.SetValue(fmt.Sprintf("%.0f", m.scenario.House.Cost))
			t.Validate = validateNumeric
		case 2:
			t.Prompt = "Down payment proportion > "
			t.CharLimit = 4
			t.SetValue(fmt.Sprintf("%.2f", m.scenario.House.DownPaymentProportion))
			t.Validate = validateNumeric
		case 3:
			t.Prompt = "Monthly maintenance > "
			t.CharLimit = 20
			t.SetValue(fmt.Sprintf("%.0f", m.scenario.House.MaintenanceMonthly))
			t.Validate = validateNumeric
		case 4:
			t.Prompt = "Proportion of down payment invested if renting > "
			t.CharLimit = 20
			t.SetValue(fmt.Sprintf("%.2f", m.scenario.Assumptions.DownPaymentInvestedPropIfRenting))
			t.Validate = validateNumeric
		case 5:
			t.Prompt = "Display savings after months > "
			t.CharLimit = 10
			t.SetValue(fmt.Sprintf("%d", m.savingsAfterMonths))
			t.Validate = validateNumeric
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

type OpenedChartMsg *os.File

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case OpenedChartMsg:
		m.openedChart = OpenedChartMsg(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.openedChart != nil {
				os.Remove(m.openedChart.Name())
				m.openedChart = nil
			}
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down", "j", "k":
			s := msg.String()

			if m.focusIndex == len(m.inputs) && s == "enter" {
				if m.openedChart != nil {
					os.Remove(m.openedChart.Name())
					m.openedChart = nil
				}
				return m, tea.Batch(func() tea.Msg {
					f := m.scenario.GenerateChart()
					cmd := exec.Command("open", f.Name())
					cmd.Run()
					return OpenedChartMsg(f)
				})

			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" || s == "k" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		if m.inputs[i].Focused() && m.inputs[i].Err == nil {
			v, _ := strconv.ParseFloat(m.inputs[i].Value(), 64)
			switch i {
			case 0:
				m.scenario.Rent = v
			case 1:
				m.scenario.House.Cost = v
			case 2:
				m.scenario.House.DownPaymentProportion = v
			case 3:
				m.scenario.House.MaintenanceMonthly = v
			case 4:
				m.scenario.Assumptions.DownPaymentInvestedPropIfRenting = v
			}
		}

	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder
	d := m.scenario.run()
	breakEvenMonths, savingsAfter := -1, -1.0
	for i, v := range d {
		if breakEvenMonths != -1 && i > m.savingsAfterMonths {
			break
		}
		if v >= 0 && breakEvenMonths == -1 {
			breakEvenMonths = i
		}
		if i == m.savingsAfterMonths {
			savingsAfter = v
		}
	}
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		PaddingLeft(1).
		PaddingRight(1).
		Width(40)
	heading := lipgloss.JoinHorizontal(lipgloss.Bottom,
		style.Render("Months to break even:"),
		style.Render(fmt.Sprintf("Savings from buying after %d years:",
			m.savingsAfterMonths/12,
		)))
	dataColor := lipgloss.Color("#04B575")
	if savingsAfter < 0 {
		dataColor = lipgloss.Color("#8B0000")
	}
	style = style.Foreground(dataColor)
	data := lipgloss.JoinHorizontal(lipgloss.Bottom,
		style.Render(fmt.Sprintf("%d (%d years)",
			breakEvenMonths,
			breakEvenMonths/12),
		),
		style.Render(fmt.Sprintf("$%.0f",
			savingsAfter,
		)))
	b.WriteString(lipgloss.NewStyle().PaddingTop(1).PaddingBottom(1).Render(
		lipgloss.JoinVertical(lipgloss.Center, heading, data)))
	b.WriteString("\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}
	button := blurredButton
	if m.focusIndex == len(m.inputs) {
		button = focusedButton
	}
	fmt.Fprintf(&b, "\n%s\n\n", button)

	return b.String()
}
