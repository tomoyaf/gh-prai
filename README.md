## PRAI
[GitHub CLI Extension](https://docs.github.com/ja/github-cli/github-cli/about-github-cli)

ChatGPT API generates Pull Requests title and description.

1. Generate the title and body of the Pull Request with ChatGPT API based on git diff
2. Create a new Pull Request or update an existing Pull Request based on this title and body.


`↓ default`

https://github.com/user-attachments/assets/b09d5c22-6711-4bbc-b652-3675b922c0fe


`↓ gh prai config language en`

https://github.com/user-attachments/assets/011514dc-348e-417d-b753-ee3e7144c87f


### How to install?

```shell
gh extension install tomoyaf/gh-prai
```

### How to use?

```shell
gh prai config api_key YOUR_OPENAI_API_KEY
```

```shell
gh prai # or 'gh prai create'
```

```shell
gh prai config language ja  # or 'en'
```
```shell
gh prai config template default  # or 'local'
```
```shell
gh prai config prompt "Your custom prompt"
```
```shell
gh prai --help # `gh prai -h` and `gh prai config -h` and `gh prai create -h`
```

### Description

- Uses the ChatGPT API to generate the body of a Pull Request
- Contains `gh pr create` (and `gh pr`), which performs up to PR creation.
