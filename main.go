package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF69B4")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB6C1"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#98FB98"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDA0DD"))
)

type state int

const (
	mainMenu state = iota
	installedView
	searchView
	installing
	uninstalling
)

type brewPackage struct {
	name        string
	description string
	installed   bool
	pkgType     string
}

func (p brewPackage) Title() string       { return p.name }
func (p brewPackage) Description() string { return p.description }
func (p brewPackage) FilterValue() string { return p.name }

type model struct {
	state         state
	list          list.Model
	textInput     textinput.Model
	spinner       spinner.Model
	loading       bool
	statusMessage string
	err           error
	packages      []brewPackage
	searchResults []brewPackage
	installedPkgs []brewPackage
}

type packagesLoadedMsg struct {
	packages []brewPackage
}

type searchResultsMsg struct {
	results []brewPackage
}

type installCompleteMsg struct {
	success bool
	pkg     string
	err     error
}

type uninstallCompleteMsg struct {
	success bool
	pkg     string
	err     error
}

func initialModel() model {
	items := []list.Item{}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Brew-Tea Package Manager"
	l.SetShowStatusBar(false)

	ti := textinput.New()
	ti.Placeholder = "Search for packages..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Moon

	return model{
		state:     mainMenu,
		list:      l,
		textInput: ti,
		spinner:   s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		loadInstalledPackages,
		m.spinner.Tick,
	)
}

func loadInstalledPackages() tea.Msg {
	var packages []brewPackage

	cmd := exec.Command("brew", "list", "--formula")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				packages = append(packages, brewPackage{
					name:        strings.TrimSpace(line),
					description: "Installed formula",
					installed:   true,
					pkgType:     "formula",
				})
			}
		}
	}

	cmd = exec.Command("brew", "list", "--cask")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				packages = append(packages, brewPackage{
					name:        strings.TrimSpace(line),
					description: "Installed cask",
					installed:   true,
					pkgType:     "cask",
				})
			}
		}
	}

	return packagesLoadedMsg{packages: packages}
}

func searchPackages(query string) tea.Cmd {
	return func() tea.Msg {
		if query == "" {
			return searchResultsMsg{results: []brewPackage{}}
		}

		var results []brewPackage

		cmd := exec.Command("brew", "search", "--formula", query)
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.Contains(line, "==>") {
					results = append(results, brewPackage{
						name:        line,
						description: "Available formula",
						installed:   false,
						pkgType:     "formula",
					})
				}
			}
		}

		cmd = exec.Command("brew", "search", "--cask", query)
		output, err = cmd.Output()
		if err == nil {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.Contains(line, "==>") {
					results = append(results, brewPackage{
						name:        line,
						description: "Available cask",
						installed:   false,
						pkgType:     "cask",
					})
				}
			}
		}

		return searchResultsMsg{results: results}
	}
}

func installPackage(pkg brewPackage) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		if pkg.pkgType == "cask" {
			cmd = exec.Command("brew", "install", "--cask", pkg.name)
		} else {
			cmd = exec.Command("brew", "install", pkg.name)
		}
		err := cmd.Run()
		return installCompleteMsg{
			success: err == nil,
			pkg:     pkg.name,
			err:     err,
		}
	}
}

