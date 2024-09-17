package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func loadTemplate(templateName string) string {
	if templateName == "default" {
		return getDefaultTemplate()
	}

	templatePath := filepath.Join(".github", "PULL_REQUEST_TEMPLATE.md")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		fmt.Println("Using default template instead.")
		return getDefaultTemplate()
	}

	return string(content)
}

func getDefaultTemplate() string {
	return `## 概要
<!-- 変更の概要を記述してください -->

## 変更内容
<!-- 具体的な変更内容を箇条書きで記述してください -->

## 関連するIssue
<!-- 関連するIssueがあれば記述してください -->

## その他
<!-- その他、レビュアーに伝えたいことがあれば記述してください -->
`
}

func getDefaultPrompt() string {
	return `You are an AI assistant that generates concise and informative Pull Request descriptions based on the provided diff and template. Please fill in the template with relevant information extracted from the diff. Be specific and focus on the key changes and their impact.`
}
