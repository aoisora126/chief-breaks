package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/minicodemonkey/chief/embed"
)

// InitOptions contains configuration for the init command.
type InitOptions struct {
	Name    string // PRD name (default: "main")
	Context string // Optional context to pass to Claude
	BaseDir string // Base directory for .chief/prds/ (default: current directory)
}

// RunInit creates a new PRD by launching an interactive Claude session.
func RunInit(opts InitOptions) error {
	// Set defaults
	if opts.Name == "" {
		opts.Name = "main"
	}
	if opts.BaseDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		opts.BaseDir = cwd
	}

	// Validate name (alphanumeric, -, _)
	if !isValidPRDName(opts.Name) {
		return fmt.Errorf("invalid PRD name %q: must contain only letters, numbers, hyphens, and underscores", opts.Name)
	}

	// Create directory structure: .chief/prds/<name>/
	prdDir := filepath.Join(opts.BaseDir, ".chief", "prds", opts.Name)
	if err := os.MkdirAll(prdDir, 0755); err != nil {
		return fmt.Errorf("failed to create PRD directory: %w", err)
	}

	// Check if prd.md already exists
	prdMdPath := filepath.Join(prdDir, "prd.md")
	if _, err := os.Stat(prdMdPath); err == nil {
		return fmt.Errorf("PRD already exists at %s. Use 'chief edit %s' to modify it", prdMdPath, opts.Name)
	}

	// Get the init prompt
	prompt := embed.GetInitPrompt(opts.Context)

	// Launch interactive Claude session
	fmt.Printf("Creating PRD in %s...\n", prdDir)
	fmt.Println("Launching Claude to help you create your PRD...")
	fmt.Println()

	if err := runInteractiveClaude(prdDir, prompt); err != nil {
		return fmt.Errorf("Claude session failed: %w", err)
	}

	// Check if prd.md was created
	if _, err := os.Stat(prdMdPath); os.IsNotExist(err) {
		fmt.Println("\nNo prd.md was created. Run 'chief init' again to try again.")
		return nil
	}

	fmt.Println("\nPRD created successfully!")

	// Run conversion from prd.md to prd.json
	fmt.Println("Converting prd.md to prd.json...")
	if err := RunConvert(prdDir); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	fmt.Printf("\nYour PRD is ready! Run 'chief' or 'chief %s' to start working on it.\n", opts.Name)
	return nil
}

// runInteractiveClaude launches an interactive Claude session in the specified directory.
func runInteractiveClaude(workDir, prompt string) error {
	cmd := exec.Command("claude", "-p", prompt)
	cmd.Dir = workDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ConvertOptions contains configuration for the conversion command.
type ConvertOptions struct {
	PRDDir string // PRD directory containing prd.md
	Merge  bool   // Auto-merge without prompting on conversion conflicts
	Force  bool   // Auto-overwrite without prompting on conversion conflicts
}

// RunConvert converts prd.md to prd.json using Claude.
// This is a placeholder that will be fully implemented in US-018.
func RunConvert(prdDir string) error {
	return RunConvertWithOptions(ConvertOptions{PRDDir: prdDir})
}

// RunConvertWithOptions converts prd.md to prd.json using Claude with options.
// The Merge and Force flags will be fully implemented in US-019.
func RunConvertWithOptions(opts ConvertOptions) error {
	prdMdPath := filepath.Join(opts.PRDDir, "prd.md")
	prdJsonPath := filepath.Join(opts.PRDDir, "prd.json")

	// Check if prd.md exists
	if _, err := os.Stat(prdMdPath); os.IsNotExist(err) {
		return fmt.Errorf("prd.md not found in %s", opts.PRDDir)
	}

	// Read prd.md content
	content, err := os.ReadFile(prdMdPath)
	if err != nil {
		return fmt.Errorf("failed to read prd.md: %w", err)
	}

	// Create conversion prompt
	conversionPrompt := fmt.Sprintf(`Convert the following PRD markdown to a valid JSON file.

Output ONLY the JSON content, no markdown code blocks or explanations.

The JSON should have this structure:
{
  "project": "Project Name",
  "description": "Brief project description",
  "userStories": [
    {
      "id": "US-001",
      "title": "Story Title",
      "description": "Full description",
      "acceptanceCriteria": ["criterion 1", "criterion 2"],
      "priority": 1,
      "passes": false
    }
  ]
}

PRD Content:
%s`, string(content))

	// Run Claude one-shot conversion
	cmd := exec.Command("claude",
		"--dangerously-skip-permissions",
		"-p", conversionPrompt,
		"--output-format", "text",
	)
	cmd.Dir = opts.PRDDir

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("conversion failed: %s", string(exitErr.Stderr))
		}
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Clean up output (remove any markdown code blocks if present)
	jsonContent := cleanJSONOutput(string(output))

	// Write prd.json
	if err := os.WriteFile(prdJsonPath, []byte(jsonContent), 0644); err != nil {
		return fmt.Errorf("failed to write prd.json: %w", err)
	}

	// Verify it's valid JSON by attempting to load it
	// This will be done by the PRD loader when the TUI starts

	return nil
}

// cleanJSONOutput removes markdown code blocks and trims whitespace.
func cleanJSONOutput(output string) string {
	output = strings.TrimSpace(output)

	// Remove markdown code blocks if present
	if strings.HasPrefix(output, "```json") {
		output = strings.TrimPrefix(output, "```json")
	} else if strings.HasPrefix(output, "```") {
		output = strings.TrimPrefix(output, "```")
	}

	if strings.HasSuffix(output, "```") {
		output = strings.TrimSuffix(output, "```")
	}

	return strings.TrimSpace(output)
}

// isValidPRDName checks if the name contains only valid characters.
func isValidPRDName(name string) bool {
	if name == "" {
		return false
	}
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}
