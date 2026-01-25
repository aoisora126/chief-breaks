package embed

import (
	_ "embed"
	"strings"
)

//go:embed prompt.txt
var promptTemplate string

//go:embed init_prompt.txt
var initPromptTemplate string

// GetPrompt returns the agent prompt with the PRD path substituted.
func GetPrompt(prdPath string) string {
	return strings.ReplaceAll(promptTemplate, "{{PRD_PATH}}", prdPath)
}

// GetInitPrompt returns the PRD generator prompt with optional context substituted.
func GetInitPrompt(context string) string {
	if context == "" {
		context = "No additional context provided. Ask the user what they want to build."
	}
	return strings.ReplaceAll(initPromptTemplate, "{{CONTEXT}}", context)
}
