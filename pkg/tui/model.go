package tui

import (
	"net"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dhanushnehru/domain-doctor/internal/health"
	"github.com/dhanushnehru/domain-doctor/internal/resolver"
)

type state int

const (
	stateInput state = iota
	stateChecking
	stateDone
)

type (
	checkDoneMsg struct{ result *health.DomainHealth }
	checkErrMsg  struct{ err error }
)

type Model struct {
	state         state
	focused       int
	errMsg        string
	domainInput   textinput.Model
	resolverInput textinput.Model
	spinner       spinner.Model
	result        *health.DomainHealth
	err           error
}

func New() Model {
	dI := textinput.New()
	dI.Placeholder = "e.g. google.com"
	dI.Focus()
	dI.CharLimit = 253

	rI := textinput.New()
	rI.Placeholder = "e.g. 1.1.1.1:53"
	rI.CharLimit = 15

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))

	return Model{state: stateInput, domainInput: dI, resolverInput: rI, spinner: sp}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyShiftTab:
			m.focused = (m.focused + 1) % 2

			if m.focused == 0 {
				m.domainInput.Focus()
				m.resolverInput.Blur()
			} else {
				m.resolverInput.Focus()
				m.domainInput.Blur()
			}
			return m, nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.state == stateInput {
				domain := strings.TrimSpace(m.domainInput.Value())
				resolver := strings.TrimSpace(m.resolverInput.Value())

				if !validDomain(domain) {
					m.errMsg = "Invalid domain format"
					return m, nil
				}

				if !validResolver(resolver) {
					m.errMsg = "Resolver must be a valid IP address"
					return m, nil
				}

				m.errMsg = ""
				m.state = stateChecking
				return m, tea.Batch(m.spinner.Tick, runCheck(domain, resolver))
			}
			if m.state == stateDone {
				return m, tea.Quit
			}
		}

	case checkDoneMsg:
		m.state = stateDone
		m.result = msg.result
		return m, nil

	case checkErrMsg:
		m.err = msg.err
		m.state = stateDone
		return m, nil

	case spinner.TickMsg:
		if m.state == stateChecking {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if m.state == stateInput {
		var cmd1, cmd2 tea.Cmd
		m.domainInput, cmd1 = m.domainInput.Update(msg)
		m.resolverInput, cmd2 = m.resolverInput.Update(msg)
		return m, tea.Batch(cmd1, cmd2)
	}

	return m, nil
}

func (m Model) View() string {
	switch m.state {
	case stateInput:
		return inputView(m)
	case stateChecking:
		return checkingView(m)
	case stateDone:
		return doneView(m)
	}
	return ""
}

func inputView(m Model) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		promptTitleStyle.Render("🩺 Domain Doctor"),

		promptLabelStyle.Render("Domain"),
		m.domainInput.View(),

		promptLabelStyle.Render("Resolver (optional)"),
		m.resolverInput.View(),
	)

	if m.errMsg != "" {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			errorStyle.Render(m.errMsg),
		)
	}

	return "\n" + promptCardStyle.Render(content) + "\n"
}

func checkingView(m Model) string {
	row := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.spinner.View(),
		spinnerLabelStyle.Render(" Checking "+m.domainInput.Value()+"..."),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		promptTitleStyle.Render("🩺 Domain Doctor"),
		row,
	)

	return "\n" + promptCardStyle.Render(content) + "\n"
}

func doneView(m Model) string {
	var body string

	if m.err != nil {
		body = errorStyle.Render(m.err.Error())
	} else {
		body = m.result.RenderReport(true)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		promptTitleStyle.Render("🩺 Domain Doctor"),
		body,
		promptLabelStyle.Render("Press Enter to exit"),
	)

	return "\n" + promptCardStyle.Render(content) + "\n"
}

func runCheck(domain string, resolverName string) tea.Cmd {
	return func() tea.Msg {

		r := resolver.New(resolverName)
		result := health.Check(domain, *r)

		return checkDoneMsg{result: result}
	}
}

func validDomain(d string) bool {
	if len(d) == 0 || len(d) > 253 {
		return false
	}
	parts := strings.Split(d, ".")
	if len(parts) < 2 {
		return false
	}
	for _, p := range parts {
		if len(p) == 0 || len(p) > 63 {
			return false
		}
	}
	return true
}

func validResolver(r string) bool {
	if r == "" {
		return true
	}

	host, port, err := net.SplitHostPort(r)
	if err != nil {
		return false
	}

	// Validate IP
	if net.ParseIP(host) == nil {
		return false
	}

	// Validate port range
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return false
	}

	return true
}
