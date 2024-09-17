## PRAI

### How to install?

`$ gh extension install tomoyaf/gh-prai`

### How to use?

`$ gh prai config api_key YOUR_OPENAI_API_KEY`

`$ gh prai # or 'gh prai create'`

```
$ gh prai config language ja  # or 'en'
$ gh prai config template default  # or 'local'
$ gh prai config prompt "Your custom prompt"
```

### Description

- ChatGPT API を使って、Pull Request の本文を生成します
- `gh pr create`を内包しており、PR 作成まで実行します
