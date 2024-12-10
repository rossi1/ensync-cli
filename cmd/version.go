package cmd

import (
	"github.com/rossi1/ensync-cli/pkg/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	var jsonFormat bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonFormat {
				return printJSON(cmd.OutOrStdout(), version.Get())
			}

			cmd.Println(version.String())
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonFormat, "json", false, "Output version information as JSON")

	return cmd
}
