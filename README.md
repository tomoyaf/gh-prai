## PRAI
1. Generate the title and body of the Pull Request with ChatGPT API based on git diff
2. Create a new Pull Request or update an existing Pull Request based on this title and body.


`↓ default`

https://github.com/user-attachments/assets/b09d5c22-6711-4bbc-b652-3675b922c0fe


`↓ gh prai config language en`

https://github.com/user-attachments/assets/a203b48e-5937-4cd7-a109-c184a6571e6c



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

### Description

- Uses the ChatGPT API to generate the body of a Pull Request
- Contains `gh pr create`, which performs up to PR creation.