func uninstallPackage(pkg brewPackage) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		if pkg.pkgType == "cask" {
			cmd = exec.Command("brew", "uninstall", "--cask", pkg.name)
		} else {
			cmd = exec.Command("brew", "uninstall", pkg.name)
		}
		err := cmd.Run()
		return uninstallCompleteMsg{
			success: err == nil,
			pkg:     pkg.name,
			err:     err,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case mainMenu:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "1":
				m.state = installedView
				items := make([]list.Item, len(m.installedPkgs))
				for i, pkg := range m.installedPkgs {
					items[i] = pkg
				}
				m.list.SetItems(items)
				m.list.Title = "Installed Packages (press 'u' to uninstall, 'b' to go back)"
			case "2":
				m.state = searchView
				m.textInput.Focus()
				m.list.Title = "Search Results (press 'ctrl+i' to install, 'b' to go back)"
			}

		case installedView:
			switch msg.String() {
			case "b":
				m.state = mainMenu
				m.list.Title = "Brew-Tea Package Manager"
				m.list.SetItems([]list.Item{})
			case "u":
				if len(m.list.Items()) > 0 {
					selected := m.list.SelectedItem().(brewPackage)
					m.state = uninstalling
					m.loading = true
					m.statusMessage = fmt.Sprintf("Uninstalling %s...", selected.name)
					return m, tea.Batch(uninstallPackage(selected), m.spinner.Tick)
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				cmds = append(cmds, cmd)
			}

		case searchView:
			switch msg.String() {
			case "b":
				m.state = mainMenu
				m.list.Title = "Brew-Tea Package Manager"
				m.list.SetItems([]list.Item{})
				m.textInput.SetValue("")
			case "ctrl+i":
				if len(m.list.Items()) > 0 {
					selected := m.list.SelectedItem().(brewPackage)
					m.state = installing
					m.loading = true
					m.statusMessage = fmt.Sprintf("Installing %s...", selected.name)
					return m, tea.Batch(installPackage(selected), m.spinner.Tick)
				}
			case "enter":
				query := m.textInput.Value()
				if query != "" {
					m.loading = true
					m.statusMessage = "Searching..."
					return m, tea.Batch(searchPackages(query), m.spinner.Tick)
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				cmds = append(cmds, cmd)

				var listCmd tea.Cmd
				m.list, listCmd = m.list.Update(msg)
				cmds = append(cmds, listCmd)
			}
		}

	case packagesLoadedMsg:
		m.installedPkgs = msg.packages

	case searchResultsMsg:
		m.loading = false
		m.searchResults = msg.results
		items := make([]list.Item, len(msg.results))
		for i, pkg := range msg.results {
			items[i] = pkg
		}
		m.list.SetItems(items)

	case installCompleteMsg:
		m.loading = false
		m.state = searchView
		if msg.success {
			m.statusMessage = successStyle.Render(fmt.Sprintf("✓ Successfully installed %s", msg.pkg))
			cmds = append(cmds, loadInstalledPackages)
		} else {
			m.statusMessage = errorStyle.Render(fmt.Sprintf("✗ Failed to install %s: %v", msg.pkg, msg.err))
		}

	case uninstallCompleteMsg:
		m.loading = false
		m.state = installedView
		if msg.success {
			m.statusMessage = successStyle.Render(fmt.Sprintf("✓ Successfully uninstalled %s", msg.pkg))
			cmds = append(cmds, loadInstalledPackages)
		} else {
			m.statusMessage = errorStyle.Render(fmt.Sprintf("✗ Failed to uninstall %s: %v", msg.pkg, msg.err))
		}

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s strings.Builder

	// Title
	s.WriteString(`   ( (
    ) )
  ........
  |      |]
  \      /
   '----'
`)
	s.WriteString(titleStyle.Render("� Brew-Tea Package Manager"))
	s.WriteString("\n\n")

	if m.statusMessage != "" {
		s.WriteString(m.statusMessage)
		s.WriteString("\n\n")
	}

	switch m.state {
	case mainMenu:
		s.WriteString("What would you like to do?\n\n")
		s.WriteString("1. View installed packages\n")
		s.WriteString("2. Search and install packages\n\n")
		s.WriteString(helpStyle.Render("Press 'q' to quit"))

	case installedView, searchView:
		if m.state == searchView {
			s.WriteString("Search: ")
			s.WriteString(m.textInput.View())
			s.WriteString("\n")
			if !m.loading {
				s.WriteString(helpStyle.Render("Press enter to search"))
				s.WriteString("\n")
			}
			s.WriteString("\n")
		}

		if m.loading {
			s.WriteString(m.spinner.View() + " " + m.statusMessage)
		} else {
			s.WriteString(m.list.View())
		}

	case installing, uninstalling:
		s.WriteString(m.spinner.View() + " " + m.statusMessage)
	}

	return s.String()
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help", "-h":
			fmt.Println("brew-tea: A Bubble Tea package manager for Homebrew")
			fmt.Println("\nUsage:")
			fmt.Println("  brew-tea        Run the interactive TUI")
			fmt.Println("  brew-tea --help Show this help message")
			os.Exit(0)
		}
	}

	if _, err := exec.LookPath("brew"); err != nil {
		fmt.Println("Error: Homebrew is not installed or not in PATH")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
