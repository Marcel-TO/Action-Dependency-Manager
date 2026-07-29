# Action-Dependency-Manager
This repository aims to make updating GH Action dependencies across multiple repositories easier

## Prerequisites
The Action Manager relies on the gh cli tool to extract required information for upgrading the Action versions. Specifically it focuses on release tag and its corresponding commit SHA.

- GitHub Account
- gh CLI

### Process
```bash
# To fetch the latest tag name
gh release view --repo actions/checkout --json tagName --jq .tagName

# To fetch not annotated tag commit hash (type == "commit")
gh api repos/actions/checkout/git/ref/tags/v7.0.1 --jq .object.sha


# To fetch annotated tag commit hash
gh api repos/actions/checkout/git/refs/tags/v7.0.1 --jq '.object.sha'
# then, if type == "tag":
gh api repos/actions/checkout/git/tags/<sha> --jq '.object.sha'
```
