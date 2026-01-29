package tui

import (
	"fmt"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/minicodemonkey/chief/internal/git"
)

// FirstTimeSetupResult contains the result of the first-time setup flow.
type FirstTimeSetupResult struct {
	PRDName        string
	AddedGitignore bool
	Cancelled      bool
}

// FirstTimeSetupStep represents the current step in the setup flow.
type FirstTimeSetupStep int

const (
	StepGitignore FirstTimeSetupStep = iota
	StepPRDName
)

// FirstTimeSetup is a TUI for first-time project setup.
type FirstTimeSetup struct {
	width  int
	height int

	step          FirstTimeSetupStep
	showGitignore bool // Whether to show the gitignore step

	// Gitignore step
	gitignoreSelected int // 0 = Yes, 1 = No

	// PRD name step
	prdName      string
	prdNameError string

	// Result
	result FirstTimeSetupResult

	baseDir string
}

// NewFirstTimeSetup creates a new first-time setup TUI.
func NewFirstTimeSetup(baseDir string, showGitignore bool) *FirstTimeSetup {
	step := StepPRDName
	if showGitignore {
		step = StepGitignore
	}
	return &FirstTimeSetup{
		baseDir:           baseDir,
		showGitignore:     showGitignore,
		step:              step,
		gitignoreSelected: 0, // Default to "Yes"
		prdName:           "main",
	}
}

// Init initializes the model.
func (f FirstTimeSetup) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update handles messages.
func (f FirstTimeSetup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
		f.height = msg.Height
		return f, nil

	case tea.KeyMsg:
		switch f.step {
		case StepGitignore:
			return f.handleGitignoreKeys(msg)
		case StepPRDName:
			return f.handlePRDNameKeys(msg)
		}
	}
	return f, nil
}

func (f FirstTimeSetup) handleGitignoreKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c", "esc":
		f.result.Cancelled = true
		return f, tea.Quit

	case "up", "k", "left", "h":
		if f.gitignoreSelected > 0 {
			f.gitignoreSelected--
		}
		return f, nil

	case "down", "j", "right", "l":
		if f.gitignoreSelected < 1 {
			f.gitignoreSelected++
		}
		return f, nil

	case "y", "Y":
		f.gitignoreSelected = 0
		return f.confirmGitignore()

	case "n", "N":
		f.gitignoreSelected = 1
		return f.confirmGitignore()

	case "enter":
		return f.confirmGitignore()
	}
	return f, nil
}

func (f FirstTimeSetup) confirmGitignore() (tea.Model, tea.Cmd) {
	if f.gitignoreSelected == 0 {
		// User wants to add .chief to gitignore
		if err := git.AddChiefToGitignore(f.baseDir); err != nil {
			// Show error but continue
			f.prdNameError = "Warning: failed to add .chief to .gitignore"
		} else {
			f.result.AddedGitignore = true
		}
	}
	f.step = StepPRDName
	return f, nil
}

func (f FirstTimeSetup) handlePRDNameKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		f.result.Cancelled = true
		return f, tea.Quit

	case "esc":
		if f.showGitignore {
			// Go back to gitignore step
			f.step = StepGitignore
			f.prdNameError = ""
			return f, nil
		}
		f.result.Cancelled = true
		return f, tea.Quit

	case "enter":
		// Validate PRD name
		name := strings.TrimSpace(f.prdName)
		if name == "" {
			f.prdNameError = "Name cannot be empty"
			return f, nil
		}
		if !isValidPRDName(name) {
			f.prdNameError = "Name can only contain letters, numbers, hyphens, and underscores"
			return f, nil
		}
		f.result.PRDName = name
		return f, tea.Quit

	case "backspace":
		if len(f.prdName) > 0 {
			f.prdName = f.prdName[:len(f.prdName)-1]
			f.prdNameError = ""
		}
		return f, nil

	default:
		// Handle character input
		if len(msg.String()) == 1 {
			r := rune(msg.String()[0])
			// Only allow valid characters
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '-' || r == '_' {
				f.prdName += string(r)
				f.prdNameError = ""
			}
		}
		return f, nil
	}
}

// isValidPRDName checks if a name is valid for a PRD.
func isValidPRDName(name string) bool {
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return validName.MatchString(name)
}

// View renders the TUI.
func (f FirstTimeSetup) View() string {
	switch f.step {
	case StepGitignore:
		return f.renderGitignoreStep()
	case StepPRDName:
		return f.renderPRDNameStep()
	default:
		return ""
	}
}

