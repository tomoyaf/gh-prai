## PRAI - Automate Your Pull Request Titles and Descriptions

[![release](https://github.com/tomoyaf/gh-prai/actions/workflows/release.yml/badge.svg)](https://github.com/tomoyaf/gh-prai/actions/workflows/release.yml)

[GitHub CLI Extension](https://docs.github.com/ja/github-cli/github-cli/about-github-cli)

PRAI (Pull Request AI) is a GitHub CLI extension that automates the creation of pull request titles and descriptions by leveraging the power of OpenAI's ChatGPT API. Whether you're a solo developer or part of a team, PRAI helps you save time by generating thoughtful, concise, and well-structured PR summaries based on your git diffs.

### Key Features

- **Automatic PR Title and Description Generation:** PRAI generates pull request titles and bodies based on your git diff, removing the need for manual writing.
- **Integration with GitHub CLI:** Easily create new pull requests or update existing ones with a few simple commands.
- **Flexible Configuration:** Customize the PR templates, language, and prompts to fit your workflow.

### DEMO
### Default Usage

https://github.com/user-attachments/assets/b09d5c22-6711-4bbc-b652-3675b922c0fe


### Configure Language to English

https://github.com/user-attachments/assets/011514dc-348e-417d-b753-ee3e7144c87f

## Why Use PRAI?
- **Time-saving:** Focus on coding while PRAI writes your PR descriptions.
- **Consistency:** Maintain a uniform style and tone in your PRs across your projects.
- **AI-Powered:** PRAI uses ChatGPT to analyze your code changes and summarize them effectively.

## Installation
Install the extension using GitHub CLI:

```bash
gh extension install tomoyaf/gh-prai
```

## Usage
### Step 1: Set up Your OpenAI API Key
Before you can start using PRAI, configure your OpenAI API key:
```bash
gh prai config api_key YOUR_OPENAI_API_KEY
```

### Step 2: Generate a Pull Request
To automatically generate the title and body of your pull request, simply run:

```bash
gh prai # or 'gh prai create'
```

### Additional Configurations
**Language:** Set the language for the PR title and description (default: English).
```bash
gh prai config language en  # or 'ja'
```
**Template:** Customize the template used for PR descriptions.
```bash
gh prai config template default  # or './.github/PULL_REQUEST_TEMPLATE/mytemplate.md'
```
**Custom Prompts:** Tailor the AI's behavior by providing a custom prompt.
```bash
gh prai config prompt "Your custom prompt"
```

## Help and Documentation
For more details on available commands and options:
```bash
gh prai --help
```

## Contributing
We welcome contributions! Please feel free to submit issues or pull requests to help improve PRAI.

## License
This project is licensed under the MIT License.
