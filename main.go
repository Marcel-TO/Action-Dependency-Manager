package main

import (
	"flag"
	"fmt"
	"os"

	"marcel-to/actions-manager/config"
	"marcel-to/actions-manager/exec"
	"marcel-to/actions-manager/logger"
	"marcel-to/actions-manager/operations"
)

const (
	appName    = "adm"
	configFile = "action-repositories.yaml"
)

// main is the entry point of the Action Dependency Manager (ADM) CLI application. It parses command-line arguments and dispatches to the appropriate operation (check or update).
func main() {
	log := logger.NewLogger("ADM", false)

	if len(os.Args) < 2 {
		printUsage(log)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help", "-h", "--help":
		// Support `adm help <command>` as well as bare `adm help`.
		if len(os.Args) > 2 {
			if !printCommandUsage(log, os.Args[2]) {
				log.Error("Unknown command: " + os.Args[2])
				printUsage(log)
				os.Exit(1)
			}
			return
		}
		printUsage(log)
	case "update":
		runUpdate(log, os.Args[2:])
	case "check":
		runCheck(log, os.Args[2:])
	default:
		log.Error("Unknown command: " + os.Args[1])
		printUsage(log)
		os.Exit(1)
	}
}

// runUpdate parses the update flags and executes the update operation. A
// -h/--help flag is handled by the flag package via the FlagSet's Usage func.
func runUpdate(log *logger.Logger, argv []string) {
	cmd := flag.NewFlagSet("update", flag.ExitOnError)
	cmd.Usage = func() { updateUsage(log) }
	args := config.ParseUpdateArguments(cmd, argv)

	requireGH(log)
	cfg := mustLoadConfig(log)
	if err := operations.RunUpdate(log, args, cfg); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

// runCheck parses the check flags and executes the check operation. A -h/--help
// flag is handled by the flag package via the FlagSet's Usage func.
func runCheck(log *logger.Logger, argv []string) {
	cmd := flag.NewFlagSet("check", flag.ExitOnError)
	cmd.Usage = func() { checkUsage(log) }
	args := config.ParseCheckArguments(cmd, argv)

	requireGH(log)
	cfg := mustLoadConfig(log)
	if err := operations.RunCheck(log, args, cfg); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

// requireGH ensures the GitHub CLI is installed, exiting with an error message
// when it is not, since every operation depends on it.
func requireGH(log *logger.Logger) {
	if err := exec.EnsureGHInstalled(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

// mustLoadConfig loads the repository configuration or exits with an error.
func mustLoadConfig(log *logger.Logger) *config.Config {
	log.Info("GH Action Dependency Manager")

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Error("Failed to load configuration: " + err.Error())
		os.Exit(1)
	}
	log.Info("Configuration loaded successfully!")
	return cfg
}

// printUsage logs the top-level help.
func printUsage(log *logger.Logger) {
	log.Info(fmt.Sprintf(`Action Dependency Manager (%s)
Keep the GitHub Action dependencies in your workflow files up to date.

Usage:
  %s <command> [flags]

Commands:
  check     Report which dependencies have a newer version (read-only).
  update    Pin every dependency to its latest release commit SHA.
  help      Show this help message.

Run '%s help <command>' or '%s <command> --help' for details on a command.

Examples:
  %s check
  %s check --select="1,3"
  %s update --select="1" --debug`,
		appName, appName, appName, appName, appName, appName, appName))
}

// printCommandUsage logs command-specific help for name and reports whether
// name was a known command.
func printCommandUsage(log *logger.Logger, name string) bool {
	switch name {
	case "update":
		updateUsage(log)
	case "check":
		checkUsage(log)
	case "help":
		printUsage(log)
	default:
		return false
	}
	return true
}

// checkUsage logs help for the check command.
func checkUsage(log *logger.Logger) {
	log.Info(fmt.Sprintf(`Report which GitHub Action dependencies have a newer version available.
This command is read-only and never modifies any files.

Usage:
  %s check [flags]

Flags:
  --select string   Comma-separated repository ids to check (default: all).
  --debug           Enable verbose debug logging.

Examples:
  %s check
  %s check --select="1,3"`,
		appName, appName, appName))
}

// updateUsage logs help for the update command.
func updateUsage(log *logger.Logger) {
	log.Info(fmt.Sprintf(`Pin every GitHub Action dependency to its latest release commit SHA,
keeping the release tag as a trailing comment (e.g. uses: actions/checkout@<sha> # v7.0.1).

Usage:
  %s update [flags]

Flags:
  --select string   Comma-separated repository ids to update (default: all).
  --commit          Commit the changes after updating (default: true).
  --debug           Enable verbose debug logging.

Examples:
  %s update
  %s update --select="1" --debug`,
		appName, appName, appName))
}
