package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/puppetlabs/wash/cmd/util"
	"github.com/puppetlabs/wash/plugin"
)

func rmCommand() *cobra.Command {
	use, aliases := generateShellAlias("rm")
	rmCmd := &cobra.Command{
		Use:     use + " <path> [<path>]",
		Aliases: aliases,
		Short:   "Deletes the entries at the specified paths",
		Long: `Deletes the entries at the specified paths, prompting the user for confirmation
before deleting each entry.`,
		Args: cobra.MinimumNArgs(1),
		RunE: toRunE(rmMain),
	}
	rmCmd.Flags().BoolP("recurse", "r", false, "Delete directories (parent entries)")
	rmCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	return rmCmd
}

func rmMain(cmd *cobra.Command, args []string) exitCode {
	paths := args
	recurse, err := cmd.Flags().GetBool("recurse")
	if err != nil {
		panic(err.Error())
	}
	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		panic(err.Error())
	}

	conn := cmdutil.NewClient()

	// ec => exitCode
	ec := 0
	for _, path := range paths {
		var deletionConfirmed bool
		if force || !plugin.IsInteractive() {
			deletionConfirmed = true
		} else {
			msg := fmt.Sprintf("remove %v?", path)
			input, err := plugin.Prompt(msg)
			if err != nil {
				cmdutil.ErrPrintf("failed to get confirmation: %v", err)
				return exitCode{1}
			}
			// Assume confirmation if input starts with "y" or "Y". This matches the built-in
			// rm.
			deletionConfirmed = len(input) > 0 && (input[0] == 'y' || input[0] == 'Y')
		}
		if !deletionConfirmed {
			continue
		}
		if err := conn.Delete(path, recurse); err != nil {
			ec = 1
			cmdutil.ErrPrintf("%v\n", err)
		}
		// Delete was successful
	}
	return exitCode{ec}
}
