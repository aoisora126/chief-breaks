package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// BranchWarningOption represents an option in the branch warning dialog.
type BranchWarningOption int

const (
	BranchOptionCreateBranch BranchWarningOption = iota
	BranchOptionContinue
	BranchOptionCancel
)

// BranchWarning manages the branch warning dialog state.
type BranchWarning struct {
	width         int
	height        int
	currentBranch string
	prdName       string
	selectedIndex int
	editMode      bool   // Whether we're editing the branch name
	branchName    string // The current branch name (editable)
}

// NewBranchWarning creates a new branch warning dialog.
func NewBranchWarning() *BranchWarning {
	return &BranchWarning{
		selectedIndex: 0, // Default to "Create branch" option
	}
}

// SetSize sets the dialog dimensions.
func (b *BranchWarning) SetSize(width, height int) {
	b.width = width
	b.height = height
}

// SetContext sets the branch and PRD context for the warning.
func (b *BranchWarning) SetContext(currentBranch, prdName string) {
	b.currentBranch = currentBranch
	b.prdName = prdName
	b.branchName = fmt.Sprintf("chief/%s", prdName)
}

// GetSuggestedBranch returns the branch name (may be edited by user).
func (b *BranchWarning) GetSuggestedBranch() string {
	return b.branchName
}

// MoveUp moves selection up.
func (b *BranchWarning) MoveUp() {
	if b.selectedIndex > 0 {
		b.selectedIndex--
	}
}

// MoveDown moves selection down.
func (b *BranchWarning) MoveDown() {
	if b.selectedIndex < 2 {
		b.selectedIndex++
	}
}

// GetSelectedOption returns the currently selected option.
func (b *BranchWarning) GetSelectedOption() BranchWarningOption {
	return BranchWarningOption(b.selectedIndex)
}

// Reset resets the dialog state.
func (b *BranchWarning) Reset() {
	b.selectedIndex = 0
	b.editMode = false
	b.branchName = fmt.Sprintf("chief/%s", b.prdName)
}

// IsEditMode returns true if the branch name is being edited.
func (b *BranchWarning) IsEditMode() bool {
	return b.editMode
}

// StartEditMode enters edit mode for the branch name.
func (b *BranchWarning) StartEditMode() {
	b.editMode = true
}

// CancelEditMode exits edit mode.
func (b *BranchWarning) CancelEditMode() {
	b.editMode = false
}

// AddInputChar adds a character to the branch name.
func (b *BranchWarning) AddInputChar(ch rune) {
	// Only allow valid git branch name characters
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') || ch == '-' || ch == '_' || ch == '/' {
		b.branchName += string(ch)
	}
}

// DeleteInputChar removes the last character from the branch name.
func (b *BranchWarning) DeleteInputChar() {
	if len(b.branchName) > 0 {
		b.branchName = b.branchName[:len(b.branchName)-1]
	}
}

// Render renders the branch warning dialog.
func (b *BranchWarning) Render() string {
	// Modal dimensions
	modalWidth := min(60, b.width-10)
	modalHeight := min(16, b.height-6)

	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalHeight < 12 {
		modalHeight = 12
	}

	// Build modal content
	var content strings.Builder

	// Warning icon and title
	warningStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(WarningColor)
	content.WriteString(warningStyle.Render("⚠️  Protected Branch Warning"))
	content.WriteString("\n")
	content.WriteString(DividerStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n\n")

	// Warning message
	messageStyle := lipgloss.NewStyle().Foreground(TextColor)
	content.WriteString(messageStyle.Render(fmt.Sprintf("You are on the '%s' branch.", b.currentBranch)))
	content.WriteString("\n")
	content.WriteString(messageStyle.Render("Starting the loop will make changes directly to this branch."))
	content.WriteString("\n\n")

	// Options
	optionStyle := lipgloss.NewStyle().Foreground(TextColor)
	selectedOptionStyle := lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)

	// Render the "Create branch" option with editable field
	if b.selectedIndex == 0 {
		content.WriteString(selectedOptionStyle.Render("▶ Create branch "))
		if b.editMode {
			// Show editable input field
			inputStyle := lipgloss.NewStyle().
				Foreground(TextBrightColor).
				Background(lipgloss.Color("237"))
			cursorStyle := lipgloss.NewStyle().Foreground(PrimaryColor).Blink(true)
			content.WriteString(inputStyle.Render(b.branchName))
			content.WriteString(cursorStyle.Render("▌"))
		} else {
			// Show branch name with edit hint
			content.WriteString(selectedOptionStyle.Render(fmt.Sprintf("'%s'", b.branchName)))
			content.WriteString(" ")
			content.WriteString(lipgloss.NewStyle().Foreground(SuccessColor).Render("(Recommended)"))
			content.WriteString(" ")
			content.WriteString(lipgloss.NewStyle().Foreground(MutedColor).Render("[e: edit]"))
		}
	} else {
		content.WriteString(optionStyle.Render(fmt.Sprintf("  Create branch '%s'", b.branchName)))
		content.WriteString(" ")
		content.WriteString(lipgloss.NewStyle().Foreground(MutedColor).Render("(Recommended)"))
	}
	content.WriteString("\n")

	// Render "Continue on current branch" option
	if b.selectedIndex == 1 {
		content.WriteString(selectedOptionStyle.Render(fmt.Sprintf("▶ Continue on '%s'", b.currentBranch)))
	} else {
		content.WriteString(optionStyle.Render(fmt.Sprintf("  Continue on '%s'", b.currentBranch)))
	}
	content.WriteString("\n")

	// Render "Cancel" option
	if b.selectedIndex == 2 {
		content.WriteString(selectedOptionStyle.Render("▶ Cancel"))
	} else {
		content.WriteString(optionStyle.Render("  Cancel"))
	}
	content.WriteString("\n")

	// Footer
	content.WriteString("\n")
	content.WriteString(DividerStyle.Render(strings.Repeat("─", modalWidth-4)))
	content.WriteString("\n")

	footerStyle := lipgloss.NewStyle().
		Foreground(MutedColor)
	if b.editMode {
		content.WriteString(footerStyle.Render("Enter: confirm  Esc: cancel edit"))
	} else {
		content.WriteString(footerStyle.Render("↑/↓: Navigate  Enter: Select  Esc: Cancel"))
	}

	// Modal box style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(WarningColor).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight)

	modal := modalStyle.Render(content.String())

	// Center the modal on screen
	return b.centerModal(modal)
}

// centerModal centers the modal on the screen.
func (b *BranchWarning) centerModal(modal string) string {
	lines := strings.Split(modal, "\n")
	modalHeight := len(lines)
	modalWidth := 0
	for _, line := range lines {
		if lipgloss.Width(line) > modalWidth {
			modalWidth = lipgloss.Width(line)
		}
	}

	// Calculate padding
	topPadding := (b.height - modalHeight) / 2
	leftPadding := (b.width - modalWidth) / 2

	if topPadding < 0 {
		topPadding = 0
	}
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Build centered content
	var result strings.Builder

	// Top padding
	for i := 0; i < topPadding; i++ {
		result.WriteString("\n")
	}

	// Modal lines with left padding
	leftPad := strings.Repeat(" ", leftPadding)
	for _, line := range lines {
		result.WriteString(leftPad)
		result.WriteString(line)
		result.WriteString("\n")
	}

	return result.String()
}
