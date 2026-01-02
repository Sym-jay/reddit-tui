package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type post struct {
	title     string
	subreddit string
	author    string
	upvotes   int
	comments  int
}

type model struct {
	sidebarItems  []string
	posts         []post
	sidebarCursor int
	postsCursor   int
	activePane    string
	width         int
	height        int
}

func initialModel() model {
	return model{
		sidebarItems: []string{"Home", "Popular", "Explore", "Settings", "Login/Auth"},
		posts: []post{
			{title: "Building a Reddit TUI with Go and Bubble Tea", subreddit: "r/golang", author: "gopher_dev", upvotes: 342, comments: 45},
			{title: "What are your favorite terminal tools?", subreddit: "r/commandline", author: "cli_enthusiast", upvotes: 528, comments: 89},
			{title: "Show HN: My weekend project - a Reddit client for the terminal", subreddit: "r/programming", author: "weekend_coder", upvotes: 1205, comments: 134},
			{title: "TUI vs GUI: The eternal debate", subreddit: "r/linux", author: "terminal_lover", upvotes: 876, comments: 201},
			{title: "Charm libraries are amazing for building TUIs", subreddit: "r/golang", author: "bubble_fan", upvotes: 445, comments: 67},
			{title: "Ask Reddit: What's your development setup?", subreddit: "r/AskReddit", author: "curious_dev", upvotes: 2301, comments: 456},
			{title: "Vim vs Emacs: A comprehensive comparison", subreddit: "r/programming", author: "editor_wars", upvotes: 689, comments: 342},
			{title: "Why I switched from GUI apps to terminal", subreddit: "r/commandline", author: "minimalist_dev", upvotes: 934, comments: 178},
		},
		sidebarCursor: 0,
		postsCursor:   0,
		activePane:    "sidebar",
		width:         80,
		height:        24,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.activePane == "sidebar" {
				m.activePane = "posts"
			} else {
				m.activePane = "sidebar"
			}
		case "up", "k":
			if m.activePane == "sidebar" {
				if m.sidebarCursor > 0 {
					m.sidebarCursor--
				}
			} else {
				if m.postsCursor > 0 {
					m.postsCursor--
				}
			}
		case "down", "j":
			if m.activePane == "sidebar" {
				if m.sidebarCursor < len(m.sidebarItems)-1 {
					m.sidebarCursor++
				}
			} else {
				if m.postsCursor < len(m.posts)-1 {
					m.postsCursor++
				}
			}
		}
	}

	return m, nil
}
func renderPane(content string, width, height int, borderColor string, active bool) string {
	innerWidth := width - 2
	innerHeight := height - 2

	if innerWidth < 1 {
		innerWidth = 1
	}
	if innerHeight < 1 {
		innerHeight = 1
	}

	lines := strings.Split(content, "\n")

	result := make([]string, innerHeight)
	for i := 0; i < innerHeight; i++ {
		if i < len(lines) {
			line := lines[i]

			if lipgloss.Width(line) > innerWidth {
				// Simple truncation - cut characters
				runes := []rune(line)
				if len(runes) > innerWidth {
					line = string(runes[:innerWidth])
				}
			}

			result[i] = line + strings.Repeat(" ", innerWidth-lipgloss.Width(line))
		} else {

			result[i] = strings.Repeat(" ", innerWidth)
		}
	}

	innerContent := strings.Join(result, "\n")

	color := lipgloss.Color(borderColor)
	if active {
		color = lipgloss.Color("205")
	}

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Width(innerWidth).
		Height(innerHeight)

	return style.Render(innerContent)
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	controlPaneHeight := 3

	sidebarWidth := m.width / 5
	if sidebarWidth < 15 {
		sidebarWidth = 15
	}
	remainingWidth := m.width - sidebarWidth
	postsWidth := remainingWidth / 2
	previewWidth := remainingWidth - postsWidth

	paneHeight := m.height - controlPaneHeight

	postsPaneHeading := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginLeft(2)
	navPaneHeading := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginLeft(2)
	previewPaneHeading := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginLeft(2)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	postTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	subredditStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	sidebarContent := navPaneHeading.Render("NAVIGATION") + "\n\n"
	for i, item := range m.sidebarItems {
		cursor := "  "
		if m.sidebarCursor == i {
			cursor = cursorStyle.Render("> ")
		}
		sidebarContent += cursor + item + "\n"
	}

	postsContent := postsPaneHeading.Render("POSTS") + "\n\n"
	for i, post := range m.posts {
		cursor := "  "
		if m.postsCursor == i {
			cursor = cursorStyle.Render("> ")
		}
		postsContent += cursor + postTitleStyle.Render(post.title) + "\n"
		postsContent += "   " + subredditStyle.Render(post.subreddit) + " by u/" + post.author + "\n"
		postsContent += "   " + metaStyle.Render(fmt.Sprintf("%d upvotes | %d comments", post.upvotes, post.comments)) + "\n\n"
	}

	var previewContent string
	if m.postsCursor >= 0 && m.postsCursor < len(m.posts) {
		selectedPost := m.posts[m.postsCursor]
		previewContent = previewPaneHeading.Render("PREVIEW") + "\n\n"
		previewContent += postTitleStyle.Render(selectedPost.title) + "\n\n"
		previewContent += subredditStyle.Render(selectedPost.subreddit) + " by u/" + selectedPost.author + "\n"
		previewContent += metaStyle.Render(fmt.Sprintf("%d upvotes | %d comments", selectedPost.upvotes, selectedPost.comments)) + "\n\n"
		previewContent += strings.Repeat("-", 20) + "\n\n"
		previewContent += "Lorem ipsum dolor sit amet,\n"
		previewContent += "consectetur adipiscing elit.\n\n"
		previewContent += "Sed do eiusmod tempor incididunt\n"
		previewContent += "ut labore et dolore magna aliqua."
	} else {
		previewContent = "PREVIEW\n\nSelect a post to view"
	}

	sidebar := renderPane(sidebarContent, sidebarWidth, paneHeight, "63", m.activePane == "sidebar")
	posts := renderPane(postsContent, postsWidth, paneHeight, "63", m.activePane == "posts")
	preview := renderPane(previewContent, previewWidth, paneHeight, "63", false)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, posts, preview)

	controlText := metaStyle.Render("Tab: switch panes | ↑↓/j/k: navigate | q: quit")
	controlPane := renderPane(controlText, m.width, controlPaneHeight, "63", false)

	return lipgloss.JoinVertical(lipgloss.Left, mainContent, controlPane)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error occurred: %v", err)
		os.Exit(1)
	}
}
