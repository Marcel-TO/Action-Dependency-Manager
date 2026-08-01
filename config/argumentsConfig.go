package config

import "flag"

type UpdateArgumentConfig struct {
	Commit   bool
	Selected string
	Debug    bool
	Pull     bool
}

type CheckArgumentConfig struct {
	Selected string
	Debug    bool
	Pull     bool
}

// ParseUpdateArguments parses the command-line arguments for the update operation.
func ParseUpdateArguments(cmd *flag.FlagSet, args []string) UpdateArgumentConfig {
	commit := cmd.Bool("commit", true, "Indicates whether to commit the changes to the repository after updating the game version and building")
	selected := cmd.String("select", "", "Comma-separated list of selected repositories to update. If empty, all repositories will be updated.")
	debug := cmd.Bool("debug", false, "Enable debug mode for verbose logging")
	pull := cmd.Bool("pull", false, "Indicates whether to pull the latest changes from the remote repository before updating the game version and building")

	cmd.Parse(args)

	return UpdateArgumentConfig{
		Commit:   *commit,
		Selected: *selected,
		Debug:    *debug,
		Pull:     *pull,
	}
}

// ParseCheckArguments parses the command-line arguments for the check operation.
func ParseCheckArguments(cmd *flag.FlagSet, args []string) CheckArgumentConfig {
	selected := cmd.String("select", "", "Comma-separated list of selected repositories to check. If empty, all repositories will be checked.")
	debug := cmd.Bool("debug", false, "Enable debug mode for verbose logging")
	pull := cmd.Bool("pull", false, "Indicates whether to pull the latest changes from the remote repository before checking the game version and building")

	cmd.Parse(args)

	return CheckArgumentConfig{
		Selected: *selected,
		Debug:    *debug,
		Pull:     *pull,
	}
}
