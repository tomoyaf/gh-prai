package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func loadTemplate(templatePath string) string {
	if templatePath == "default" {
		return getDefaultTemplate()
	}
	
	if templatePath == "" {
		templatePath = filepath.Join(".github", "pull_request_template.md")
	}
	content, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		fmt.Println("Using default template instead.")
		return getDefaultTemplate()
	}

	fmt.Printf("Using template from %s\n", templatePath)
	fmt.Printf("Template content:\n%s\n", string(content))

	return string(content)
}

func getDefaultTemplate() string {
	return `## 概要
<!-- 変更の概要簡潔に記述してください -->

## 変更内容
<!-- 具体的な変更内容を箇条書きで記述してください -->

## その他
<!-- その他、変更の詳細や注意すべき点などレビュアーに伝えたいことがあれば簡潔に記述してください -->
`
}

func getDefaultPrompt() string {
	return `You are an AI assistant that generates concise and informative Pull Request descriptions based on the provided diff and template. Please fill in the template with relevant information extracted from the diff. Be specific and focus on the key changes and their impact.`
}