func (f FirstTimeSetup) renderGitignoreStep() string {
	modalWidth := min(65, f.width-10)
	if modalWidth < 45 {
		modalWidth = 45
	}

	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor)
	content.WriteString(titleStyle.Render("Welcome to Chief!"))
	content.WriteString("\n")
	content.WriteString(DividerStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n\n")

	// Message
	messageStyle := lipgloss.NewStyle().Foreground(TextColor)
	content.WriteString(messageStyle.Render("Would you like to add .chief to .gitignore?"))
	content.WriteString("\n\n")

	descStyle := lipgloss.NewStyle().Foreground(MutedColor)
	content.WriteString(descStyle.Render("This keeps your PRD plans local and out of version control."))
	content.WriteString("\n")
	content.WriteString(descStyle.Render("Not required, but recommended if you prefer local-only plans."))
	content.WriteString("\n\n")

	// Options
	optionStyle := lipgloss.NewStyle().Foreground(TextColor)
	selectedStyle := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)

	options := []struct {
		label string
		desc  string
	}{
		{"Yes, add .chief to .gitignore", "(Recommended)"},
		{"No, keep .chief in version control", ""},
	}

	for i, opt := range options {
		var line string
		if i == f.gitignoreSelected {
			line = selectedStyle.Render(fmt.Sprintf("▶ %s", opt.label))
			if opt.desc != "" {
				line += " " + lipgloss.NewStyle().Foreground(SuccessColor).Render(opt.desc)
			}
		} else {
			line = optionStyle.Render(fmt.Sprintf("  %s", opt.label))
			if opt.desc != "" {
				line += " " + lipgloss.NewStyle().Foreground(MutedColor).Render(opt.desc)
			}
		}
		content.WriteString(line)
		content.WriteString("\n")
	}

	// Footer
	content.WriteString("\n")
	content.WriteString(DividerStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n")

	footerStyle := lipgloss.NewStyle().Foreground(MutedColor)
	content.WriteString(footerStyle.Render("↑/↓: Navigate  Enter: Select  y/n: Quick select  Esc: Cancel"))

	// Modal box
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2).
		Width(modalWidth)

	modal := modalStyle.Render(content.String())

	return f.centerModal(modal)
}

func (f FirstTimeSetup) renderPRDNameStep() string {
	modalWidth := min(60, f.width-10)
	if modalWidth < 45 {
		modalWidth = 45
	}

	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor)

	if f.showGitignore && f.result.AddedGitignore {
		content.WriteString(lipgloss.NewStyle().Foreground(SuccessColor).Render("✓ Added .chief to .gitignore"))
		content.WriteString("\n\n")
	}

	content.WriteString(titleStyle.Render("Create Your First PRD"))
	content.WriteString("\n")
	content.WriteString(DividerStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n\n")

	// Message
	messageStyle := lipgloss.NewStyle().Foreground(TextColor)
	content.WriteString(messageStyle.Render("Enter a name for your PRD:"))
	content.WriteString("\n\n")

	// Input field
	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(0, 1).
		Width(modalWidth - 8)

	displayName := f.prdName
	if displayName == "" {
		displayName = " " // Show cursor position
	}
	content.WriteString(inputStyle.Render(displayName + "█"))
	content.WriteString("\n")

	// Error message
	if f.prdNameError != "" {
		errorStyle := lipgloss.NewStyle().Foreground(ErrorColor)
		content.WriteString("\n")
		content.WriteString(errorStyle.Render(f.prdNameError))
	}

	// Hint
	content.WriteString("\n")
	hintStyle := lipgloss.NewStyle().Foreground(MutedColor)
	content.WriteString(hintStyle.Render("PRD will be created at: .chief/prds/" + f.prdName + "/"))

	// Footer
	content.WriteString("\n\n")
	content.WriteString(DividerStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n")

	footerStyle := lipgloss.NewStyle().Foreground(MutedColor)
	if f.showGitignore {
		content.WriteString(footerStyle.Render("Enter: Create PRD  Esc: Back  Ctrl+C: Cancel"))
	} else {
		content.WriteString(footerStyle.Render("Enter: Create PRD  Esc/Ctrl+C: Cancel"))
	}

	// Modal box
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2).
		Width(modalWidth)

	modal := modalStyle.Render(content.String())

	return f.centerModal(modal)
}

func (f FirstTimeSetup) centerModal(modal string) string {
	lines := strings.Split(modal, "\n")
	modalHeight := len(lines)
	modalWidth := 0
	for _, line := range lines {
		if lipgloss.Width(line) > modalWidth {
			modalWidth = lipgloss.Width(line)
		}
	}

	topPadding := (f.height - modalHeight) / 2
	leftPadding := (f.width - modalWidth) / 2

	if topPadding < 0 {
		topPadding = 0
	}
	if leftPadding < 0 {
		leftPadding = 0
	}

	var result strings.Builder

	for i := 0; i < topPadding; i++ {
		result.WriteString("\n")
	}

	leftPad := strings.Repeat(" ", leftPadding)
	for _, line := range lines {
		result.WriteString(leftPad)
		result.WriteString(line)
		result.WriteString("\n")
	}

	return result.String()
}

// GetResult returns the setup result.
func (f FirstTimeSetup) GetResult() FirstTimeSetupResult {
	return f.result
}

// RunFirstTimeSetup runs the first-time setup TUI and returns the result.
func RunFirstTimeSetup(baseDir string, showGitignore bool) (FirstTimeSetupResult, error) {
	setup := NewFirstTimeSetup(baseDir, showGitignore)
	p := tea.NewProgram(setup, tea.WithAltScreen())

	model, err := p.Run()
	if err != nil {
		return FirstTimeSetupResult{Cancelled: true}, err
	}

	if finalSetup, ok := model.(FirstTimeSetup); ok {
		return finalSetup.GetResult(), nil
	}

	return FirstTimeSetupResult{Cancelled: true}, nil
}
