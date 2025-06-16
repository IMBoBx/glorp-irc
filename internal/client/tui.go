package client

import (
	"fmt"
	"net"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textarea    textarea.Model
	viewport    viewport.Model
	messages    []string
	senderStyle lipgloss.Style
	err         error
	conn        net.Conn
}

const gap = "\n\n"

func InitialModel(conn net.Conn) model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 280
	ta.SetWidth(30)
	ta.SetHeight(3)
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(30, 5)
	vp.SetContent("               _                        _        \n" +
" __ __ __ ___ | | __  ___  _ __   ___  | |_  ___ \n" +
" \\ V  V // -_)| |/ _|/ _ \\| '  \\ / -_) |  _|/ _ \\\n" +
"  \\_/\\_/ \\___||_|\\__|\\___/|_|_|_|\\___|  \\__|\\___/\n" +
"          __ _  | |  ___   _ _   _ __            \n" +
"         / _` | | | / _ \\ | '_| | '_ \\           \n" +
"         \\__, | |_| \\___/ |_|   | .__/           \n" +
"         |___/                  |_|              \n" + 
"Type `/join <username> <room-name>` for joining a channel.")

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Bold(true),
		err:          nil,
		conn:        conn,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

type IncomingMsg string

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case IncomingMsg:
		m.messages = append(m.messages, string(msg))
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "")))
		m.viewport.GotoBottom()

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "")))
		}
		m.viewport.GotoBottom()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter:
			if val := m.textarea.Value(); len(val) > 0 {
				// m.messages = append(m.messages, m.senderStyle.Render("You:")+" "+val)
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
				m.conn.Write([]byte(val + "\n"))
				m.textarea.Reset()
				m.viewport.GotoBottom()
			}
		}

	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
