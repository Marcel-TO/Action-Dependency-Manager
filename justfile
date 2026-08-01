# Synchronization of dependencies.
[group("Dev")]
tidy:
    go mod tidy

# Run Go vet to examine Go source code and report suspicious constructs.
[group("Dev")]
vet:
    go vet ./...

# Format Go source code files.
[group("Dev")]
format:
    go fmt ./...

# Run check operation to discover dependency and required updates
[group("Run")]
check pull="false":
    go run main.go check --pull="{{pull}}"

# Run check operation to discover dependency and required updates for a specific repos.
[group("Run")]
check-selected select pull="false":
    go run main.go check --select="{{select}}" --pull="{{pull}}"

# Run update operation to update dependencies to the latest versions
[group("Run")]
update pull="false" commit="false":
    go run main.go update --pull="{{pull}}" --commit="{{commit}}"

# Run update operation to update dependencies to the latest versions for a specific repos.
[group("Run")]
update-selected select pull="false" commit="false":
    go run main.go update --select="{{select}}" --pull="{{pull}}" --commit="{{commit}}"

# Run the main.go file with additional arguments passed to it.
[group("Run")]
run +ARGS:
    go run main.go {{ARGS}}
