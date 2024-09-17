package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
)

var baseBranch string

func init() {
	flag.StringVar(&baseBranch, "base", "", "Specify the base branch for the PR")
}

func createPR() {
	colorPrint := color.New(color.FgHiGreen, color.Bold)
	errorPrint := color.New(color.FgHiRed, color.Bold)
	
	flag.Parse()

	config := loadConfig()

	if config.APIKey == "" {
		errorPrint.Println("OpenAI API key is not set. Please set it using 'gh prai config api_key YOUR_API_KEY'\nsee: https://platform.openai.com/api-keys")
		os.Exit(1)
	}

	if baseBranch == "" {
		var err error
		baseBranch, err = getDefaultBranch()
		if err != nil {
			errorPrint.Printf("Error getting default branch: %v\n", err)
			os.Exit(1)
		}
	}
	
	headBranch, _ := getCurrentBranch()
	existingPR, err := checkExistingPR(baseBranch, headBranch)
	if err != nil {
		errorPrint.Printf("Error checking for existing PR: %v\n", err)
		os.Exit(1)
	}
	if existingPR != nil {
		fmt.Printf("An existing PR (#%d) was found:\n\n", existingPR.Number)
		pullRequestUrl := getPullRequestUrl(existingPR.Number)
		colorPrint.Printf("%s #%d\n", existingPR.Title, existingPR.Number)
		colorPrint.Println(pullRequestUrl)
		fmt.Print("\n")
		if !promptUser("\nDo you want to update this PR? ([y]/n): ") {
			fmt.Println("Operation cancelled.")
			return
		}
	}

	fmt.Print("\n")

	diff, err := getPRDiff(baseBranch)
	if err != nil {
		errorPrint.Printf("Error getting PR diff: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n🤖 Title")
	template := loadTemplate(config.Template)
	title, err := generatePRTitle(diff, config)
	if err != nil {
		errorPrint.Printf("Error generating PR title: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n🤖 Description")
	description, err := generatePRDescription(diff, template, config)
	if err != nil {
		errorPrint.Printf("Error generating PR description: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("\n")

	prompt := "\nDo you want to create a PR with this title and description? ([y]/n): "
	if existingPR != nil {
		prompt = fmt.Sprintf("\nDo you want to update the existing PR (#%d) with this title and description? ([y]/n): ", existingPR.Number)
	}
	confirmCreate := promptUser(prompt)
	for !confirmCreate {
		title = promptForEdit("title", title)
		description = promptForEdit("description", description)

		fmt.Println("🤖 Title")
		colorPrint.Print(title)
		fmt.Println("\n🤖 Description:")
		colorPrint.Print(description)

		confirmCreate = promptUser(prompt)
	}

	if existingPR != nil {
		fmt.Print("\n\n")
		err = updatePR(existingPR.Number, title, description)
		if err != nil {
			errorPrint.Printf("Error updating PR: %v\n", err)
			os.Exit(1)
		}
		
		pullRequestUrl := getPullRequestUrl(existingPR.Number)

		colorPrint.Printf("\n\n%s #%d\n%s\n\n", title, existingPR.Number, pullRequestUrl)
		fmt.Printf("Pull Request updated successfully!\n")
	} else {
		fmt.Print("\n\n")
		err = executePRCreate(title, description, baseBranch)
		if err != nil {
			errorPrint.Printf("Error creating PR: %v\n", err)
			os.Exit(1)
		}
		createdPR, err := checkExistingPR(baseBranch, headBranch)
		if err != nil {
			errorPrint.Printf("Error checking for created PR: %v\n", err)
			os.Exit(1)
		}

		pullRequestUrl := getPullRequestUrl(createdPR.Number)

		colorPrint.Printf("\n\n%s #%d\n%s\n\n", title, createdPR.Number, pullRequestUrl)
		fmt.Println("Pull Request created successfully!")
	}
}

type PullRequest struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
}

func getPullRequestUrl(pullRequestNumber int) string {
	cmd := exec.Command("gh", "pr", "view", fmt.Sprintf("%d", pullRequestNumber), "--json", "url", "--jq", ".url")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func checkExistingPR(baseBranch, headBranch string) (*PullRequest, error) {
	cmd := exec.Command(
		"gh", "pr", "list",
		"--state", "open", "--json", "number,title",
		"-B", baseBranch, "-H", headBranch,
	)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Exit status 1 means no PR exists, which is not an error for us
			if exitError.ExitCode() == 1 {
				return nil, nil
			}
		}
		return nil, err
	}

	var pullRequests []PullRequest
	err = json.Unmarshal(output, &pullRequests)
	if err != nil {
		return nil, fmt.Errorf("error parsing PR data: %v", err)
	}
	if len(pullRequests) == 0 {
		return nil, nil
	}

	return &pullRequests[0], nil
}

func updatePR(number int, title, body string) error {
	cmd := exec.Command("gh", "pr", "edit", fmt.Sprintf("%d", number), "--title", title, "--body", body)
	return cmd.Run()
}

func getDefaultBranch() (string, error) {
	cmd := exec.Command("gh", "repo", "view", "--json=defaultBranchRef", "--jq", ".defaultBranchRef.name")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getPRDiff(baseBranch string) (string, error) {
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return "", err
	}

	ignoreFiles := []string{
		// TODO: Move this to a config file and allow users to customize
		"package-lock.json",
		"composer.lock",
		"*.lock",
		"go.sum",
		"go.mod",
	}
	args := []string{"diff", fmt.Sprintf("origin/%s...%s", baseBranch, currentBranch), "--"}
	for _, file := range ignoreFiles {
		args = append(args, fmt.Sprintf(":!%s", file))
	}

	cmd := exec.Command(
		"git", args...,
	)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	diff := string(output)
	if diff == "" {
		fmt.Printf("origin/%s...%s: No changes to create a PR for.\n", baseBranch, currentBranch)
		os.Exit(0)
	}
	return diff, nil
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func generatePRDescription(diff, template string, config Config) (string, error) {
	colorPrint := color.New(color.FgHiGreen, color.Bold)
	client := openai.NewClient(config.APIKey)

	req := openai.ChatCompletionRequest{
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
		Stream: true,
	}
		
	ctx := context.Background()

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var fullResponse strings.Builder
	
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Print("\n")
			return fullResponse.String(), nil
		}

		if err != nil {
			return "", err
		}

		content := response.Choices[0].Delta.Content
		colorPrint.Print(content)
		fullResponse.WriteString(content)
	}
}

func generatePRTitle(diff string, config Config) (string, error) {
	colorPrint := color.New(color.FgHiGreen, color.Bold)
	client := openai.NewClient(config.APIKey)

	req := openai.ChatCompletionRequest{
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
		Stream: true,
	}
		
	ctx := context.Background()

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", err
	}
	defer stream.Close()
	
	var fullResponse strings.Builder
	
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Print("\n")
			return fullResponse.String(), nil
		}

		if err != nil {
			return "", err
		}

		content := response.Choices[0].Delta.Content
		colorPrint.Print(content)
		fullResponse.WriteString(content)
	}
}

func executePRCreate(title, body, baseBranch string) error {
	cmd := exec.Command("gh", "pr", "create", "--title", title, "--body", body, "--base", baseBranch)
	return cmd.Run()
}

func promptForEdit(fieldName, content string) string {
	errorPrint := color.New(color.FgHiRed, color.Bold)

	for {
		fmt.Printf("Do you want to edit the %s? ([y]/n): ", fieldName)
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) == "n" {
			return content
		}

		editedContent, err := editInEditor(content)
		if err != nil {
			errorPrint.Printf("Error editing %s: %v\n", fieldName, err)
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

