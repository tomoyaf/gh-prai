package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func createPR() {
	config := loadConfig()

	if config.APIKey == "" {
		fmt.Println("OpenAI API key is not set. Please set it using 'gh prai config api_key YOUR_API_KEY'\nsee: https://platform.openai.com/api-keys")
		os.Exit(1)
	}

	diff, err := getPRDiff()
	if err != nil {
		fmt.Printf("Error getting PR diff: %v\n", err)
		os.Exit(1)
	}

	template := loadTemplate(config.Template)
	description, err := generatePRDescription(diff, template, config)
	if err != nil {
		fmt.Printf("Error generating PR description: %v\n", err)
		os.Exit(1)
	}

	title, err := generatePRTitle(diff, config)
	if err != nil {
		fmt.Printf("Error generating PR title: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("title:")
	fmt.Println(title)
	fmt.Println("\ndescription:")
	fmt.Println(description)

	confirmCreate := promptUser("Do you want to create a PR with this title and description? ([y]/n): ")
	for !confirmCreate {
		title = promptForEdit("title", title)
		description = promptForEdit("description", description)

		fmt.Println("title:")
		fmt.Println(title)
		fmt.Println("\ndescription:")
		fmt.Println(description)

		confirmCreate = promptUser("Do you want to create a PR with this title and description? ([y]/n): ")
	}

	err = executePRCreate(title, description)
	if err != nil {
		fmt.Printf("Error creating PR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Pull Request created successfully!")
}

func getPRDiff() (string, error) {
	cmd := exec.Command("gh", "pr", "diff")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func generatePRDescription(diff, template string, config Config) (string, error) {
	client := openai.NewClient(config.APIKey)
	resp, err := client.CreateChatCompletion(
    context.Background(),
    openai.ChatCompletionRequest{
        Model: openai.GPT4oMini,
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleSystem,
                Content: `You are an AI assistant specialized in creating concise and informative Pull Request (PR) descriptions. Your task is to analyze the provided code diff and generate a clear, structured PR description that focuses on essential information. Follow these guidelines:

1. Title: Use the first line of the description as a clear, concise title that summarizes the main purpose of the changes.

2. Overview: Provide a brief, high-level summary of the changes and their purpose. Limit this to 1-2 sentences.

3. Detailed Changes: List all specific changes made, using bullet points. Group related changes and use sub-bullets for more detailed points. Focus on:
   - Files modified
   - Sections added, removed, or renamed
   - New features or functionality added
   - Specific improvements or clarifications made

4. Keep the description concise and to the point. Avoid unnecessary explanations or background information unless absolutely crucial for understanding the changes.

5. Use appropriate Markdown syntax for formatting, especially for bullet points and sub-bullets.

6. Do not include sections for 'Related Issues', 'Testing Instructions', 'Performance Impact', or any other additional information unless explicitly present in the diff or absolutely necessary for understanding the changes.

7. Ensure the description is free of grammatical errors and uses clear, professional language.

The goal is to create a PR description that provides all necessary information about the changes in a brief, easily scannable format.`,
            },
            {
                Role:    openai.ChatMessageRoleUser,
                Content: fmt.Sprintf("Generate a Pull Request description in %s for the following diff, using this template:\n\nTemplate:\n%s\n\nDiff:\n%s", config.Language, template, diff),
            },
        },
        MaxTokens: 800,
    },
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func generatePRTitle(diff string, config Config) (string, error) {
	client := openai.NewClient(config.APIKey)
	resp, err := client.CreateChatCompletion(
    context.Background(),
    openai.ChatCompletionRequest{
        Model: openai.GPT4oMini,
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleSystem,
                Content: `You are an AI assistant that generates concise, informative, and impactful Pull Request titles based on the provided diff. Strictly adhere to these rules:
                1. Start with an English type prefix (feat, fix, docs, style, refactor, test, chore) followed by a colon and a space.
                2. Use the specified language for the main content of the title, unless English terms are more appropriate or widely used in the tech context.
                3. Use present tense, imperative mood verbs (e.g., "Add", "Update", "Fix", "Implement" or their equivalents in the specified language).
                4. Be extremely specific about the changes, focusing on the most important aspect.
                5. Include the affected component, file, or module name.
                6. Keep the title between 30-50 characters, prioritizing brevity and clarity.
                7. If there's a ticket number, include it in square brackets at the beginning.
                8. For breaking changes, start with "[BREAKING]".
                9. Mention the programming language or framework only if it's the main focus of the change.
                10. Use technical terms precisely and avoid general descriptions.
                11. Exclude articles and unnecessary words to maximize information density.
                12. For documentation changes, specify the exact nature of the update.
                13. Blend languages appropriately if technical terms are better left in English.
                Remember, the title should allow developers to immediately understand the core change without reading the full diff.`,
            },
            {
                Role:    openai.ChatMessageRoleUser,
                Content: fmt.Sprintf("Generate a short, impactful, and descriptive Pull Request title in %s for the following diff:\n\n%s", config.Language, diff),
            },
        },
        MaxTokens: 60,
    },
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func executePRCreate(title, body string) error {
	cmd := exec.Command("gh", "pr", "create", "--title", title, "--body", body)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func promptForEdit(fieldName, content string) string {
	for {
		fmt.Printf("Do you want to edit the %s? (y/n): ", fieldName)
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" {
			return content
		}

		editedContent, err := editInEditor(content)
		if err != nil {
			fmt.Printf("Error editing %s: %v\n", fieldName, err)
			continue
		}

		return editedContent
	}
}

func editInEditor(content string) (string, error) {
	tempFile, err := os.CreateTemp("", "pr-edit-*")
	
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(content); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %v", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"  // Default to vim if EDITOR is not set
	}

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run editor: %v", err)
	}

	editedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read edited file: %v", err)
	}

	return string(editedContent), nil
}

func promptUser(prompt string) bool {
	fmt.Print(prompt)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) != "n"
}

